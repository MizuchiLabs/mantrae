-- name: CreateTraefikConfig :exec
INSERT INTO
  traefik (
    profile_id,
    source,
    entrypoints,
    overview,
    version,
    config
  )
VALUES
  (?, ?, ?, ?, ?, ?);

-- name: GetTraefikConfig :one
SELECT
  *
FROM
  traefik
WHERE
  id = ?;

-- name: ListTraefikIDs :many
SELECT
  id
FROM
  traefik;

-- name: GetTraefikConfigBySource :one
SELECT
  *
FROM
  traefik
WHERE
  profile_id = ?
  AND source = ?;

-- name: UpdateTraefikConfig :exec
UPDATE traefik
SET
  entrypoints = ?,
  overview = ?,
  config = ?,
  version = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  profile_id = ?
  AND source = ?;

-- name: DeleteTraefikConfig :exec
DELETE FROM traefik
WHERE
  id = ?;
