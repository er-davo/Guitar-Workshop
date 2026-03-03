// ===== CONFIG =====
const API_BASE = 'http://localhost:8080';
const POLL_INTERVAL = 2000;
const TASK_TIMEOUT = 10 * 60 * 1000; // 10 минут

let pollingInterval = null;
let pollingStartedAt = null;
let currentTaskId = null;
let isSubmitting = false;

// ===== STATUS MAP =====
function mapStatus(status) {
    switch (status) {
        case 'created': return 'Задача создана';
        case 'pending': return 'Обработка';
        case 'waiting_for_separation': return 'Сепарация аудио';
        case 'done': return 'Готово';
        case 'fail': return 'Ошибка';
        default: return status;
    }
}

// ===== UI HELPERS =====
function setLoading(active, text = '') {
    const loadingDiv = document.getElementById('loading');
    const loadingText = loadingDiv.querySelector('.loading-text');

    if (active) {
        loadingDiv.classList.add('active');
        if (loadingText && text) loadingText.textContent = text;
    } else {
        loadingDiv.classList.remove('active');
    }
}

function showError(message) {
    const resultDiv = document.getElementById('tabResult');
    resultDiv.innerHTML = `<div class="error-message">Ошибка: ${message}</div>`;
}

function showStatus(status) {
    const resultDiv = document.getElementById('tabResult');
    resultDiv.innerHTML = `<div class="placeholder-text">${mapStatus(status)}</div>`;
}

function displayTab(tabText) {
    const resultDiv = document.getElementById('tabResult');
    resultDiv.innerHTML = '';

    if (!tabText) {
        resultDiv.innerHTML = '<div class="placeholder-text">Пустой результат</div>';
        return;
    }

    tabText.split('\n').filter(l => l.trim()).forEach(line => {
        const div = document.createElement('div');
        div.className = 'tab-line';
        div.innerHTML = line.replace(/(\d+)/g, '<span class="note">$1</span>');
        resultDiv.appendChild(div);
    });
}

// ===== API =====
async function createTask(file, separation) {
    const formData = new FormData();
    formData.append('audio_file', file);
    formData.append('separation', separation ? 'true' : 'false');

    const response = await fetch(`${API_BASE}/generation`, {
        method: 'POST',
        body: formData
    });

    if (!response.ok) {
        const err = await response.json().catch(() => ({}));
        throw new Error(err.error || `Ошибка сервера: ${response.status}`);
    }

    return response.json(); // TabGenTask
}

async function getTask(taskId) {
    const response = await fetch(`${API_BASE}/generation/${taskId}`);

    if (!response.ok) {
        throw new Error(`Ошибка получения статуса: ${response.status}`);
    }

    return response.json(); // TabGenResponse
}

async function loadTabFromURL(url) {
    const resp = await fetch(url);
    if (!resp.ok) throw new Error('Ошибка загрузки таба');
    return resp.text();
}

// ===== POLLING =====
function stopPolling() {
    if (pollingInterval) {
        clearInterval(pollingInterval);
        pollingInterval = null;
    }
}

function startPolling(taskId) {
    pollingStartedAt = Date.now();
    currentTaskId = taskId;

    pollingInterval = setInterval(async () => {
        try {
            if (Date.now() - pollingStartedAt > TASK_TIMEOUT) {
                stopPolling();
                setLoading(false);
                showError('Превышено время ожидания задачи');
                isSubmitting = false;
                return;
            }

            const data = await getTask(taskId);
            const task = data.task;

            if (!task) {
                throw new Error('Некорректный ответ сервера');
            }

            showStatus(task.status);

            if (task.status === 'done') {
                stopPolling();

                if (!data.tab || !data.tab.presigned_url) {
                    throw new Error('Результат отсутствует');
                }

                const tabText = await loadTabFromURL(data.tab.presigned_url);

                displayTab(tabText);

                document.getElementById('copyButton').disabled = false;
                document.getElementById('saveSection').style.display = 'flex';
                document.getElementById('saveTabButton').dataset.taskId = taskId;

                setLoading(false);
                isSubmitting = false;
            }

            if (task.status === 'fail') {
                stopPolling();
                setLoading(false);
                showError(task.error || 'Ошибка генерации');
                isSubmitting = false;
            }

        } catch (err) {
            stopPolling();
            setLoading(false);
            showError(err.message);
            isSubmitting = false;
        }

    }, POLL_INTERVAL);
}

// ===== SUBMIT =====
document.getElementById('uploadForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (isSubmitting) return;

    const fileInput = document.getElementById('audioFile');
    if (!fileInput.files.length) return;

    isSubmitting = true;
    stopPolling();

    document.getElementById('tabResult').innerHTML = '';
    document.getElementById('copyButton').disabled = true;
    document.getElementById('saveSection').style.display = 'none';

    setLoading(true, 'Создание задачи');

    try {
        const task = await createTask(
            fileInput.files[0],
            document.getElementById('enableSeparation')?.checked
        );

        showStatus(task.status);
        startPolling(task.id);

    } catch (err) {
        setLoading(false);
        showError(err.message);
        isSubmitting = false;
    }
});

// ===== SAVE =====
document.getElementById('saveTabButton').addEventListener('click', async () => {
    const taskId = document.getElementById('saveTabButton').dataset.taskId;
    if (!taskId) return;

    try {
        const resp = await fetch(`${API_BASE}/generation/save/${taskId}`, {
            method: 'POST'
        });

        if (!resp.ok) throw new Error('Ошибка сохранения');

        alert('Таб сохранён');
        document.getElementById('saveSection').style.display = 'none';

    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
});

// ===== CLEANUP =====
window.addEventListener('beforeunload', () => {
    stopPolling();
});