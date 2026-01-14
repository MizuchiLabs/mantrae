-- +goose Up
-- +goose StatementBegin
CREATE TABLE "agents_new" (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT,
  public_ip TEXT,
  containers TEXT,
  active_ip TEXT,
  private_ip TEXT,
  token TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

-- Copy data from old table to new table
INSERT INTO
  agents_new
SELECT
  *
FROM
  agents;

-- Drop old table
DROP TABLE agents;

-- Rename new table to original name
ALTER TABLE agents_new
RENAME TO agents;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Revert back to original schema (this assumes TEXT can be converted back)
CREATE TABLE "agents_new" (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT,
  public_ip TEXT,
  containers TEXT,
  active_ip TEXT,
  private_ip TEXT,
  token TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

INSERT INTO
  agents_new
SELECT
  *
FROM
  agents;

DROP TABLE agents;

ALTER TABLE agents_new
RENAME TO agents;

-- +goose StatementEnd
