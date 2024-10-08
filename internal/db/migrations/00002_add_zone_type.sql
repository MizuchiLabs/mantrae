-- +goose Up
ALTER TABLE providers
ADD COLUMN zone_type TEXT;

-- +goose Down
ALTER TABLE providers
DROP COLUMN zone_type;
