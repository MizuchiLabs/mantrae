-- name: CreateHttpMiddleware :one
INSERT INTO
  http_middlewares (
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

-- name: UnsetDefaultHttpMiddleware :exec
UPDATE http_middlewares
SET
  is_default = FALSE
WHERE
  is_default = TRUE;

-- name: ListHttpMiddlewares :many
SELECT
  *
FROM
  http_middlewares
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpMiddlewaresByAgent :many
SELECT
  *
FROM
  http_middlewares
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpMiddlewaresEnabled :many
SELECT
  *
FROM
  http_middlewares
WHERE
  profile_id = ?
  AND enabled = TRUE
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

-- name: CountHttpMiddlewaresByProfile :one
SELECT
  COUNT(*)
FROM
  http_middlewares
WHERE
  profile_id = ?;

-- name: CountHttpMiddlewaresByAgent :one
SELECT
  COUNT(*)
FROM
  http_middlewares
WHERE
  agent_id = ?;

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
