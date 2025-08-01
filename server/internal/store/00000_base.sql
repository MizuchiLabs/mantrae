-- +goose Up
CREATE TABLE dns_providers (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  type VARCHAR(255) NOT NULL,
  config JSON NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_default BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE "settings" (
  key VARCHAR(255) PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "agents" (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  hostname TEXT,
  public_ip TEXT,
  containers JSONB,
  active_ip TEXT,
  token TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  private_ip TEXT,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE errors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  profile_id INTEGER NOT NULL,
  category TEXT NOT NULL,
  message TEXT NOT NULL,
  details TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE entry_points (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  address TEXT,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE TABLE http_routers (
  id TEXT PRIMARY KEY,
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

CREATE TABLE tcp_routers (
  id TEXT PRIMARY KEY,
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

CREATE TABLE udp_routers (
  id TEXT PRIMARY KEY,
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

CREATE TABLE http_services (
  id TEXT PRIMARY KEY,
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

CREATE TABLE tcp_services (
  id TEXT PRIMARY KEY,
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

CREATE TABLE udp_services (
  id TEXT PRIMARY KEY,
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

CREATE TABLE http_middlewares (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE TABLE tcp_middlewares (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  agent_id TEXT,
  name TEXT NOT NULL,
  config TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  FOREIGN KEY (agent_id) REFERENCES agents (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE TABLE http_servers_transports (
  id TEXT PRIMARY KEY,
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

CREATE TABLE tcp_servers_transports (
  id TEXT PRIMARY KEY,
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

CREATE TABLE traefik_instances (
  id TEXT PRIMARY KEY,
  profile_id INTEGER NOT NULL,
  name TEXT NOT NULL,
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
  FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
  UNIQUE (profile_id, name)
);

CREATE TABLE audit_logs (
  id INTEGER PRIMARY KEY,
  profile_id INTEGER,
  user_id TEXT,
  agent_id TEXT,
  event TEXT NOT NULL,
  details TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profiles (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  token TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id) REFERENCES profiles (id) ON DELETE CASCADE
);

CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT,
  otp VARCHAR(6),
  otp_expiry TIMESTAMP,
  last_login TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE http_router_dns_providers (
  http_router_id TEXT NOT NULL,
  dns_provider_id INTEGER NOT NULL,
  PRIMARY KEY (http_router_id, dns_provider_id),
  FOREIGN KEY (http_router_id) REFERENCES http_routers (id) ON DELETE CASCADE,
  FOREIGN KEY (dns_provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE
);

CREATE TABLE tcp_router_dns_providers (
  tcp_router_id TEXT NOT NULL,
  dns_provider_id INTEGER NOT NULL,
  PRIMARY KEY (tcp_router_id, dns_provider_id),
  FOREIGN KEY (tcp_router_id) REFERENCES tcp_routers (id) ON DELETE CASCADE,
  FOREIGN KEY (dns_provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE
);

CREATE INDEX idx_http_routers_profile_name ON http_routers (profile_id, name);

CREATE INDEX idx_tcp_routers_profile_name ON tcp_routers (profile_id, name);

CREATE INDEX idx_udp_routers_profile_name ON udp_routers (profile_id, name);

CREATE INDEX idx_http_services_profile_name ON http_services (profile_id, name);

CREATE INDEX idx_tcp_services_profile_name ON tcp_services (profile_id, name);

CREATE INDEX idx_udp_services_profile_name ON udp_services (profile_id, name);

CREATE INDEX idx_http_middlewares_profile_name ON http_middlewares (profile_id, name);

CREATE INDEX idx_tcp_middlewares_profile_name ON tcp_middlewares (profile_id, name);

CREATE INDEX idx_http_servers_transports_profile_name ON http_servers_transports (profile_id, name);

CREATE INDEX idx_tcp_servers_transports_profile_name ON tcp_servers_transports (profile_id, name);

-- +goose Down
DROP TABLE IF EXISTS dns_providers;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS agents;

DROP TABLE IF EXISTS errors;

DROP TABLE IF EXISTS entry_points;

DROP TABLE IF EXISTS http_routers;

DROP TABLE IF EXISTS tcp_routers;

DROP TABLE IF EXISTS udp_routers;

DROP TABLE IF EXISTS http_services;

DROP TABLE IF EXISTS tcp_services;

DROP TABLE IF EXISTS udp_services;

DROP TABLE IF EXISTS http_middlewares;

DROP TABLE IF EXISTS tcp_middlewares;

DROP TABLE IF EXISTS http_servers_transports;

DROP TABLE IF EXISTS tcp_servers_transports;

DROP TABLE IF EXISTS traefik_instances;

DROP TABLE IF EXISTS audit_logs;

DROP TABLE IF EXISTS profiles;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS http_router_dns_providers;

DROP TABLE IF EXISTS tcp_router_dns_providers;
