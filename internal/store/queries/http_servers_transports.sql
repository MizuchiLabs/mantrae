-- name: CreateHttpServersTransport :one
INSERT INTO
  http_servers_transports (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetHttpServersTransport :one
SELECT
  *
FROM
  http_servers_transports
WHERE
  id = ?;

-- name: ListHttpServersTransports :many
SELECT
  *
FROM
  http_servers_transports
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpServersTransportsByAgent :many
SELECT
  *
FROM
  http_servers_transports
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpServersTransportsEnabled :many
SELECT
  *
FROM
  http_servers_transports
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountHttpServersTransports :one
SELECT
  COUNT(*)
FROM
  http_servers_transports;

-- name: CountHttpServersTransportsByProfile :one
SELECT
  COUNT(*)
FROM
  http_servers_transports
WHERE
  profile_id = ?;

-- name: CountHttpServersTransportsByAgent :one
SELECT
  COUNT(*)
FROM
  http_servers_transports
WHERE
  agent_id = ?;

-- name: UpdateHttpServersTransport :one
UPDATE http_servers_transports
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpServersTransport :exec
DELETE FROM http_servers_transports
WHERE
  id = ?;
