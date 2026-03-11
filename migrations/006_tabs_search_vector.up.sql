ALTER TABLE tabs
ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (to_tsvector('simple', name)) STORED;

CREATE INDEX idx_tabs_search_vector
ON tabs USING gin(search_vector)
WHERE deleted_at IS NULL;