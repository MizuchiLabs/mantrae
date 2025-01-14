-- +goose Up
CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    username VARCHAR(255),
    password TEXT,
    tls BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS traefik_config (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    source TEXT NOT NULL,
    entrypoints JSON,
    overview JSON,
    config JSON,
    last_sync TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
    UNIQUE (profile_id, source)
);

CREATE TABLE IF NOT EXISTS dns_providers (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(255) NOT NULL,
    config JSON NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS router_dns_provider (
    traefik_id INTEGER NOT NULL,
    provider_id INTEGER NOT NULL,
    router_name TEXT NOT NULL,
    FOREIGN KEY (traefik_id) REFERENCES traefik (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE,
    UNIQUE (traefik_id, router_name)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email VARCHAR(255),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS agents (
    id TEXT PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    hostname TEXT NOT NULL,
    public_ip TEXT,
    private_ips JSON,
    containers JSON,
    active_ip TEXT,
    token TEXT NOT NULL,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE
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

DROP TABLE IF EXISTS dns_providers;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS agents;

DROP TRIGGER IF EXISTS ensure_single_active_insert;

DROP TRIGGER IF EXISTS ensure_single_active_update;
