-- name: CreateTcpRouter :one
INSERT INTO
  tcp_routers (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetTcpRouter :one
SELECT
  *
FROM
  tcp_routers
WHERE
  id = ?;

-- name: GetTcpRoutersUsingEntryPoint :many
WITH
  ep_name AS (
    SELECT
      name
    FROM
      entry_points
    WHERE
      entry_points.id = ?
      AND entry_points.profile_id = ?
  )
SELECT
  r.id,
  r.name,
  r.config,
  r.enabled
FROM
  tcp_routers r
  JOIN json_each (r.config, '$.entryPoints') je
  JOIN ep_name ep ON je.value = ep.name;

-- name: GetTcpRoutersUsingMiddleware :many
WITH
  mw_name AS (
    SELECT
      name
    FROM
      tcp_middlewares
    WHERE
      tcp_middlewares.id = ?
      AND tcp_middlewares.profile_id = ?
  )
SELECT
  r.id,
  r.name,
  r.config,
  r.enabled
FROM
  tcp_routers r
  JOIN json_each (r.config, '$.middlewares') je
  JOIN mw_name mw ON je.value = mw.name;

-- name: ListTcpRouters :many
SELECT
  *
FROM
  tcp_routers
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpRoutersByAgent :many
SELECT
  *
FROM
  tcp_routers
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListTcpRoutersEnabled :many
SELECT
  *
FROM
  tcp_routers
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountTcpRouters :one
SELECT
  COUNT(*)
FROM
  tcp_routers;

-- name: CountTcpRoutersByProfile :one
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  profile_id = ?;

-- name: CountTcpRoutersByAgent :one
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  agent_id = ?;

-- name: UpdateTcpRouter :one
UPDATE tcp_routers
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteTcpRouter :exec
DELETE FROM tcp_routers
WHERE
  id = ?;
