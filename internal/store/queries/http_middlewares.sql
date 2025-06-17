-- name: CreateHttpMiddleware :one
INSERT INTO
  http_middlewares (
    profile_id,
    name,
    config,
    source,
    enabled,
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

-- name: GetHttpMiddleware :one
SELECT
  *
FROM
  http_middlewares
WHERE
  id = ?;

-- name: GetHttpMiddlewareByProfile :one
SELECT
  *
FROM
  http_middlewares
WHERE
  profile_id = ?;

-- name: ListHttpMiddlewares :many
SELECT
  *
FROM
  http_middlewares
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountHttpMiddlewares :one
SELECT
  COUNT(*)
FROM
  http_middlewares;

-- name: UpdateHttpMiddleware :one
UPDATE http_middlewares
SET
  name = ?,
  config = ?,
  source = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpMiddleware :exec
DELETE FROM http_middlewares
WHERE
  id = ?;
