package change

import (
	"archive/zip"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
)

// job — одна картинка в очереди
type job struct {
	index int
	file  *multipart.FileHeader
}

// result — готовый байткод или ошибка
type result struct {
	index int
	data  []byte
	err   error
}

func Size(w http.ResponseWriter, r *http.Request) {
	// увеличил лимит до 50 и форму до 50мб
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "ошибка получения файлов", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	if len(files) == 0 || len(files) > 50 {
		http.Error(w, "нужно загрузить от 1 до 50 изображений", http.StatusBadRequest)
		return
	}

	width, err := strconv.Atoi(r.FormValue("width"))
	if err != nil || width <= 0 {
		http.Error(w, "некорректная ширина", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil || height <= 0 {
		http.Error(w, "некорректная высота", http.StatusBadRequest)
		return
	}

	// каналы для очереди задач и результатов
	tasks := make(chan job, len(files))
	out := make(chan result, len(files))

	// сколько воркеров запускать
	workers := runtime.NumCPU()
	if workers > len(files) {
		workers = len(files)
	}

	var wg sync.WaitGroup

	// запускаем пул воркеров
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				data, err := sizePhoto(task.file, width, height)
				out <- result{
					index: task.index,
					data:  data,
					err:   err,
				}
			}
		}()
	}

	// кидаем влидные файлы в очередь
	var validCount int
	for _, fh := range files {
		if fh.Size > 10<<20 {
			continue
		}
		tasks <- job{index: validCount, file: fh}
		validCount++
	}
	close(tasks)

	// чек если файл большой
	if validCount == 0 {
		http.Error(w, "все файлы слишком большие", http.StatusBadRequest)
		return
	}

	// вг ждет завершения для закрытия каналов
	go func() {
		wg.Wait()
		close(out)
	}()

	// собираем всё обратно в правильном порядке
	final := make([][]byte, validCount)
	var failed bool

	for res := range out {
		if res.err != nil {
			failed = true
			break
		}
		final[res.index] = res.data
	}

	if failed {
		http.Error(w, "ошибка обработки одного из изображений", http.StatusBadRequest)
		return
	}

	// пакуем в зип
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for i, data := range final {
		if data == nil {
			continue
		}
		f, err := zw.Create(fmt.Sprintf("image_%d.jpg", i+1))
		if err != nil {
			http.Error(w, "ошибка создания zip", http.StatusInternalServerError)
			return
		}
		if _, err := f.Write(data); err != nil {
			http.Error(w, "ошибка записи в zip", http.StatusInternalServerError)
			return
		}
	}

	if err := zw.Close(); err != nil {
		http.Error(w, "ошибка завершения zip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="images.zip"`)
	w.Write(buf.Bytes())
}

func sizePhoto(fh *multipart.FileHeader, w int, h int) ([]byte, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}

	resized := imaging.Resize(img, w, h, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.JPEG); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
