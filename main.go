package main

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"tidy/change"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/size", change.Size)
	mux.HandleFunc("/", home)

	// открываю порт перед запуском сервера
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}

	fmt.Println("PixelSize запущен: http://localhost:8080")

	if err := http.Serve(listener, mux); err != nil {
		fmt.Println("Ошибка работы сервера:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// не отдаю стрн а адресах
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, page)
}

const page = `<!doctype html>
<html lang="ru">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>PixelSize | редактор размера изображений</title>
    <link rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Crect width='32' height='32' rx='8' fill='%23111'/%3E%3Cpath fill='%23fff' fill-rule='evenodd' d='M8 7H18A6 6 0 0 1 18 19H13V25H8ZM13 11H17A2 2 0 0 1 17 15H13Z'/%3E%3Crect x='20' y='21' width='4' height='4' fill='%23fff'/%3E%3C/svg%3E">
    <style>
        :root { color-scheme: light; }

        * { box-sizing: border-box; }

        body {
            min-width: 320px;
            margin: 0;
            background: #fff;
            color: #111;
            font-family: "Segoe UI", system-ui, -apple-system, Roboto, Arial, sans-serif;
            font-size: 15px;
            line-height: 1.5;
        }

        button, input { font: inherit; }
        button { cursor: pointer; }

        button:focus-visible, .drop-area:focus-visible {
            outline: 2px solid #111;
            outline-offset: 3px;
        }

        .page {
            max-width: 1040px;
            margin: 0 auto;
            padding: 0 24px 40px;
        }

        /* шапка */
        .topbar {
            display: flex;
            height: 72px;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            border-bottom: 1px solid #e6e6e6;
        }

        .brand {
            display: inline-flex;
            align-items: center;
            gap: 10px;
            color: #111;
            font-size: 20px;
            font-weight: 700;
            letter-spacing: -0.5px;
            text-decoration: none;
        }

        .brand svg { display: block; flex: none; }

        .top-note { color: #777; font-size: 14px; }

        .intro { padding: 72px 0 48px; }

        h1 {
            max-width: 640px;
            margin: 0;
            font-size: clamp(36px, 5.5vw, 52px);
            font-weight: 700;
            letter-spacing: -1.5px;
            line-height: 1.08;
        }

        .intro-copy {
            max-width: 500px;
            margin: 18px 0 0;
            color: #666;
            font-size: 17px;
            line-height: 1.6;
        }

        .workspace {
            display: grid;
            grid-template-columns: minmax(0, 1fr) 360px;
            gap: 48px;
            align-items: start;
        }

        .source-panel { min-width: 0; }

        .section-heading {
            display: flex;
            align-items: baseline;
            justify-content: space-between;
            gap: 12px;
            margin-bottom: 12px;
        }

        .section-heading h2 { margin: 0; font-size: 16px; font-weight: 600; }
        .counter { color: #888; font-size: 13px; }

        /* зона загрузки */
        .drop-area {
            display: flex;
            min-height: 200px;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 28px 20px;
            border: 1.5px dashed #c9c9c9;
            border-radius: 14px;
            text-align: center;
            cursor: pointer;
            transition: border-color .15s, background-color .15s;
        }

        .drop-area:hover, .drop-area.dragging { border-color: #111; background: #fafafa; }

        .upload-icon { margin-bottom: 14px; }
        .drop-title { margin: 0; font-size: 17px; font-weight: 600; }
        .drop-subtitle { margin: 6px 0 0; color: #666; font-size: 14px; }

        .text-button {
            padding: 0;
            border: 0;
            background: none;
            color: #111;
            font-weight: 600;
            text-decoration: underline;
            text-underline-offset: 3px;
        }
        .text-button:hover { color: #555; }

        .drop-hint { margin: 14px 0 0; color: #999; font-size: 13px; }

        .files-section { margin-top: 36px; }

        .empty-files {
            padding: 18px 0;
            border-top: 1px solid #ececec;
            border-bottom: 1px solid #ececec;
            color: #888;
            font-size: 14px;
        }

        .file-row {
            display: flex;
            min-width: 0;
            align-items: center;
            gap: 14px;
            padding: 10px 0;
            border-top: 1px solid #ececec;
        }
        .file-row:last-child { border-bottom: 1px solid #ececec; }

        .file-preview {
            width: 48px;
            height: 48px;
            flex: 0 0 48px;
            border: 1px solid #e6e6e6;
            border-radius: 8px;
            background: #f5f5f5;
            object-fit: cover;
        }

        .file-info { min-width: 0; flex: 1; }

        .file-name {
            overflow: hidden;
            font-size: 14px;
            font-weight: 500;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .file-size { margin-top: 2px; color: #888; font-size: 13px; }

        .remove-button {
            width: 32px;
            height: 32px;
            flex: 0 0 32px;
            border: 0;
            border-radius: 8px;
            background: none;
            color: #888;
            font-size: 20px;
            line-height: 1;
        }
        .remove-button:hover { background: #f2f2f2; color: #111; }

        /* настройки справа */
        .settings-panel {
            position: sticky;
            top: 24px;
            padding: 28px;
            border: 1px solid #e3e3e3;
            border-radius: 16px;
        }

        .settings-title { margin: 0 0 8px; font-size: 20px; font-weight: 700; letter-spacing: -0.3px; }
        .panel-description { margin: 0 0 22px; color: #666; font-size: 14px; line-height: 1.6; }

        .settings { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
        .field { display: grid; min-width: 0; gap: 6px; }
        .field label { color: #333; font-size: 13px; font-weight: 500; }

        .field input {
            width: 100%;
            min-width: 0;
            height: 48px;
            padding: 0 12px;
            border: 1px solid #d2d2d2;
            border-radius: 10px;
            background: #fff;
            color: #111;
            font-size: 16px;
            -moz-appearance: textfield;
            appearance: textfield;
        }
        .field input::-webkit-outer-spin-button,
        .field input::-webkit-inner-spin-button { margin: 0; -webkit-appearance: none; }
        .field input::placeholder { color: #b0b0b0; }
        .field input:focus { border-color: #111; outline: 1px solid #111; }

        .field-hint { margin: 12px 0 0; color: #888; font-size: 13px; line-height: 1.55; }

        .result-line {
            display: flex;
            justify-content: space-between;
            gap: 12px;
            margin-top: 22px;
            padding: 14px 0;
            border-top: 1px solid #ececec;
            border-bottom: 1px solid #ececec;
            font-size: 14px;
        }
        .result-line span:first-child { color: #666; }
        .result-line span:last-child { font-weight: 600; }

        .actions { margin-top: 18px; }

        .submit-button {
            width: 100%;
            min-height: 50px;
            padding: 12px 16px;
            border: 0;
            border-radius: 10px;
            background: #111;
            color: #fff;
            font-size: 15px;
            font-weight: 600;
            transition: background-color .15s;
        }
        .submit-button:hover:not(:disabled) { background: #333; }
        .submit-button:disabled { background: #ececec; color: #a0a0a0; cursor: not-allowed; }

        .status {
            display: none;
            margin-top: 14px;
            color: #666;
            font-size: 14px;
            line-height: 1.5;
            overflow-wrap: anywhere;
        }
        .status.visible { display: block; }
        .status.error { color: #c62828; }

        .footer {
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            gap: 8px 24px;
            margin-top: 72px;
            padding-top: 20px;
            border-top: 1px solid #e6e6e6;
            color: #888;
            font-size: 13px;
            line-height: 1.6;
        }
        .footer span:first-child { max-width: 560px; }
        .footer span:last-child { color: #111; font-weight: 500; }

        /* планшеты и телефоны */
        @media (max-width: 820px) {
            .intro { padding: 48px 0 36px; }
            .workspace { grid-template-columns: 1fr; gap: 36px; }
            .settings-panel { position: static; }
        }

        @media (max-width: 520px) {
            .page { padding: 0 16px 32px; }
            .topbar { height: 64px; }
            .top-note { display: none; }
            .intro { padding: 36px 0 28px; }
            .intro-copy { font-size: 16px; }
            .drop-area { min-height: 170px; }
            .settings-panel { padding: 20px; }
            .footer { margin-top: 48px; }
        }
    </style>
</head>
<body>
    <main class="page">
        <header class="topbar">
            <a class="brand" href="/" aria-label="PixelSize, главная">
                <!-- логотип: буква P и маленький квадрат-пиксель -->
                <svg width="32" height="32" viewBox="0 0 32 32" aria-hidden="true">
                    <rect width="32" height="32" rx="8" fill="#111"></rect>
                    <path fill="#fff" fill-rule="evenodd" d="M8 7H18A6 6 0 0 1 18 19H13V25H8ZM13 11H17A2 2 0 0 1 17 15H13Z"></path>
                    <rect x="20" y="21" width="4" height="4" fill="#fff"></rect>
                </svg>
                PixelSize
            </a>
            <span class="top-note">Редактор размера изображений</span>
        </header>

        <section class="intro">
            <h1>Изображения в нужном размере</h1>
            <p class="intro-copy">Добавьте до 50-ти изображений, задайте ширину и высоту. Готовые файлы скачаются одним ZIP-архивом.</p>
        </section>

        <!-- слева файлы, справа настройки -->
        <form id="resize-form" class="workspace">
            <div class="source-panel">
                <section>
                    <div class="section-heading">
                        <h2>Исходные файлы</h2>
                        <span id="file-count" class="counter">0 из 50</span>
                    </div>
                    <div id="drop-area" class="drop-area" tabindex="0" role="region" aria-label="Выбрать изображения или перетащить их сюда">
                        <svg class="upload-icon" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                            <path d="M12 16V4"></path><path d="m7 9 5-5 5 5"></path><path d="M5 14v5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-5"></path>
                        </svg>
                        <p class="drop-title">Перетащите изображения сюда</p>
                        <p class="drop-subtitle">или <button id="choose-files" class="text-button" type="button">выберите файлы</button></p>
                        <p class="drop-hint">Можно выбрать несколько файлов сразу</p>
                        <input id="file-input" type="file" accept="image/*" multiple hidden>
                    </div>
                </section>

                <section class="files-section">
                    <div class="section-heading">
                        <h2>Ваши изображения</h2>
                    </div>
                    <div id="file-list" class="file-list">
                        <div class="empty-files">Пока ничего не добавлено. Загрузите изображения выше.</div>
                    </div>
                </section>
            </div>

            <aside class="settings-panel">
                <h2 class="settings-title">Настройки изображения</h2>
                <p class="panel-description">Укажите ширину и высоту в пикселях. Значения применятся ко всем файлам.</p>

                <div class="settings">
                    <div class="field">
                        <label for="width">Ширина, px</label>
                        <input id="width" name="width" type="number" min="1" max="5000" step="1" inputmode="numeric" placeholder="1200" required>
                    </div>
                    <div class="field">
                        <label for="height">Высота, px</label>
                        <input id="height" name="height" type="number" min="1" max="5000" step="1" inputmode="numeric" placeholder="800" required>
                    </div>
                </div>
                <p class="field-hint">Максимум 5000 px по каждой стороне. Готовые файлы сохраняются в формате JPEG.</p>

                <div class="result-line">
                    <span>Итоговый размер</span>
                    <span id="result-size">не задан</span>
                </div>

                <div class="actions">
                    <button id="submit-button" class="submit-button" type="submit" disabled>Изменить размер и скачать ZIP</button>
                    <div id="status" class="status" role="status" aria-live="polite"></div>
                </div>
            </aside>
        </form>

        <footer class="footer">
            <span>Не теряйте время. Меняем размер всех картинок одновременно.</span>
            <span>mctavich / size editor</span>
        </footer>
    </main>

    <script>
        const fileInput = document.getElementById("file-input");
        const chooseFiles = document.getElementById("choose-files");
        const dropArea = document.getElementById("drop-area");
        const fileList = document.getElementById("file-list");
        const fileCount = document.getElementById("file-count");
        const form = document.getElementById("resize-form");
        const submitButton = document.getElementById("submit-button");
        const statusBox = document.getElementById("status");
        const widthInput = document.getElementById("width");
        const heightInput = document.getElementById("height");
        const resultSize = document.getElementById("result-size");

        // сюда складываю выбранные изображенияю
        let selectedFiles = [];
        const maxFiles = 50;

        // 5к макс размер
        const maxSize = 5000;

        function formatSize(size) {
            if (size < 1024) return size + " Б";
            if (size < 1024 * 1024) return (size / 1024).toFixed(1) + " КБ";
            return (size / (1024 * 1024)).toFixed(1) + " МБ";
        }

        function showStatus(message, isError) {
            statusBox.textContent = message;
            statusBox.className = "status visible" + (isError ? " error" : "");
        }

        function hideStatus() {
            statusBox.textContent = "";
            statusBox.className = "status";
        }

        // показываю итоговый размер, пока вводят числа
        function showResultSize() {
            const width = Number(widthInput.value);
            const height = Number(heightInput.value);

            if (width > 0 && height > 0) {
                resultSize.textContent = width + " × " + height + " px";
            } else {
                resultSize.textContent = "не задан";
            }
        }

        function addFiles(files) {
            hideStatus();

            for (const file of files) {
                if (!file.type.startsWith("image/")) {
                    showStatus("Можно добавлять только изображения", true);
                    continue;
                }

                if (selectedFiles.length >= maxFiles) {
                    showStatus("Можно загрузить не больше 50-ти изображений", true);
                    break;
                }

                selectedFiles.push(file);
            }

            // обновляю очередь после добавления файлоы
            drawFiles();
            fileInput.value = "";
        }

        function drawFiles() {
            fileList.replaceChildren();
            fileCount.textContent = selectedFiles.length + " из " + maxFiles;
            submitButton.disabled = selectedFiles.length === 0;

            if (selectedFiles.length === 0) {
                const empty = document.createElement("div");
                empty.className = "empty-files";
                empty.textContent = "Пока ничего не добавлено. Загрузите изображения выше.";
                fileList.append(empty);
                return;
            }

            selectedFiles.forEach(function (file, index) {
                const row = document.createElement("div");
                row.className = "file-row";

                const preview = document.createElement("img");
                preview.className = "file-preview";
                preview.alt = "";
                const previewUrl = URL.createObjectURL(file);
                preview.src = previewUrl;
                preview.onload = function () { URL.revokeObjectURL(previewUrl); };
                preview.onerror = function () { URL.revokeObjectURL(previewUrl); };

                const info = document.createElement("div");
                info.className = "file-info";

                const name = document.createElement("div");
                name.className = "file-name";
                name.textContent = file.name;
                name.title = file.name;

                const size = document.createElement("div");
                size.className = "file-size";
                size.textContent = formatSize(file.size);

                // создаю кнопку удаления для каждой картинки
                const removeButton = document.createElement("button");
                removeButton.className = "remove-button";
                removeButton.type = "button";
                removeButton.textContent = "×";
                removeButton.setAttribute("aria-label", "Удалить " + file.name);
                removeButton.addEventListener("click", function () {
                    selectedFiles.splice(index, 1);
                    hideStatus();
                    drawFiles();
                });

                info.append(name, size);
                row.append(preview, info, removeButton);
                fileList.append(row);
            });
        }

        chooseFiles.addEventListener("click", function () {
            fileInput.click();
        });

        fileInput.addEventListener("change", function () {
            addFiles(fileInput.files);
        });

        dropArea.addEventListener("click", function (event) {
            if (event.target.closest("button")) return;
            fileInput.click();
        });

        dropArea.addEventListener("keydown", function (event) {
            if (event.key === "Enter" || event.key === " ") {
                event.preventDefault();
                fileInput.click();
            }
        });

        ["dragenter", "dragover"].forEach(function (eventName) {
            dropArea.addEventListener(eventName, function (event) {
                event.preventDefault();
                dropArea.classList.add("dragging");
            });
        });

        ["dragleave", "drop"].forEach(function (eventName) {
            dropArea.addEventListener(eventName, function (event) {
                event.preventDefault();
                dropArea.classList.remove("dragging");
            });
        });

        dropArea.addEventListener("drop", function (event) {
            addFiles(event.dataTransfer.files);
        });

        widthInput.addEventListener("input", showResultSize);
        heightInput.addEventListener("input", showResultSize);

        form.addEventListener("submit", async function (event) {
            event.preventDefault();
            hideStatus();

            if (selectedFiles.length === 0) {
                showStatus("Сначала добавьте хотя бы одно изображение", true);
                return;
            }

            const width = Number(widthInput.value);
            const height = Number(heightInput.value);

            if (!Number.isInteger(width) || width <= 0 || !Number.isInteger(height) || height <= 0) {
                showStatus("Введите положительные целые значения ширины и высоты", true);
                return;
            }

            if (width > maxSize || height > maxSize) {
                showStatus("Максимальный размер по каждой стороне: " + maxSize + " px", true);
                return;
            }

            const data = new FormData();
            selectedFiles.forEach(function (file) {
                data.append("files", file);
            });
            data.append("width", String(width));
            data.append("height", String(height));

            submitButton.disabled = true;
            submitButton.textContent = "Обработка файлов";

            try {
                // отправляю файлы и размер одним запросом
                const response = await fetch("/size", {
                    method: "POST",
                    body: data
                });

                if (!response.ok) {
                    const message = await response.text();
                    throw new Error(message || "Не удалось обработать изображения");
                }

                const blob = await response.blob();
                const url = URL.createObjectURL(blob);
                const link = document.createElement("a");
                link.href = url;
                link.download = "images.zip";
                document.body.append(link);
                link.click();
                link.remove();
                setTimeout(function () { URL.revokeObjectURL(url); }, 1000);
                showStatus("Готово, архив скачивается", false);
            } catch (error) {
                showStatus(error.message || "Не удалось обработать изображения", true);
            } finally {
                submitButton.textContent = "Изменить размер и скачать ZIP";
                submitButton.disabled = selectedFiles.length === 0;
            }
        });
    </script>
</body>
</html>`
