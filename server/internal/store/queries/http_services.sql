-- name: CreateHttpService :one
INSERT INTO
  http_services (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetHttpServiceByID :one
SELECT
  *
FROM
  http_services
WHERE
  profile_id = ?
  AND id = ?;

-- name: GetHttpServiceByName :one
SELECT
  *
FROM
  http_services
WHERE
  profile_id = ?
  AND name = ?;

-- name: ListHttpServices :many
SELECT
  *
FROM
  http_services
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  )
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: ListHttpServicesEnabled :many
SELECT
  *
FROM
  http_services
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountHttpServices :one
SELECT
  COUNT(*)
FROM
  http_services
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

-- name: UpdateHttpService :one
UPDATE http_services
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpService :exec
DELETE FROM http_services
WHERE
  id = ?;
