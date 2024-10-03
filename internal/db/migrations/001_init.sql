-- +goose Up
CREATE TABLE profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(100) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(100),
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE config (
  profile_id INTEGER NOT NULL,
  overview JSONB,
  entrypoints JSONB,
  routers JSONB,
  services JSONB,
  middlewares JSONB,
  tls JSONB,
  version TEXT,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE providers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(100) NOT NULL UNIQUE,
  type VARCHAR(50) NOT NULL,
  external_ip TEXT NOT NULL,
  api_key TEXT NOT NULL,
  api_url TEXT,
  proxied BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(100) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(100) UNIQUE,
  type VARCHAR(50) NOT NULL
);

CREATE TABLE settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key VARCHAR(100) NOT NULL UNIQUE,
  value TEXT NOT NULL
);

CREATE TABLE agents (
  id TEXT PRIMARY KEY,
  hostname VARCHAR(100) NOT NULL,
  public_ip TEXT,
  private_ips JSONB,
  containers JSONB,
  last_seen DATETIME
);

-- +goose StatementBegin
CREATE TRIGGER add_profile_config AFTER INSERT ON profiles FOR EACH ROW BEGIN
INSERT INTO
  config (profile_id)
VALUES
  (NEW.id);

END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER ensure_single_active_insert BEFORE INSERT ON providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER ensure_single_active_update BEFORE
UPDATE ON providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

-- +goose StatementEnd
-- +goose Down
DROP TABLE profiles;

DROP TABLE config;

DROP TABLE providers;

DROP TABLE users;

DROP TABLE settings;

DROP TABLE agents;

DROP TABLE containers;

DROP TRIGGER add_profile_config;

DROP TRIGGER ensure_single_active_insert;

DROP TRIGGER ensure_single_active_update;
