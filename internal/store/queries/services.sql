-- name: ListServicesByProfile :many
SELECT
  *,
  'http' AS type
FROM
  http_services
WHERE
  http_services.profile_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_services
WHERE
  tcp_services.profile_id = ?
UNION ALL
SELECT
  *,
  'udp' AS type
FROM
  udp_services
WHERE
  udp_services.profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListServicesByAgent :many
SELECT
  *,
  'http' AS type
FROM
  http_services
WHERE
  http_services.agent_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_services
WHERE
  tcp_services.agent_id = ?
UNION ALL
SELECT
  *,
  'udp' AS type
FROM
  udp_services
WHERE
  udp_services.agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountServicesByProfile :one
SELECT
  COUNT(*)
FROM
  http_services
WHERE
  http_services.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_services
WHERE
  tcp_services.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  udp_services
WHERE
  udp_services.profile_id = ?;

-- name: CountServicesByAgent :one
SELECT
  COUNT(*)
FROM
  http_services
WHERE
  http_services.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_services
WHERE
  tcp_services.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  udp_services
WHERE
  udp_services.agent_id = ?;
