-- +goose Up
CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    username VARCHAR(255),
    password TEXT,
    tls BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS traefik_api (
    profile_id INTEGER NOT NULL,
    entrypoints TEXT,
    overview TEXT,
    external TEXT,
    internal TEXT,
    FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
    UNIQUE (profile_id)
);

CREATE TABLE IF NOT EXISTS router_dns_provider (
    profile_id INTEGER NOT NULL,
    provider_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles (id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE CASCADE,
    UNIQUE (profile_id, name)
);

CREATE TABLE IF NOT EXISTS providers (
    id INTEGER PRIMARY KEY,
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
    id INTEGER PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email VARCHAR(255),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY,
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS agents (
    id TEXT PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    hostname TEXT NOT NULL,
    public_ip TEXT,
    private_ips TEXT,
    containers TEXT,
    active_ip TEXT,
    token TEXT NOT NULL,
    last_seen TIMESTAMP,
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

DROP TABLE IF EXISTS providers;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS settings;

DROP TABLE IF EXISTS agents;

DROP TRIGGER IF EXISTS ensure_single_active_insert;

DROP TRIGGER IF EXISTS ensure_single_active_update;
