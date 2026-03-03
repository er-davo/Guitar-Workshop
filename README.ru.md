# 🎸 Guitar Workshop

![CI](https://github.com/er-davo/guitar-workshop/actions/workflows/go.yml/badge.svg)

**Guitar Workshop** — микросервисное приложение на Go для генерации гитарных табулатур из аудио.

Проект включает:
- Асинхронную обработку аудио
- ML-анализ нот (Basic Pitch)
- Генерацию играбельных табов с аппликатурами

Производительность (по умолчанию):
- До 2 задач сепарации аудио одновременно
- До 5 задач генерации табов одновременно (worker pool + semaphore)

Благодаря Kafka и асинхронной архитектуре приложение может обрабатывать больше задач, просто потребуется больше времени.  
Лёгкое масштабирование через увеличение числа воркеров/семафоров или развертывание нескольких реплик сервисов (например, в Kubernetes) позволяет повысить параллельность обработки.

## Возможности
- Загрузка аудио (WAV/MP3)
- Опциональная изоляция гитары (source separation)
- Извлечение нот с помощью ML (Basic Pitch)
- Генерация играбельных табулатур с аппликатурами
- Веб-интерфейс для отображения табов
- Сохранение и поиск табов в PostgreSQL
- REST API для управления табами и задачами
- Полностью асинхронный event-driven pipeline через Kafka

## Архитектура

Проект построен на микросервисной архитектуре с event-driven pipeline. Компоненты общаются через `Kafka` и `gRPC`. `REST API` предоставляется через `Echo`.

| Компонент         | Язык | Назначение                                        |
| ----------------- | -------- | -------------------------------------------------- |
| `api-gateway`     | Go       | REST API, создание задач, просмотр статуса         |
| `orchestrator`    | Go       | Координация генерации табов и аудио-сепарации      |
| `tab-generator`   | Go       | Генерация табов, worker pool, gRPC к note-analyzer |
| `audio-separator` | Python   | Сепарация аудио с Spleeter, concurrency=2          |
| `note-analyzer`   | Python   | Извлечение нот (Basic Pitch)                       |
| `PostgreSQL`      | SQL      | Сохранение метаданных табов                        |
| `S3`              | —        | Хранение файлов и presigned URLs                   |

### Мини описание микросервисов
- **api-gateway** — REST API; создаёт задачи генерации табов и сепарации; публикует события в Kafka; возвращает presigned URL после завершения задач.
- **orchestrator** — читает события из Kafka; решает, когда запускать генерацию табов: если есть аудио-сепарация — меняет статус задачи на `waiting for separation`, иначе сразу публикует `tab.generation.start`.
- **tab-generator** — потребляет tab.generation.start; использует worker pool для генерации табов и semaphore для gRPC вызовов note-analyzer; после обработки сохраняет табы в S3 и обновляет статус.
- **audio-separator** — потребляет audio.separation.start; делает сепарацию аудио через Spleeter с semaphore; публикует `audio.separation.complete`; сохраняет stems в S3.
- note-analyzer — извлекает ноты из аудио (ML модель Basic Pitch) по запросу tab-generator.
- **S3** — хранение аудио и табов; presigned URLs используются для безопасного доступа.
- **PostgreSQL** — хранение метаданных и статусов задач и табов.



---


## REST API (Echo)
### Tab Management `/tab`
| Метод | Эндпоинт            | Описание              |
| ------ | ------------------- | ------------------------ |
| GET    | `/tab/search?name=` | Поиск сохранённых табов  |
| GET    | `/tab/:id`          | Получить сохранённый таб |
| DELETE | `/tab/:id`          | Удалить таб              |

### Tab Generation `/generation`

| Метод | Эндпоинт            | Описание              |
| ------ | ------------------- | ------------------------ |
| POST   | `/generation`          | Создать задачу генерации таба |
| GET    | `/generation/:id`      | Получить статус задачи        |
| POST   | `/generation/save/:id` | Сохранить готовый таб         |

POST /generation behavior:

- Создаётся TabGenTask

- Если separation=true:

    - создаётся AudioSepTask

    - публикуется audio.separation.start

- Иначе сразу публикуется tab.generation.request

GET /generation/:id example:
```json
{
  "task": {
    "id": "ucrhu9743uj",
    "status": "waiting_for_separation",
  }
}
```
Если задача donw:
```json
{
  "task": {
    "id": "ucrhu9743uj",
    "status": "done",
  },
  "tab": {
    "presigned_url": "some url"
  }
}
```

### Audio Separation `/audio/separation`
| Метод | Эндпоинт            | Описание              |
| ------ | ------------------- | ------------------------ |
| POST   | `/audio/separation`     | Создать задачу сепарации |
| GET    | `/audio/separation/:id` | Получить задачу   |

GET /audio/separation/:id example:
```json
{
  "id": "sep1234",
  "status": "done",
  "separated_aduio_signed_urls": {
      "vocal": "https://s3.example.com/presigned/guitar_sep1234.wav",
      "drum": "https://s3.example.com/presigned/accompaniment_sep1234.wav",
      "bass": "https://s3.example.com/presigned/bass_sep1234.wav",
      "other": "https://s3.example.com/presigned/other_sep1234.wav"
  }
}
```
---
## Async Flow

### Генерация таба без сепарации
```
POST /generation
   ↓ tab.generation.request (kafka)
orchestrator
   ↓ tab.generation.start (kafka)
tab-generator (
    worker pool 5,
    semaphore 2 → gRPC note-analyzer
)
   ↓ S3 (tab saved)
```

### Генерация с сепарацией
```
POST /generation (if separation == true)
   ↓ AudioSepTask (db) → audio.separation.start (kafka)
audio-separator (semaphore=2)
   ↓ audio.separation.complete (kafka)
orchestrator
   ↓ tab.generation.start (kafka)
tab-generator
   ↓ S3 (tab saved)
```
#### Отдельная аудио-сепарация
```
POST /audio/separation
   ↓ audio.separation.start (kafka)
audio-separator
   ↓ audio.separation.complete (kafka)
```
Все сервисы работают через Kafka + presigned URLs, без блокирующих HTTP операций.

---

## Схема базы данных

### Табы
```sql
CREATE TABLE tabs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP,
    UNIQUE(name)
);


```
### Статусы
```sql
CREATE TYPE status AS ENUM(
    'CREATED',
    'PENDING',
    'WAITING_FOR_SEPARATION',
    'PROCESSING',
    'DONE',
    'FAILED'
);
```
### Задачи сепарации
```sql
CREATE TABLE audio_sep_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    status status NOT NULL DEFAULT 'PENDING',

    input_audio_name TEXT NOT NULL,

    separated_dir_name TEXT NOT NULL,

    error TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);
```
### Задачи генерации табов
```sql
CREATE TABLE tab_gen_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    audio_sep_task_id UUID REFERENCES audio_sep_tasks(id),

    status status NOT NULL DEFAULT 'PENDING',

    input_audio_name TEXT NOT NULL,

    result_tab_name TEXT,

    saved BOOLEAN NOT NULL DEFAULT FALSE,

    error TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);
```

## Примеры

### Пример генерации табов для *Maneskin - Coriline*.
Если сыграть это на гитаре, то звучание будет похожим.
![Пример табулатуры](./docs/example-tab.png)

### Поиск табулатур и их просмотр
![Searching tab](./docs/example-searching-saved-tabs-and-viewing.png)

### Разделение аудио
![Пример разделения](./docs/example-separation-part-1.png)

![Пример разделения](./docs/example-separation-part-2.png)

# Алгоритм генерации табов
Проект использует библиотеку на go [guitar](https://github.com/er-davo/guitar), написанную мной
1. Каждая нота получает возможные позиции на грифе (`Playable`)
2. Используется динамическое программирование для поиска оптимальной последовательности
3. Алгоритм минимизирует стоимость переходов между позициями (учитываются: перемещения по грифу, смена струн, растяжка пальцев)