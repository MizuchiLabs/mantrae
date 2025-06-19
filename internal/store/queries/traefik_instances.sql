-- name: CreateTraefikInstance :one
INSERT INTO
  traefik_instances (
    profile_id,
    url,
    username,
    password,
    tls,
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
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ) RETURNING *;

-- name: GetTraefikInstance :one
SELECT
  *
FROM
  traefik_instances
WHERE
  id = ?;

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

-- name: UpdateTraefikInstance :one
UPDATE traefik_instances
SET
  url = ?,
  username = ?,
  password = ?,
  tls = ?,
  entrypoints = ?,
  overview = ?,
  config = ?,
  version = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTraefikInstance :exec
DELETE FROM traefik_instances
WHERE
  id = ?;
