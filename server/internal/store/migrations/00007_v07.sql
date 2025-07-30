-- +goose Up
-- Create new tables that don't exist in the old schema
CREATE TABLE entry_points (
  id INTEGER PRIMARY KEY,
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

CREATE TABLE http_services (
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

CREATE TABLE tcp_services (
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

CREATE TABLE udp_services (
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

CREATE TABLE http_middlewares (
  id INTEGER PRIMARY KEY,
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
  id INTEGER PRIMARY KEY,
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

CREATE TABLE tcp_servers_transports (
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

CREATE TABLE traefik_instances (
  id INTEGER PRIMARY KEY,
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

-- Update profiles table structure
CREATE TABLE new_profiles (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  token TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (id) REFERENCES profiles (id) ON DELETE CASCADE
);

INSERT INTO
  new_profiles (id, name, token)
SELECT
  id,
  name,
  ''
FROM
  profiles;

DROP TABLE profiles;

ALTER TABLE new_profiles
RENAME TO profiles;

-- Update agents table - change JSON columns to TEXT
UPDATE agents
SET
  containers = CASE
    WHEN containers IS NOT NULL THEN containers
    ELSE NULL
  END;

ALTER TABLE agents
DROP COLUMN private_ips;

ALTER TABLE agents
ADD COLUMN private_ip TEXT;

-- Update dns_providers table 
-- Drop old triggers for DNS providers 
DROP TRIGGER ensure_single_active_insert;

DROP TRIGGER ensure_single_active_update;

ALTER TABLE dns_providers
DROP COLUMN is_active;

ALTER TABLE dns_providers
ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE;

-- New users table
DROP TABLE users;

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

-- Update settings table - remove description column
ALTER TABLE settings
DROP COLUMN description;

-- Create new DNS provider association tables
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

-- Create indexes
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

-- Drop old tables and constraints
DROP TABLE router_dns_provider;

DROP TABLE traefik;

-- Remove old unique index on errors table
DROP INDEX unique_dns_error;
