const TYPE_FILE = 0;
const TYPE_YOUTUBE = 1;

document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('audioFile');
    const youtubeInput = document.getElementById('youtubeUrl');
    const resultDiv = document.getElementById('tabResult');
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');

    // Переключение между вкладками
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            // Обновляем активные кнопки
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Обновляем активные контенты
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabId}-tab`).classList.add('active');
        });
    });

    // Обработка формы
    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const activeTab = document.querySelector('.tab-button.active').getAttribute('data-tab');
        let formData = new FormData();
        const endpoint = 'http://localhost:8081/tab-generate';

        resultDiv.textContent = 'Обработка...';

        try {
            if (activeTab === 'file') {
                if (!fileInput.files.length) {
                    throw new Error('Выберите аудиофайл!');
                }
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
                formData.append('source_type', TYPE_YOUTUBE);
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
            resultDiv.textContent = data.tab || 'Табулатура успешно сгенерирована!';
        } catch (err) {
            resultDiv.textContent = 'Ошибка: ' + err.message;
            console.error('Error:', err);
        }
    });

    // Проверка YouTube URL
    function isValidYouTubeURL(url) {
        const pattern = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$/;
        return pattern.test(url);
    }
});