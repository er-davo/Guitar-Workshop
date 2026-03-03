CREATE TABLE audio_sep_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    status status NOT NULL DEFAULT 'PENDING',

    input_audio_name TEXT NOT NULL,

    separated_dir_name TEXT NOT NULL,

    error TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);