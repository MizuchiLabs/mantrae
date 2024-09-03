-- name: GetProfileByID :one
SELECT
  *
FROM
  profiles
WHERE
  id = ?
LIMIT
  1;

-- name: GetProfileByName :one
SELECT
  *
FROM
  profiles
WHERE
  name = ?
LIMIT
  1;

-- name: ListProfiles :many
SELECT
  *
FROM
  profiles;

-- name: CreateProfile :one
INSERT INTO
  profiles (name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
  name = ?,
  url = ?,
  username = ?,
  password = ?,
  tls = ?
WHERE
  id = ? RETURNING *;

-- name: DeleteProfileByID :exec
DELETE FROM profiles
WHERE
  id = ?;

-- name: DeleteProfileByName :exec
DELETE FROM profiles
WHERE
  name = ?;

-- name: GetConfigByProfileID :one
SELECT
  *
FROM
  config
WHERE
  profile_id = ?
LIMIT
  1;

-- name: GetConfigByProfileName :one
SELECT
  *
FROM
  config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  )
LIMIT
  1;

-- name: ListConfigs :many
SELECT
  *
FROM
  config;

-- name: CreateConfig :one
INSERT INTO
  config (
    profile_id,
    entrypoints,
    routers,
    services,
    middlewares,
    version
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateConfig :one
UPDATE config
SET
  entrypoints = ?,
  routers = ?,
  services = ?,
  middlewares = ?,
  version = ?
WHERE
  profile_id = ? RETURNING *;

-- name: DeleteConfigByProfileID :exec
DELETE FROM config
WHERE
  profile_id = ?;

-- name: DeleteConfigByProfileName :exec
DELETE FROM config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  );

-- name: GetProviderByID :one
SELECT
  *
FROM
  providers
WHERE
  id = ?
LIMIT
  1;

-- name: GetProviderByName :one
SELECT
  *
FROM
  providers
WHERE
  name = ?
LIMIT
  1;

-- name: ListProviders :many
SELECT
  *
FROM
  providers;

-- name: CreateProvider :one
INSERT INTO
  providers (name, type, external_ip, api_key, api_url)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProvider :one
UPDATE providers
SET
  name = ?,
  type = ?,
  external_ip = ?,
  api_key = ?,
  api_url = ?
WHERE
  id = ? RETURNING *;

-- name: DeleteProviderByID :exec
DELETE FROM providers
WHERE
  id = ?;

-- name: DeleteProviderByName :exec
DELETE FROM providers
WHERE
  name = ?;

-- name: GetCredentialByID :one
SELECT
  *
FROM
  credentials
WHERE
  id = ?
LIMIT
  1;

-- name: GetCredentialByUsername :one
SELECT
  *
FROM
  credentials
WHERE
  username = ?
LIMIT
  1;

-- name: ListCredentials :many
SELECT
  *
FROM
  credentials;

-- name: CreateCredential :exec
INSERT INTO
  credentials (username, password)
VALUES
  (?, ?);

-- name: UpdateCredential :exec
UPDATE credentials
SET
  username = ?,
  password = ?
WHERE
  id = ?;

-- name: DeleteCredentialByID :exec
DELETE FROM credentials
WHERE
  id = ?;

-- name: DeleteCredentialByUsername :exec
DELETE FROM credentials
WHERE
  username = ?;

-- name: ValidateAuth :one
SELECT
  id,
  username
FROM
  credentials
WHERE
  username = ?
  AND password = ?;
