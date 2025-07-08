-- name: CreateUdpRouter :one
INSERT INTO
  udp_routers (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetUdpRouter :one
SELECT
  *
FROM
  udp_routers
WHERE
  id = ?;

-- name: GetUdpRoutersUsingEntryPoint :many
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
  udp_routers r
  JOIN json_each (r.config, '$.entryPoints') je
  JOIN ep_name ep ON je.value = ep.name;

-- name: ListUdpRouters :many
SELECT
  *
FROM
  udp_routers
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListUdpRoutersByAgent :many
SELECT
  *
FROM
  udp_routers
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListUdpRoutersEnabled :many
SELECT
  *
FROM
  udp_routers
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountUdpRouters :one
SELECT
  COUNT(*)
FROM
  udp_routers;

-- name: CountUdpRoutersByProfile :one
SELECT
  COUNT(*)
FROM
  udp_routers
WHERE
  profile_id = ?;

-- name: CountUdpRoutersByAgent :one
SELECT
  COUNT(*)
FROM
  udp_routers
WHERE
  agent_id = ?;

-- name: UpdateUdpRouter :one
UPDATE udp_routers
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteUdpRouter :exec
DELETE FROM udp_routers
WHERE
  id = ?;
