
-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP WITH TIME ZONE null;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;