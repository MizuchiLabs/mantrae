-- name: CreateTcpService :one
INSERT INTO
  tcp_services (
    profile_id,
    name,
    config,
    source,
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

-- name: GetTcpServiceByProfile :one
SELECT
  *
FROM
  tcp_services
WHERE
  profile_id = ?;

-- name: ListTcpServices :many
SELECT
  *
FROM
  tcp_services
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

-- name: UpdateTcpService :one
UPDATE tcp_services
SET
  name = ?,
  config = ?,
  source = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpService :exec
DELETE FROM tcp_services
WHERE
  id = ?;
