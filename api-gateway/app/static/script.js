document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('audioFile');
    const resultDiv = document.getElementById('tabResult');
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');
    const copyButton = document.getElementById('copyButton');
    const loadingDiv = document.getElementById('loading');
    const saveSection = document.getElementById('saveSection');
    const tabNameInput = document.getElementById('tabNameInput');
    const saveTabButton = document.getElementById('saveTabButton');

    const separationForm = document.getElementById('separationForm');
    const separationFileInput = document.getElementById('separationFile');
    const separationFileNameDiv = document.getElementById('separationFileName');
    const separationResultDiv = document.getElementById('separationResult');
    const separationLoading = document.getElementById('separationLoading');

    setLanguage('ru');
    document.getElementById("langRu").addEventListener("click", () => setLanguage("ru"));
    document.getElementById("langEn").addEventListener("click", () => setLanguage("en"));

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.dataset.tab;
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabId}-tab`).classList.add('active');
        });
    });

    // Drag&drop для таб-генератора
    const fileUpload = document.querySelector('#tab-generator-tab .file-upload');
    ['dragover', 'dragleave', 'drop'].forEach(evt => {
        fileUpload.addEventListener(evt, (e) => {
            e.preventDefault();
            fileUpload.classList.remove('dragover');
            if (evt === 'dragover') fileUpload.classList.add('dragover');
            if (evt === 'drop' && e.dataTransfer.files.length) {
                fileInput.files = e.dataTransfer.files;
                document.getElementById('fileName').textContent = e.dataTransfer.files[0].name;
            }
        });
    });

    fileInput.addEventListener('change', (e) => {
        document.getElementById('fileName').textContent = e.target.files[0]?.name || '';
    });

    // Генерация таба
    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        if (!fileInput.files.length) return alert(translations[currentLang].errorNoFile);

        const formData = new FormData();
        formData.append('audio_file', fileInput.files[0]);
        formData.append('separation', document.getElementById('enableSeparation')?.checked ? 'true' : 'false');

        resultDiv.innerHTML = '';
        loadingDiv.classList.add('active');
        copyButton.disabled = true;
        saveSection.style.display = 'none';

        try {
            const response = await fetch('http://localhost:8080/generation', { method: 'POST', body: formData });
            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || `Ошибка сервера: ${response.status}`);
            }

            const data = await response.json();
            displayTab(data.result || '');
            copyButton.disabled = false;
            saveSection.style.display = 'flex';
            saveTabButton.disabled = true;

            saveTabButton.dataset.taskId = data.id;

        } catch (err) {
            resultDiv.innerHTML = `<div class="error-message">Ошибка: ${err.message}</div>`;
            console.error(err);
        } finally {
            loadingDiv.classList.remove('active');
        }
    });

    // Сохранение таба
    tabNameInput.addEventListener('input', () => {
        saveTabButton.disabled = !tabNameInput.value.trim();
    });

    saveTabButton.addEventListener('click', async () => {
        const name = tabNameInput.value.trim();
        if (!name) return alert(translations[currentLang].tabNamePlaceholder);

        const tabText = Array.from(resultDiv.querySelectorAll('.tab-line')).map(l => l.textContent).join('\n');
        const taskId = saveTabButton.dataset.taskId;
        if (!taskId) return alert('Нет ID задачи');

        try {
            const resp = await fetch(`http://localhost:8080/generation/save/${taskId}`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({name, body: tabText})
            });
            if (!resp.ok) throw new Error('Ошибка сохранения');
            alert('Таб успешно сохранён!');
            saveTabButton.disabled = true;
            tabNameInput.value = '';
            saveSection.style.display = 'none';
        } catch (err) {
            alert('Ошибка: ' + err.message);
        }
    });

    // Копирование таба
    copyButton.addEventListener('click', () => {
        const tabText = Array.from(resultDiv.querySelectorAll('.tab-line')).map(l => l.textContent).join('\n');
        navigator.clipboard.writeText(tabText).then(() => {
            copyButton.classList.add('copied');
            setTimeout(() => copyButton.classList.remove('copied'), 1500);
        });
    });

    // Drag&drop и выбор файла для сепарации
    separationFileInput.addEventListener('change', (e) => separationFileNameDiv.textContent = e.target.files[0]?.name || '');
    const separationFileUpload = document.querySelector('#audio-separation-tab .file-upload');
    ['dragover', 'dragleave', 'drop'].forEach(evt => {
        separationFileUpload.addEventListener(evt, (e) => {
            e.preventDefault();
            separationFileUpload.classList.remove('dragover');
            if (evt === 'dragover') separationFileUpload.classList.add('dragover');
            if (evt === 'drop' && e.dataTransfer.files.length) {
                separationFileInput.files = e.dataTransfer.files;
                separationFileNameDiv.textContent = e.dataTransfer.files[0].name;
            }
        });
    });

    // Сепарация аудио
    separationForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        if (!separationFileInput.files.length) return alert(translations[currentLang].errorNoFile);

        const selectedStems = Array.from(separationForm.querySelectorAll('input[name="stems"]:checked')).map(cb => cb.value);
        if (!selectedStems.length) return alert(translations[currentLang].errorNoStems);

        const formData = new FormData();
        formData.append('audio_file', separationFileInput.files[0]);
        formData.append('stems', selectedStems.join(','));

        separationLoading.style.display = 'flex';
        separationResultDiv.innerHTML = '';
        try {
            const resp = await fetch('http://localhost:8080/audio/separation', { method: 'POST', body: formData });
            if (!resp.ok) throw new Error('Ошибка сервера');
            const data = await resp.json();

            separationResultDiv.innerHTML = '';
            for (const [stem, url] of Object.entries(data.stems)) {
                const blobUrl = base64ToBlobUrl(url);
                const div = document.createElement('div');
                div.className = 'stem-item';
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
        } finally {
            separationLoading.style.display = 'none';
        }
    });

    function base64ToBlobUrl(base64String, mimeType = 'audio/wav') {
        const base64 = base64String.split(',')[1];
        const byteCharacters = atob(base64);
        const byteNumbers = new Array(byteCharacters.length);
        for(let i = 0; i < byteCharacters.length; i++) byteNumbers[i] = byteCharacters.charCodeAt(i);
        return URL.createObjectURL(new Blob([new Uint8Array(byteNumbers)], { type: mimeType }));
    }

    function displayTab(tabText) {
        resultDiv.innerHTML = '';
        if (!tabText) return resultDiv.innerHTML = '<div class="placeholder-text">Ошибка генерации табулатуры</div>';
        tabText.split('\n').filter(l => l.trim()).forEach(line => {
            const div = document.createElement('div');
            div.className = 'tab-line';
            div.innerHTML = line.replace(/(\d+)/g, '<span class="note">$1</span>');
            resultDiv.appendChild(div);
        });
    }

    function setLanguage(lang) {
        const dict = translations[lang];
        if (!dict) return;

        document.querySelectorAll("[data-i18n]").forEach(el => {
            const key = el.getAttribute("data-i18n");
            if (dict[key]) el.textContent = dict[key];
        });

        document.querySelectorAll("[data-i18n-placeholder]").forEach(el => {
            const key = el.getAttribute("data-i18n-placeholder");
            if (dict[key]) el.placeholder = dict[key];
        });

        const ytInput = document.getElementById("youtubeUrl");
        if (ytInput) ytInput.placeholder = dict.youtubePlaceholder;

        const loadingText = document.querySelector("#loading .loading-text");
        if (loadingText) loadingText.textContent = dict.loading;

        const separationLoadingText = document.querySelector("#separationLoading .loading-text");
        if (separationLoadingText) separationLoadingText.textContent = dict.loadingSeparation;
    }
});