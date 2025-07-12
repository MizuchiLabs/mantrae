-- name: GetTraefikInstanceByID :one
SELECT
  *
FROM
  traefik_instances
WHERE
  id = ?;

-- name: GetTraefikInstanceByName :one
SELECT
  *
FROM
  traefik_instances
WHERE
  name = ?;

-- name: ListTraefikInstances :many
SELECT
  *
FROM
  traefik_instances
WHERE
  profile_id = ?
ORDER BY
  created_at DESC
LIMIT
  ?
OFFSET
  ?;

-- name: CountTraefikInstances :one
SELECT
  COUNT(*)
FROM
  traefik_instances;

-- name: UpsertTraefikInstance :one
INSERT INTO
  traefik_instances (
    profile_id,
    name,
    url,
    username,
    password,
    tls,
    entrypoints,
    overview,
    config,
    version,
    created_at,
    updated_at
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ) ON CONFLICT (profile_id, name) DO
UPDATE
SET
  url = EXCLUDED.url,
  username = EXCLUDED.username,
  password = EXCLUDED.password,
  tls = EXCLUDED.tls,
  entrypoints = EXCLUDED.entrypoints,
  overview = EXCLUDED.overview,
  config = EXCLUDED.config,
  version = EXCLUDED.version,
  updated_at = CURRENT_TIMESTAMP RETURNING *;

-- name: DeleteTraefikInstanceByID :exec
DELETE FROM traefik_instances
WHERE
  id = ?;

-- name: PurgeTraefikInstances :exec
DELETE FROM traefik_instances
WHERE
  updated_at < DATETIME ('now', '-10 minutes');
