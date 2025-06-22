-- name: ListMiddlewaresByProfile :many
SELECT
  *,
  'http' AS type
FROM
  http_middlewares
WHERE
  http_middlewares.profile_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_middlewares
WHERE
  tcp_middlewares.profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListMiddlewaresByAgent :many
SELECT
  *,
  'http' AS type
FROM
  http_middlewares
WHERE
  http_middlewares.agent_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_middlewares
WHERE
  tcp_middlewares.agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountMiddlewaresByProfile :one
SELECT
  COUNT(*)
FROM
  http_middlewares
WHERE
  http_middlewares.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_middlewares
WHERE
  tcp_middlewares.profile_id = ?;

-- name: CountMiddlewaresByAgent :one
SELECT
  COUNT(*)
FROM
  http_middlewares
WHERE
  http_middlewares.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_middlewares
WHERE
  tcp_middlewares.agent_id = ?;
