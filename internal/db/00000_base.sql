-- +goose Up
CREATE TABLE profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(255),
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
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

CREATE TABLE agents (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  public_ip TEXT,
  private_ips JSONB,
  containers JSONB,
  active_ip TEXT,
  hostname TEXT,
  token TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE dns_providers (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  type VARCHAR(255) NOT NULL,
  config JSON NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE "settings" (
  key VARCHAR(255) PRIMARY KEY,
  value TEXT NOT NULL,
  description TEXT,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "router_dns_provider" (
  traefik_id INTEGER NOT NULL,
  provider_id INTEGER NOT NULL,
  router_name TEXT NOT NULL,
  FOREIGN KEY (traefik_id) REFERENCES traefik (id) ON DELETE CASCADE,
  FOREIGN KEY (provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE,
  UNIQUE (traefik_id, router_name, provider_id)
);

CREATE UNIQUE INDEX idx_traefik_profile_source ON traefik (profile_id, source)
WHERE
  source IN ('local', 'api');

CREATE UNIQUE INDEX idx_traefik_agent ON traefik (profile_id, agent_id)
WHERE
  agent_id IS NOT NULL
  AND source = 'agent';

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
DROP TABLE IF EXISTS profiles;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS agents;

DROP TABLE IF EXISTS dns_providers;

DROP TABLE IF EXISTS traefik;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS router_dns_provider;
