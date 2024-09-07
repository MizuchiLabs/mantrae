-- db/schema.sql
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
  entrypoints JSONB,
  routers JSONB,
  services JSONB,
  middlewares JSONB,
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

-- Trigger to create an empty config when inserting a profile
CREATE TRIGGER add_profile_config AFTER INSERT ON profiles FOR EACH ROW BEGIN
INSERT INTO
  config (profile_id)
VALUES
  (NEW.id);

END;

-- Triggers to ensure only one provider is active at a time
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
