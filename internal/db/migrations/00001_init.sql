-- +goose Up
CREATE TABLE IF NOT EXISTS profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(255),
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS entrypoints (
  profile_id INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  as_default BOOLEAN DEFAULT FALSE,
  http JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE TABLE IF NOT EXISTS routers (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  provider VARCHAR(255) NOT NULL,
  protocol VARCHAR(255) NOT NULL,
  status VARCHAR(255),
  agent_id TEXT,
  entry_points JSONB,
  middlewares JSONB,
  rule TEXT NOT NULL,
  rule_syntax VARCHAR(255),
  service TEXT NOT NULL,
  priority INTEGER,
  tls JSONB,
  dns_provider INTEGER,
  errors JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (dns_provider) REFERENCES providers (id) ON DELETE SET NULL,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (name, profile_id)
);

CREATE TABLE IF NOT EXISTS services (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  provider VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  protocol VARCHAR(255) NOT NULL,
  agent_id TEXT,
  status VARCHAR(255),
  server_status JSONB,
  load_balancer JSONB,
  weighted JSONB,
  mirroring JSONB,
  failover JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (name, profile_id)
);

CREATE TABLE IF NOT EXISTS middlewares (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  provider VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  protocol VARCHAR(255) NOT NULL,
  agent_id TEXT,
  content JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (name, profile_id)
);

CREATE TABLE IF NOT EXISTS providers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  type VARCHAR(255) NOT NULL,
  external_ip TEXT NOT NULL,
  api_key TEXT NOT NULL,
  api_url TEXT,
  zone_type TEXT,
  proxied BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(255),
  type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key VARCHAR(255) NOT NULL UNIQUE,
  value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS agents (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT NOT NULL,
  public_ip TEXT,
  private_ips JSONB,
  containers JSONB,
  active_ip TEXT,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  last_seen DATETIME,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

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
DROP TABLE IF EXISTS profiles;

DROP TABLE IF EXISTS entrypoints;

DROP TABLE IF EXISTS routers;

DROP TABLE IF EXISTS services;

DROP TABLE IF EXISTS middlewares;

DROP TABLE IF EXISTS providers;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS agents;

DROP TRIGGER IF EXISTS ensure_single_active_insert;

DROP TRIGGER IF EXISTS ensure_single_active_update;
