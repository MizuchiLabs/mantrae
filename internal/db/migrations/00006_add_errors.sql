-- +goose Up
-- +goose StatementBegin
CREATE TABLE errors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  profile_id INTEGER NOT NULL,
  category TEXT NOT NULL,
  message TEXT NOT NULL,
  details TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_dns_error ON errors (profile_id, category, details);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE errors;

DROP INDEX unique_dns_error;

-- +goose StatementEnd
