const BASE_URL = 'http://localhost:8080';

// Генерация таба
export async function createTask(file, separation = false) {
    const formData = new FormData();
    formData.append('audio_file', file);
    formData.append('separation', separation ? 'true' : 'false');

    const resp = await fetch(`${BASE_URL}/generation`, {
        method: 'POST',
        body: formData
    });

    if (!resp.ok) {
        const errData = await resp.json().catch(() => ({}));
        throw new Error(errData.error || `Ошибка сервера: ${resp.status}`);
    }

    return await resp.json(); // возвращает { task: TabGenTask, tab: Tab | undefined }
}

// Получение статуса задачи (polling)
export async function getTask(taskId) {
    const resp = await fetch(`${BASE_URL}/generation/${taskId}`);
    if (!resp.ok) {
        const errData = await resp.json().catch(() => ({}));
        throw new Error(errData.error || `Ошибка сервера: ${resp.status}`);
    }
    return await resp.json(); // возвращает { task, tab }
}

// Сохранение таба
export async function saveTab(taskId, name) {
    const resp = await fetch(`${BASE_URL}/generation/save/${taskId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
            name: name
        }),
    });

    if (!resp.ok) {
        const errData = await resp.json().catch(() => ({}));
        throw new Error(errData.error || 'Ошибка сохранения');
    }

    return true;
}

// Сепарация аудио
export async function startSeparation(file) {
    const formData = new FormData()
    formData.append('audio_file', file)

    const resp = await fetch(`${BASE_URL}/audio/separation`, {
        method: 'POST',
        body: formData
    })

    if (!resp.ok) {
        throw new Error(`Ошибка сервера: ${resp.status}`)
    }

    return await resp.json() 
}

export async function getSeparationTask(id) {
    const resp = await fetch(`${BASE_URL}/audio/separation/${id}`)

    if (!resp.ok) {
        throw new Error(`Ошибка получения задачи: ${resp.status}`)
    }

    return await resp.json()
}

// Поиск табов
export async function searchTabs(name) {
    const res = await fetch(`${BASE_URL}/tab/search?name=${encodeURIComponent(name)}`);
    if (!res.ok) throw new Error((await res.json()).error || 'Ошибка поиска');
    return res.json();
}

// Получение таба
export async function getTab(id) {
    const res = await fetch(`${BASE_URL}/tab/${id}`);
    if (!res.ok) throw new Error((await res.json()).error || 'Ошибка получения таба');
    return res.json();
}

// Удаление таба
export async function deleteTab(id) {
    const res = await fetch(`${BASE_URL}/tab/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error((await res.json()).error || 'Ошибка удаления таба');
}