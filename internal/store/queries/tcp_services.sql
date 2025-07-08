-- name: CreateTcpService :one
INSERT INTO
  tcp_services (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetTcpService :one
SELECT
  *
FROM
  tcp_services
WHERE
  id = ?;

-- name: GetTcpServiceByName :one
SELECT
  *
FROM
  tcp_services
WHERE
  name = ?;

-- name: ListTcpServices :many
SELECT
  *
FROM
  tcp_services
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpServicesByAgent :many
SELECT
  *
FROM
  tcp_services
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpServicesEnabled :many
SELECT
  *
FROM
  tcp_services
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountTcpServices :one
SELECT
  COUNT(*)
FROM
  tcp_services;

-- name: CountTcpServicesByProfile :one
SELECT
  COUNT(*)
FROM
  tcp_services
WHERE
  profile_id = ?;

-- name: CountTcpServicesByAgent :one
SELECT
  COUNT(*)
FROM
  tcp_services
WHERE
  agent_id = ?;

-- name: UpdateTcpService :one
UPDATE tcp_services
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpService :exec
DELETE FROM tcp_services
WHERE
  id = ?;
