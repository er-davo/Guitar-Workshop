const TYPE_FILE = 0;
const TYPE_YOUTUBE = 1;

document.addEventListener('DOMContentLoaded', () => {
    // Элементы для таб-генератора
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
    const loadingDiv = document.getElementById('loading');
    const saveSection = document.getElementById('saveSection');
    const tabNameInput = document.getElementById('tabNameInput');
    const saveTabButton = document.getElementById('saveTabButton');

    // Элементы для разделения аудио
    const separationForm = document.getElementById('separationForm');
    const separationFileInput = document.getElementById('separationFile');
    const separationFileNameDiv = document.getElementById('separationFileName');
    const separationResultDiv = document.getElementById('separationResult');
    const separationLoading = document.getElementById('separationLoading');

    setLanguage('ru');

    document.getElementById("langRu").addEventListener("click", () => {
        setLanguage("ru");
        document.getElementById("langRu").classList.add("active");
        document.getElementById("langEn").classList.remove("active");
    });

    document.getElementById("langEn").addEventListener("click", () => {
        setLanguage("en");
        document.getElementById("langEn").classList.add("active");
        document.getElementById("langRu").classList.remove("active");
    });

    // Переключение вкладок
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabId}-tab`).classList.add('active');
        });
    });

    // Обработка выбора файла (таб-генератор)
    fileInput.addEventListener('change', (e) => {
        const fileName = e.target.files[0]?.name || '';
        document.getElementById('fileName').textContent = fileName;
    });

    // Drag and drop для таб-генератора
    const fileUpload = document.querySelector('#tab-generator-tab .file-upload');
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

    // Обработка выбора файла (разделение аудио)
    separationFileInput.addEventListener('change', (e) => {
        const fileName = e.target.files[0]?.name || '';
        separationFileNameDiv.textContent = fileName;
    });

    // Drag and drop для разделения аудио
    const separationFileUpload = document.querySelector('#audio-separation-tab .file-upload');
    separationFileUpload.addEventListener('dragover', (e) => {
        e.preventDefault();
        separationFileUpload.classList.add('dragover');
    });
    separationFileUpload.addEventListener('dragleave', () => {
        separationFileUpload.classList.remove('dragover');
    });
    separationFileUpload.addEventListener('drop', (e) => {
        e.preventDefault();
        separationFileUpload.classList.remove('dragover');
        if (e.dataTransfer.files.length) {
            separationFileInput.files = e.dataTransfer.files;
            separationFileNameDiv.textContent = e.dataTransfer.files[0].name;
        }
    });

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const activeTab = document.querySelector('.tab-button.active').getAttribute('data-tab');
        if (activeTab !== 'tab-generator') return;

        let formData = new FormData();
        const endpoint = 'http://localhost:8080/tab/generate';

        loadingDiv.classList.add('active');
        resultDiv.innerHTML = 
            '<div class="tab-line">e|---------------------------------------------------------</div>' +
            '<div class="tab-line">B|---------------------------------------------------------</div>' +
            '<div class="tab-line">G|---------------------------------------------------------</div>' +
            '<div class="tab-line">D|---------------------------------------------------------</div>' +
            '<div class="tab-line">A|---------------------------------------------------------</div>' +
            '<div class="tab-line">E|---------------------------------------------------------</div>';

        try {
            if (fileInput.files.length) {
                // Показываем прогресс-бар
                progressContainer.style.display = 'block';
                progressBar.style.width = '0%';
                progressText.textContent = '0%';

                simulateUploadProgress();

                formData.append('audio_file', fileInput.files[0]);
                formData.append('type', TYPE_FILE);
            } else {
                // Иначе берём YouTube ссылку
                const youtubeUrl = youtubeInput.value.trim();
                if (!youtubeUrl) throw new Error('Введите YouTube ссылку!');
                if (!isValidYouTubeURL(youtubeUrl)) throw new Error('Некорректная YouTube ссылка!');

                formData.append('audio_url', youtubeUrl);
                formData.append('type', TYPE_YOUTUBE);
            }

            const separationEnabled = document.getElementById('enableSeparation')?.checked;
            formData.append('separation', separationEnabled ? '1' : '0');

            const response = await fetch(endpoint, {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Ошибка сервера: ${response.status}`);
            }

            const data = await response.json();
            displayTab(data.tab);
            copyButton.disabled = false;
            saveSection.style.display = 'block';
            saveTabButton.disabled = true;
        } catch (err) {
            resultDiv.innerHTML = `<div class="error-message">Ошибка: ${err.message}</div>`;
            console.error(err);
        } finally {
            loadingDiv.classList.remove('active');
            progressContainer.style.display = 'none';
        }
    });

    tabNameInput.addEventListener('input', () => {
        saveTabButton.disabled = !tabNameInput.value.trim();
    });

    saveTabButton.addEventListener('click', async () => {
        const tabName = tabNameInput.value.trim();
        if (!tabName) {
            alert(currentLang === 'ru' ? 'Введите имя для табулатуры' : 'Enter a name for the tab');
            return;
        }

        const tabText = Array.from(resultDiv.querySelectorAll('.tab-line'))
            .map(line => line.textContent)
            .join('\n');

        try {
            const resp = await fetch('/tab/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: tabName, body: tabText }),
            });

            if (!resp.ok) {
                const errData = await resp.json();
                throw new Error(errData.error || (currentLang === 'ru' ? 'Ошибка сохранения' : 'Save error'));
            }

            const data = await resp.json();
            const viewUrl = `/tab/${data.id}`;
            alert(currentLang === 'ru'
                ? `Табулатура успешно сохранена!\nПосмотреть: ${viewUrl}`
                : `Tab saved successfully!\nView: ${viewUrl}`
            );

            saveTabButton.disabled = true;
            tabNameInput.value = '';
            saveSection.style.display = 'none';
        } catch (err) {
            alert((currentLang === 'ru' ? 'Ошибка: ' : 'Error: ') + err.message);
        }
    });


    separationForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (!separationFileInput.files.length) {
            alert('Пожалуйста, выберите аудиофайл для разделения.');
            return;
        }

        // Собираем выбранные дорожки
        const selectedStems = Array.from(separationForm.querySelectorAll('input[name="stems"]:checked'))
            .map(cb => cb.value);

        if (selectedStems.length === 0) {
            alert('Выберите хотя бы одну дорожку для сохранения.');
            return;
        }

        const formData = new FormData();
        formData.append('audio_file', separationFileInput.files[0]);
        formData.append('stems', selectedStems.join(','));

        separationLoading.style.display = 'flex';
        separationResultDiv.innerHTML = '';

        try {
            const response = await fetch('http://localhost:8080/audio/separate', {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || `Ошибка сервера: ${response.status}`);
            }

            const data = await response.json();

            separationResultDiv.innerHTML = '';

            for (const [stem, url] of Object.entries(data.stems)) {
                const div = document.createElement('div');
                div.className = 'stem-item';

                const blobUrl = base64ToBlobUrl(url);

                div.innerHTML = `
                    <h4>${stem.charAt(0).toUpperCase() + stem.slice(1)}</h4>
                    <audio controls src="${blobUrl}"></audio>
                    <br/>
                    <a href="${blobUrl}" download="${stem}.wav">Скачать ${stem}</a>
                `;
                separationResultDiv.appendChild(div);
            }
        } catch (err) {
            separationResultDiv.innerHTML = `<div class="error-message">Ошибка: ${err.message}</div>`;
            console.error(err);
        } finally {
            separationLoading.style.display = 'none';
        }
    });

    function base64ToBlobUrl(base64String, mimeType = 'audio/wav') {
        const base64 = base64String.split(',')[1];
        const byteCharacters = atob(base64);
        const byteNumbers = new Array(byteCharacters.length);
        for(let i = 0; i < byteCharacters.length; i++) {
            byteNumbers[i] = byteCharacters.charCodeAt(i);
        }
        const byteArray = new Uint8Array(byteNumbers);
        const blob = new Blob([byteArray], { type: mimeType });
        return URL.createObjectURL(blob);
    }

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

    // Отображение табулатуры с подсветкой нот
    function displayTab(tabText) {
        if (!tabText) {
            resultDiv.innerHTML = '<div style="color:red; text-align:center; padding:20px;">Ошибка генерации табулатуры</div>';
            return;
        }

        resultDiv.innerHTML = '';

        const lines = tabText.split('\n').filter(line => line.trim());

        lines.forEach(line => {
            const lineDiv = document.createElement('div');
            lineDiv.className = 'tab-line';

            const highlightedLine = line.replace(/(\d+)/g, '<span class="note">$1</span>');
            lineDiv.innerHTML = highlightedLine;

            resultDiv.appendChild(lineDiv);
        });
    }

    // Симуляция прогресса загрузки
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


