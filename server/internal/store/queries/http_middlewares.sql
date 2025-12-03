-- name: CreateHttpMiddleware :one
INSERT INTO
  http_middlewares (
    id,
    profile_id,
    agent_id,
    name,
    config,
    is_default
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetHttpMiddleware :one
SELECT
  *
FROM
  http_middlewares
WHERE
  id = ?;

-- name: UnsetDefaultHttpMiddleware :exec
UPDATE http_middlewares
SET
  is_default = FALSE
WHERE
  is_default = TRUE
  AND profile_id = ?;

-- name: ListHttpMiddlewares :many
SELECT
  *
FROM
  http_middlewares
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

-- name: ListHttpMiddlewaresEnabled :many
SELECT
  *
FROM
  http_middlewares
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountHttpMiddlewares :one
SELECT
  COUNT(*)
FROM
  http_middlewares
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

-- name: UpdateHttpMiddleware :one
UPDATE http_middlewares
SET
  name = ?,
  config = ?,
  enabled = ?,
  is_default = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpMiddleware :exec
DELETE FROM http_middlewares
WHERE
  id = ?;
