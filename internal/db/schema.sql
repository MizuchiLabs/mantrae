-- db/schema.sql
CREATE TABLE profiles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(100) NOT NULL UNIQUE,
  url TEXT NOT NULL,
  username VARCHAR(100),
  password TEXT,
  tls BOOLEAN NOT NULL
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
  api_url TEXT
);

CREATE TABLE credentials (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(100) NOT NULL UNIQUE,
  password TEXT NOT NULL
);

-- Trigger to create an empty config when inserting a profile
CREATE TRIGGER add_config AFTER INSERT ON profiles FOR EACH ROW BEGIN
INSERT INTO
  config (profile_id)
VALUES
  (NEW.id);

END;