function setLanguage(lang) {
    const dict = translations[lang];
    if (!dict) return;

    document.querySelectorAll("[data-i18n]").forEach(el => {
        const key = el.getAttribute("data-i18n");
        if (dict[key]) {
            el.textContent = dict[key];
        }
    });

    document.querySelectorAll("[data-i18n-placeholder]").forEach(el => {
        const key = el.getAttribute("data-i18n-placeholder");
        if (dict[key]) {
            el.placeholder = dict[key];
        }
    });


    const ytInput = document.getElementById("youtubeUrl");
    if (ytInput) ytInput.placeholder = dict.youtubePlaceholder;

    const loadingText = document.querySelector("#loading .loading-text");
    if (loadingText) loadingText.textContent = dict.loading;

    const separationLoadingText = document.querySelector("#separationLoading .loading-text");
    if (separationLoadingText) separationLoadingText.textContent = dict.loadingSeparation;
}

async function searchTabs() {
    const name = document.getElementById('searchInput').value;
    if (!name) return;

    const response = await fetch(`/tab/search?name=${encodeURIComponent(name)}`);
    const results = await response.json();

    const resultsList = document.getElementById('searchResults');
    resultsList.innerHTML = '';

    results.forEach(tab => {
        const li = document.createElement('li');
        const link = document.createElement('a');
        link.href = `/tab/view/${tab.id}`;
        link.innerText = tab.name;
        link.target = '_blank';

        li.appendChild(link);
        resultsList.appendChild(li);
    });
}

document.getElementById('searchInput').addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        searchTabs();
    }
});
