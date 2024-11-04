-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS config;

CREATE TABLE IF NOT EXISTS entrypoints (
  profile_id INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  as_default BOOLEAN DEFAULT FALSE,
  http JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS entrypoints;

-- +goose StatementEnd
