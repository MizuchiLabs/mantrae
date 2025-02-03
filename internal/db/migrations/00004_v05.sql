-- +goose Up
-- +goose StatementBegin
ALTER TABLE profiles
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE profiles
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create new dns_providers table and migrate data from providers
CREATE TABLE dns_providers (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  type VARCHAR(255) NOT NULL,
  config JSON NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
  dns_providers (id, name, type, config, is_active)
SELECT
  id,
  name,
  type,
  json_object (
    'traefikIp',
    external_ip,
    'apiKey',
    api_key,
    'apiUrl',
    api_url,
    'zoneType',
    zone_type,
    'proxied',
    proxied
  ),
  is_active
FROM
  providers;

-- Create new traefik table
CREATE TABLE traefik (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  source TEXT NOT NULL,
  entrypoints JSON,
  overview JSON,
  config JSON,
  version TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  CONSTRAINT valid_source CHECK (source IN ('local', 'api', 'agent'))
);

-- Create router_dns_provider table
CREATE TABLE router_dns_provider (
  traefik_id INTEGER NOT NULL,
  provider_id INTEGER NOT NULL,
  router_name TEXT NOT NULL,
  FOREIGN KEY (traefik_id) REFERENCES traefik (id) ON DELETE CASCADE,
  FOREIGN KEY (provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE,
  UNIQUE (traefik_id, router_name)
);

-- Update users table
ALTER TABLE users
ADD COLUMN otp VARCHAR(6);

ALTER TABLE users
ADD COLUMN otp_expiry TIMESTAMP;

ALTER TABLE users
ADD COLUMN last_login TIMESTAMP;

ALTER TABLE users
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE users
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Update settings table
CREATE TABLE new_settings (
  key VARCHAR(255) PRIMARY KEY,
  value TEXT NOT NULL,
  description TEXT,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
  new_settings (key, value)
SELECT
  key,
  value
FROM
  settings;

DROP TABLE settings;

ALTER TABLE new_settings
RENAME TO settings;

-- Update agents table
ALTER TABLE agents
DROP COLUMN hostname;

ALTER TABLE agents
ADD COLUMN hostname TEXT;

ALTER TABLE agents
ADD COLUMN token TEXT NOT NULL;

ALTER TABLE agents
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE agents
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE agents
DROP COLUMN deleted;

ALTER TABLE agents
DROP COLUMN last_seen;

-- Create indexes for traefik table
CREATE UNIQUE INDEX idx_traefik_profile_source ON traefik (profile_id, source)
WHERE
  source IN ('local', 'api');

CREATE UNIQUE INDEX idx_traefik_agent ON traefik (profile_id, agent_id)
WHERE
  agent_id IS NOT NULL
  AND source = 'agent';

-- Drop old tables
DROP TABLE providers;

DROP TABLE middlewares;

DROP TABLE services;

DROP TABLE routers;

DROP TABLE entrypoints;

DROP TRIGGER IF EXISTS ensure_single_active_insert;

DROP TRIGGER IF EXISTS ensure_single_active_update;

-- +goose StatementEnd
-- Create trigger
-- +goose StatementBegin
CREATE TRIGGER ensure_single_active_insert BEFORE INSERT ON dns_providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE dns_providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER ensure_single_active_update BEFORE
UPDATE ON dns_providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE dns_providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE router_dns_provider;

DROP TABLE traefik;

DROP TABLE dns_providers;

CREATE TABLE providers (
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

ALTER TABLE profiles
DROP COLUMN created_at;

ALTER TABLE profiles
DROP COLUMN updated_at;

ALTER TABLE users
DROP COLUMN is_admin;

ALTER TABLE users
DROP COLUMN otp;

ALTER TABLE users
DROP COLUMN otp_expiry;

ALTER TABLE users
DROP COLUMN last_login;

ALTER TABLE users
DROP COLUMN created_at;

ALTER TABLE users
DROP COLUMN updated_at;

ALTER TABLE users
ADD COLUMN type VARCHAR(255) NOT NULL DEFAULT 'user';

DROP TABLE settings;

CREATE TABLE settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key VARCHAR(255) NOT NULL UNIQUE,
  value TEXT NOT NULL
);

ALTER TABLE agents
DROP COLUMN token;

ALTER TABLE agents
DROP COLUMN created_at;

ALTER TABLE agents
DROP COLUMN updated_at;

ALTER TABLE agents
ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE agents
ADD COLUMN last_seen DATETIME;

-- +goose StatementEnd
