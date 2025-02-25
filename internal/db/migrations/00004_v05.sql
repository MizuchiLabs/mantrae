-- +goose Up
-- +goose StatementBegin
-- Create new profiles table with timestamp columns
CREATE TABLE new_profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(255),
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from old profiles table
INSERT INTO
  new_profiles (id, name, url, username, password, tls)
SELECT
  id,
  name,
  url,
  username,
  password,
  tls
FROM
  profiles;

-- Drop old profiles table and rename new one
DROP TABLE profiles;

ALTER TABLE new_profiles
RENAME TO profiles;

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

-- Create new users table with additional columns
CREATE TABLE new_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(255),
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  otp VARCHAR(6),
  otp_expiry TIMESTAMP,
  last_login TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from old users table
INSERT INTO
  new_users (id, username, password, email, is_admin)
SELECT
  id,
  username,
  password,
  email,
  is_admin
FROM
  users;

-- Drop old users table and rename new one
DROP TABLE users;

ALTER TABLE new_users
RENAME TO users;

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
-- Create new agents table with modified columns
CREATE TABLE new_agents (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT,
  public_ip TEXT,
  private_ips JSONB,
  containers JSONB,
  active_ip TEXT,
  token TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

-- Copy data from old agents table
INSERT INTO
  new_agents (
    id,
    profile_id,
    public_ip,
    private_ips,
    containers,
    active_ip
  )
SELECT
  id,
  profile_id,
  public_ip,
  private_ips,
  containers,
  active_ip
FROM
  agents;

-- Drop old agents table and rename new one
DROP TABLE agents;

ALTER TABLE new_agents
RENAME TO agents;

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
-- +goose StatementBegin
-- Create trigger
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
DROP TRIGGER IF EXISTS ensure_single_active_insert;

DROP TRIGGER IF EXISTS ensure_single_active_update;

-- Recreate the original tables
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

-- Restore providers data from dns_providers
INSERT INTO
  providers (
    id,
    name,
    type,
    external_ip,
    api_key,
    api_url,
    zone_type,
    proxied,
    is_active
  )
SELECT
  id,
  name,
  type,
  json_extract (config, '$.traefikIp'),
  json_extract (config, '$.apiKey'),
  json_extract (config, '$.apiUrl'),
  json_extract (config, '$.zoneType'),
  json_extract (config, '$.proxied'),
  is_active
FROM
  dns_providers;

-- Recreate original tables
CREATE TABLE routers (
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

CREATE TABLE services (
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

CREATE TABLE middlewares (
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

CREATE TABLE entrypoints (
  profile_id INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  as_default BOOLEAN DEFAULT FALSE,
  http JSONB,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

-- Restore original profiles table
CREATE TABLE new_profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(255),
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO
  new_profiles (id, name, url, username, password, tls)
SELECT
  id,
  name,
  url,
  username,
  password,
  tls
FROM
  profiles;

DROP TABLE profiles;

ALTER TABLE new_profiles
RENAME TO profiles;

-- Restore original users table
CREATE TABLE new_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email VARCHAR(255),
  is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO
  new_users (id, username, password, email, is_admin)
SELECT
  id,
  username,
  password,
  email,
  is_admin
FROM
  users;

DROP TABLE users;

ALTER TABLE new_users
RENAME TO users;

-- Restore original settings table
CREATE TABLE new_settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  key VARCHAR(255) NOT NULL UNIQUE,
  value TEXT NOT NULL
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

-- Restore original agents table
CREATE TABLE new_agents (
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

INSERT INTO
  new_agents (
    id,
    profile_id,
    hostname,
    public_ip,
    private_ips,
    containers,
    active_ip,
    deleted,
    last_seen
  )
SELECT
  id,
  profile_id,
  hostname,
  public_ip,
  private_ips,
  containers,
  active_ip,
  false,
  NULL
FROM
  agents;

DROP TABLE agents;

ALTER TABLE new_agents
RENAME TO agents;

-- Recreate triggers
CREATE TRIGGER ensure_single_active_insert BEFORE INSERT ON providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

CREATE TRIGGER ensure_single_active_update BEFORE
UPDATE ON providers FOR EACH ROW WHEN NEW.is_active = 1 BEGIN
UPDATE providers
SET
  is_active = 0
WHERE
  is_active = 1;

END;

-- Drop new tables
DROP TABLE IF EXISTS router_dns_provider;

DROP TABLE IF EXISTS traefik;

DROP TABLE IF EXISTS dns_providers;

-- +goose StatementEnd
