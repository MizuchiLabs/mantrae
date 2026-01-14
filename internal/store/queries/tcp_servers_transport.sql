-- name: CreateTcpServersTransport :one
INSERT INTO
  tcp_servers_transports (id, profile_id, agent_id, name, config)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

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

-- name: ListTcpServersTransportsEnabled :many
SELECT
  *
FROM
  tcp_servers_transports
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountTcpServersTransports :one
SELECT
  COUNT(*)
FROM
  tcp_servers_transports
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
