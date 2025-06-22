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
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteUdpRouter :exec
DELETE FROM udp_routers
WHERE
  id = ?;
