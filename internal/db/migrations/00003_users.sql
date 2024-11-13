-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN type;

ALTER TABLE users
ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_admin;

ALTER TABLE users
ADD COLUMN type VARCHAR(255) NOT NULL;

-- +goose StatementEnd
