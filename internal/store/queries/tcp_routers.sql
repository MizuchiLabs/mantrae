-- name: CreateTcpRouter :one
INSERT INTO
  tcp_routers (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetTcpRouter :one
SELECT
  *
FROM
  tcp_routers
WHERE
  id = ?;

-- name: ListTcpRouters :many
SELECT
  *
FROM
  tcp_routers
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpRoutersByAgent :many
SELECT
  *
FROM
  tcp_routers
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountTcpRouters :one
SELECT
  COUNT(*)
FROM
  tcp_routers;

-- name: CountTcpRoutersByProfile :one
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  profile_id = ?;

-- name: CountTcpRoutersByAgent :one
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  agent_id = ?;

-- name: UpdateTcpRouter :one
UPDATE tcp_routers
SET
  name = ?,
  config = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpRouter :exec
DELETE FROM tcp_routers
WHERE
  id = ?;
