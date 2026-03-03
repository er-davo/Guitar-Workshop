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