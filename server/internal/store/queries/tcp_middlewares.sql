-- name: CreateTcpMiddleware :one
INSERT INTO
  tcp_middlewares (
    id,
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
  is_default = TRUE
  AND profile_id = ?;

-- name: ListTcpMiddlewares :many
SELECT
  *
FROM
  tcp_middlewares
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

-- name: ListTcpMiddlewaresEnabled :many
SELECT
  *
FROM
  tcp_middlewares
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountTcpMiddlewares :one
SELECT
  COUNT(*)
FROM
  tcp_middlewares
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
