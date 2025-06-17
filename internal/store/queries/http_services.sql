-- name: CreateHttpService :one
INSERT INTO
  http_services (
    profile_id,
    name,
    config,
    source,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetHttpService :one
SELECT
  *
FROM
  http_services
WHERE
  id = ?;

-- name: GetHttpServiceByProfile :one
SELECT
  *
FROM
  http_services
WHERE
  profile_id = ?;

-- name: ListHttpServices :many
SELECT
  *
FROM
  http_services
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountHttpServices :one
SELECT
  COUNT(*)
FROM
  http_services;

-- name: UpdateHttpService :one
UPDATE http_services
SET
  name = ?,
  config = ?,
  source = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpService :exec
DELETE FROM http_services
WHERE
  id = ?;
