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
  profile_id = ?
  AND name = ?;

-- name: ListTraefikInstances :many
SELECT
  *
FROM
  traefik_instances
WHERE
  profile_id = sqlc.arg ('profile_id')
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountTraefikInstances :one
SELECT
  COUNT(*)
FROM
  traefik_instances
WHERE
  profile_id = ?;

-- name: UpsertTraefikInstance :one
INSERT INTO
  traefik_instances (
    id,
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

-- name: DeleteTraefikInstance :exec
DELETE FROM traefik_instances
WHERE
  id = ?;

-- name: PurgeTraefikInstances :exec
DELETE FROM traefik_instances
WHERE
  updated_at < DATETIME ('now', '-10 minutes');
