-- name: CreateUdpService :one
INSERT INTO
  udp_services (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetUdpService :one
SELECT
  *
FROM
  udp_services
WHERE
  id = ?;

-- name: GetUdpServiceByName :one
SELECT
  *
FROM
  udp_services
WHERE
  name = ?;

-- name: ListUdpServices :many
SELECT
  *
FROM
  udp_services
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListUdpServicesByAgent :many
SELECT
  *
FROM
  udp_services
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListUdpServicesEnabled :many
SELECT
  *
FROM
  udp_services
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountUdpServices :one
SELECT
  COUNT(*)
FROM
  udp_services;

-- name: CountUdpServicesByProfile :one
SELECT
  COUNT(*)
FROM
  udp_services
WHERE
  profile_id = ?;

-- name: CountUdpServicesByAgent :one
SELECT
  COUNT(*)
FROM
  udp_services
WHERE
  agent_id = ?;

-- name: UpdateUdpService :one
UPDATE udp_services
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteUdpService :exec
DELETE FROM udp_services
WHERE
  id = ?;
