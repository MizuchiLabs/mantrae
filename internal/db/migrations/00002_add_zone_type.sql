-- +goose Up
ALTER TABLE providers
ADD COLUMN IF NOT EXISTS zone_type TEXT;

-- +goose Down
ALTER TABLE providers
DROP COLUMN IF EXISTS zone_type;
