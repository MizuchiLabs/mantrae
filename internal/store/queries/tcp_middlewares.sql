-- name: CreateTcpMiddleware :one
INSERT INTO
  tcp_middlewares (
    profile_id,
    agent_id,
    name,
    config,
    is_default,
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

-- name: UnsetDefaultTcpMiddleware :exec
UPDATE tcp_middlewares
SET
  is_default = FALSE
WHERE
  is_default = TRUE;

-- name: ListTcpMiddlewares :many
SELECT
  *
FROM
  tcp_middlewares
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpMiddlewaresByAgent :many
SELECT
  *
FROM
  tcp_middlewares
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpMiddlewaresEnabled :many
SELECT
  *
FROM
  tcp_middlewares
WHERE
  profile_id = ?
  AND enabled = TRUE
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

-- name: CountTcpMiddlewaresByProfile :one
SELECT
  COUNT(*)
FROM
  tcp_middlewares
WHERE
  profile_id = ?;

-- name: CountTcpMiddlewaresByAgent :one
SELECT
  COUNT(*)
FROM
  tcp_middlewares
WHERE
  agent_id = ?;

-- name: UpdateTcpMiddleware :one
UPDATE tcp_middlewares
SET
  name = ?,
  config = ?,
  enabled = ?,
  is_default = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpMiddleware :exec
DELETE FROM tcp_middlewares
WHERE
  id = ?;
