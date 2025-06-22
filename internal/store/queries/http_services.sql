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

-- name: GetHttpService :one
SELECT
  *
FROM
  http_services
WHERE
  id = ?;

-- name: GetHttpServiceByName :one
SELECT
  *
FROM
  http_services
WHERE
  name = ?;

-- name: ListHttpServices :many
SELECT
  *
FROM
  http_services
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpServicesByAgent :many
SELECT
  *
FROM
  http_services
WHERE
  agent_id = ?
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

-- name: CountHttpServicesByProfile :one
SELECT
  COUNT(*)
FROM
  http_services
WHERE
  profile_id = ?;

-- name: CountHttpServicesByAgent :one
SELECT
  COUNT(*)
FROM
  http_services
WHERE
  agent_id = ?;

-- name: UpdateHttpService :one
UPDATE http_services
SET
  name = ?,
  config = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpService :exec
DELETE FROM http_services
WHERE
  id = ?;
