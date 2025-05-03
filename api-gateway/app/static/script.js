const TYPE_FILE = 0;
const TYPE_YOUTUBE = 1;

document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('audioFile');
    const youtubeInput = document.getElementById('youtubeUrl');
    const resultDiv = document.getElementById('tabResult');
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');
    const copyButton = document.getElementById('copyButton');
    const progressContainer = document.getElementById('uploadProgress');
    const progressBar = progressContainer.querySelector('.progress-bar');
    const progressText = progressContainer.querySelector('.progress-text');

    // Переключение между вкладками
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabId}-tab`).classList.add('active');
        });
    });

    // Обработка выбора файла
    fileInput.addEventListener('change', (e) => {
        const fileName = e.target.files[0]?.name || '';
        document.getElementById('fileName').textContent = fileName;
    });

    // Drag and drop для файлов
    const fileUpload = document.querySelector('.file-upload');
    fileUpload.addEventListener('dragover', (e) => {
        e.preventDefault();
        fileUpload.classList.add('dragover');
    });

    fileUpload.addEventListener('dragleave', () => {
        fileUpload.classList.remove('dragover');
    });

    fileUpload.addEventListener('drop', (e) => {
        e.preventDefault();
        fileUpload.classList.remove('dragover');
        if (e.dataTransfer.files.length) {
            fileInput.files = e.dataTransfer.files;
            document.getElementById('fileName').textContent = e.dataTransfer.files[0].name;
        }
    });

    // Обработка формы
    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const activeTab = document.querySelector('.tab-button.active').getAttribute('data-tab');
        let formData = new FormData();
        const endpoint = 'http://localhost:8080/generate-tab';

        // Показываем индикатор загрузки
        document.getElementById('loading').classList.add('active');
        resultDiv.innerHTML = '<div class="tab-line">e|---------------------------------------------------------</div>' +
                             '<div class="tab-line">B|---------------------------------------------------------</div>' +
                             '<div class="tab-line">G|---------------------------------------------------------</div>' +
                             '<div class="tab-line">D|---------------------------------------------------------</div>' +
                             '<div class="tab-line">A|---------------------------------------------------------</div>' +
                             '<div class="tab-line">E|---------------------------------------------------------</div>';

        try {
            if (activeTab === 'file') {
                if (!fileInput.files.length) {
                    throw new Error('Выберите аудиофайл!');
                }
                
                // Показываем прогресс-бар для загрузки файла
                progressContainer.style.display = 'block';
                progressBar.style.width = '0%';
                progressText.textContent = '0%';
                
                // Симулируем прогресс загрузки (в реальном приложении используйте xhr.upload.onprogress)
                simulateUploadProgress();
                
                formData.append('audio_url', fileInput.files[0]);
                formData.append('type', TYPE_FILE);
            } else if (activeTab === 'youtube') {
                const youtubeUrl = youtubeInput.value.trim();
                if (!youtubeUrl) {
                    throw new Error('Введите YouTube ссылку!');
                }
                if (!isValidYouTubeURL(youtubeUrl)) {
                    throw new Error('Некорректная YouTube ссылка!');
                }
                formData.append('audio_url', youtubeUrl);
                formData.append('type', TYPE_YOUTUBE);
            }

            const response = await fetch(endpoint, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Ошибка сервера: ${response.status}`);
            }
            
            const data = await response.json();
            displayTab(data.tab || 'Табулатура успешно сгенерирована!');
            copyButton.disabled = false;
        } catch (err) {
            resultDiv.innerHTML = `<div class="error-message">Ошибка: ${err.message}</div>`;
            console.error('Error:', err);
        } finally {
            document.getElementById('loading').classList.remove('active');
            progressContainer.style.display = 'none';
        }
    });

    // Копирование табулатуры
    copyButton.addEventListener('click', () => {
        const tabText = Array.from(resultDiv.querySelectorAll('.tab-line'))
                            .map(line => line.textContent)
                            .join('\n');
        
        navigator.clipboard.writeText(tabText).then(() => {
            copyButton.classList.add('copied');
            setTimeout(() => copyButton.classList.remove('copied'), 1500);
        });
    });

    // Проверка YouTube URL
    function isValidYouTubeURL(url) {
        const pattern = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$/;
        return pattern.test(url);
    }

    // Отображение табулатуры с форматированием
    function displayTab(tabText) {
        if (!tabText) {
            resultDiv.innerHTML = '<div style="color:red; text-align:center; padding:20px;">Ошибка генерации табулатуры</div>';
            return;
        }
    
        // Очищаем контейнер
        resultDiv.innerHTML = '';
        
        // Разделяем строки и добавляем их как отдельные div
        const lines = tabText.split('\n').filter(line => line.trim());
        
        lines.forEach(line => {
            const lineDiv = document.createElement('div');
            lineDiv.className = 'tab-line';
            
            // Подсвечиваем цифры (ноты) красным
            const highlightedLine = line.replace(/(\d+)/g, '<span class="note">$1</span>');
            lineDiv.innerHTML = highlightedLine;
            
            resultDiv.appendChild(lineDiv);
        });
    }

    // Симуляция прогресса загрузки (для демонстрации)
    function simulateUploadProgress() {
        let progress = 0;
        const interval = setInterval(() => {
            progress += Math.random() * 10;
            if (progress >= 100) {
                progress = 100;
                clearInterval(interval);
            }
            progressBar.style.width = `${progress}%`;
            progressText.textContent = `${Math.round(progress)}%`;
        }, 200);
    }
});