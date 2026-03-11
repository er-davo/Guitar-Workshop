DROP INDEX idx_tabs_search_vector;

ALTER TABLE tabs
DELETE COLUMN search_vector;