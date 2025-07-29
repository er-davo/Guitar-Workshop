CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE tabs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    file_path TEXT NOT NULL,
    UNIQUE(name)
);
