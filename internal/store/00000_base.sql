-- +goose Up
CREATE TABLE profiles (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE entry_points (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  address TEXT NOT NULL,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE http_routers (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_http_routers_profile_name ON http_routers (profile_id, name);

CREATE TABLE tcp_routers (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_tcp_middlewares_profile_name ON tcp_middlewares (profile_id, name);

CREATE TABLE udp_routers (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_udp_routers_profile_name ON udp_routers (profile_id, name);

CREATE TABLE http_services (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_http_services_profile_name ON http_services (profile_id, name);

CREATE TABLE tcp_services (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_tcp_services_profile_name ON tcp_services (profile_id, name);

CREATE TABLE udp_services (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_udp_services_profile_name ON udp_services (profile_id, name);

CREATE TABLE http_middlewares (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_http_middlewares_profile_name ON http_middlewares (profile_id, name);

CREATE TABLE tcp_middlewares (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE INDEX idx_tcp_middlewares_profile_name ON tcp_middlewares (profile_id, name);

CREATE TABLE traefik_instances (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  entrypoints TEXT,
  overview TEXT,
  config TEXT,
  version TEXT,
  url TEXT NOT NULL,
  username TEXT,
  password TEXT,
  tls BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  otp VARCHAR(6),
  otp_expiry TIMESTAMP,
  last_login TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE errors (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  category TEXT NOT NULL,
  message TEXT NOT NULL,
  details TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE agents (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT,
  public_ip TEXT,
  private_ip TEXT,
  containers TEXT,
  active_ip TEXT,
  token TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE dns_providers (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL,
  config TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE http_router_dns_providers (
  http_router_id INTEGER NOT NULL,
  dns_provider_id INTEGER NOT NULL,
  PRIMARY KEY (http_router_id, dns_provider_id),
  FOREIGN KEY (http_router_id) REFERENCES http_routers (id) ON DELETE CASCADE,
  FOREIGN KEY (dns_provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE
);

CREATE TABLE tcp_router_dns_providers (
  tcp_router_id INTEGER NOT NULL,
  dns_provider_id INTEGER NOT NULL,
  PRIMARY KEY (tcp_router_id, dns_provider_id),
  FOREIGN KEY (tcp_router_id) REFERENCES tcp_routers (id) ON DELETE CASCADE,
  FOREIGN KEY (dns_provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE
);

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

DROP TABLE IF EXISTS entry_points;

DROP TABLE IF EXISTS http_routers;

DROP TABLE IF EXISTS tcp_routers;

DROP TABLE IF EXISTS udp_routers;

DROP TABLE IF EXISTS http_services;

DROP TABLE IF EXISTS tcp_services;

DROP TABLE IF EXISTS udp_services;

DROP TABLE IF EXISTS http_middlewares;

DROP TABLE IF EXISTS tcp_middlewares;

DROP TABLE IF EXISTS traefik_instances;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS errors;

DROP TABLE IF EXISTS agents;

DROP TABLE IF EXISTS dns_providers;

DROP TABLE IF EXISTS http_router_dns_providers;

DROP TABLE IF EXISTS tcp_router_dns_providers;

DROP TABLE IF EXISTS udp_router_dns_providers;
