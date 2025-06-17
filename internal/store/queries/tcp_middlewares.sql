-- name: CreateTcpMiddleware :one
INSERT INTO
  tcp_middlewares (
    profile_id,
    name,
    config,
    source,
    enabled,
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

-- name: GetTcpMiddleware :one
SELECT
  *
FROM
  tcp_middlewares
WHERE
  id = ?;

-- name: GetTcpMiddlewareByProfile :one
SELECT
  *
FROM
  tcp_middlewares
WHERE
  profile_id = ?;

-- name: ListTcpMiddlewares :many
SELECT
  *
FROM
  tcp_middlewares
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountTcpMiddlewares :one
SELECT
  COUNT(*)
FROM
  tcp_middlewares;

-- name: UpdateTcpMiddleware :one
UPDATE tcp_middlewares
SET
  name = ?,
  config = ?,
  source = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpMiddleware :exec
DELETE FROM tcp_middlewares
WHERE
  id = ?;
