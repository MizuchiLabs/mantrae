-- name: CreateUdpRouter :one
INSERT INTO
  udp_routers (
    profile_id,
    name,
    config,
    source,
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

-- name: GetUdpRouterByProfile :one
SELECT
  *
FROM
  udp_routers
WHERE
  profile_id = ?;

-- name: ListUdpRouters :many
SELECT
  *
FROM
  udp_routers
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

-- name: UpdateUdpRouter :one
UPDATE udp_routers
SET
  name = ?,
  config = ?,
  source = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteUdpRouter :exec
DELETE FROM udp_routers
WHERE
  id = ?;
