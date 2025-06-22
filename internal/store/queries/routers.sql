-- name: ListRoutersByProfile :many
SELECT
  *,
  'http' AS type
FROM
  http_routers
WHERE
  http_routers.profile_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_routers
WHERE
  tcp_routers.profile_id = ?
UNION ALL
SELECT
  *,
  'udp' AS type
FROM
  udp_routers
WHERE
  udp_routers.profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListRoutersByAgent :many
SELECT
  *,
  'http' AS type
FROM
  http_routers
WHERE
  http_routers.agent_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_routers
WHERE
  tcp_routers.agent_id = ?
UNION ALL
SELECT
  *,
  'udp' AS type
FROM
  udp_routers
WHERE
  udp_routers.agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountRoutersByProfile :one
SELECT
  COUNT(*)
FROM
  http_routers
WHERE
  http_routers.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  tcp_routers.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  udp_routers
WHERE
  udp_routers.profile_id = ?;

-- name: CountRoutersByAgent :one
SELECT
  COUNT(*)
FROM
  http_routers
WHERE
  http_routers.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  tcp_routers.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  udp_routers
WHERE
  udp_routers.agent_id = ?;
