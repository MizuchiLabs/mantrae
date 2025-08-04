-- name: CreateTcpService :one
INSERT INTO
  tcp_services (
    id,
    profile_id,
    agent_id,
    name,
    config,
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
  profile_id = ?
  AND name = ?;

-- name: ListTcpServices :many
SELECT
  *
FROM
  tcp_services
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

-- name: ListTcpServicesEnabled :many
SELECT
  *
FROM
  tcp_services
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountTcpServices :one
SELECT
  COUNT(*)
FROM
  tcp_services
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
