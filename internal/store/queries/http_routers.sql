-- name: CreateHttpRouter :one
INSERT INTO
  http_routers (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetHttpRouter :one
SELECT
  *
FROM
  http_routers
WHERE
  id = ?;

-- name: ListHttpRouters :many
SELECT
  *
FROM
  http_routers
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountHttpRouters :one
SELECT
  COUNT(*)
FROM
  http_routers;

-- name: UpdateHttpRouter :one
UPDATE http_routers
SET
  name = ?,
  config = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpRouter :exec
DELETE FROM http_routers
WHERE
  id = ?;
