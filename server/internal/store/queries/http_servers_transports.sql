-- name: CreateHttpServersTransport :one
INSERT INTO
  http_servers_transports (
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

-- name: ListHttpServersTransportsEnabled :many
SELECT
  *
FROM
  http_servers_transports
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountHttpServersTransports :one
SELECT
  COUNT(*)
FROM
  http_servers_transports
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
