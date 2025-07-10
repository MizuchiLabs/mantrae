-- name: CreateTcpServersTransport :one
INSERT INTO
  tcp_servers_transports (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetTcpServersTransport :one
SELECT
  *
FROM
  tcp_servers_transports
WHERE
  id = ?;

-- name: ListTcpServersTransports :many
SELECT
  *
FROM
  tcp_servers_transports
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpServersTransportsByAgent :many
SELECT
  *
FROM
  tcp_servers_transports
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpServersTransportsEnabled :many
SELECT
  *
FROM
  tcp_servers_transports
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountTcpServersTransports :one
SELECT
  COUNT(*)
FROM
  tcp_servers_transports;

-- name: CountTcpServersTransportsByProfile :one
SELECT
  COUNT(*)
FROM
  tcp_servers_transports
WHERE
  profile_id = ?;

-- name: CountTcpServersTransportsByAgent :one
SELECT
  COUNT(*)
FROM
  tcp_servers_transports
WHERE
  agent_id = ?;

-- name: UpdateTcpServersTransport :one
UPDATE tcp_servers_transports
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpServersTransport :exec
DELETE FROM tcp_servers_transports
WHERE
  id = ?;
