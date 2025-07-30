-- name: ListServersTransportsByProfile :many
SELECT
  *,
  'http' AS type
FROM
  http_servers_transports
WHERE
  http_servers_transports.profile_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_servers_transports
WHERE
  tcp_servers_transports.profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListServersTransportsByAgent :many
SELECT
  *,
  'http' AS type
FROM
  http_servers_transports
WHERE
  http_servers_transports.agent_id = ?
UNION ALL
SELECT
  *,
  'tcp' AS type
FROM
  tcp_servers_transports
WHERE
  tcp_servers_transports.agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountServersTransportsByProfile :one
SELECT
  COUNT(*)
FROM
  http_servers_transports
WHERE
  http_servers_transports.profile_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_servers_transports
WHERE
  tcp_servers_transports.profile_id = ?;

-- name: CountServersTransportsByAgent :one
SELECT
  COUNT(*)
FROM
  http_servers_transports
WHERE
  http_servers_transports.agent_id = ?
UNION ALL
SELECT
  COUNT(*)
FROM
  tcp_servers_transports
WHERE
  tcp_servers_transports.agent_id = ?;
