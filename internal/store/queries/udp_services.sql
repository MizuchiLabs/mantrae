-- name: CreateUdpService :one
INSERT INTO
  udp_services (
    profile_id,
    name,
    config,
    source,
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

-- name: GetUdpServiceByProfile :one
SELECT
  *
FROM
  udp_services
WHERE
  profile_id = ?;

-- name: ListUdpServices :many
SELECT
  *
FROM
  udp_services
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

-- name: UpdateUdpService :one
UPDATE udp_services
SET
  name = ?,
  config = ?,
  source = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteUdpService :exec
DELETE FROM udp_services
WHERE
  id = ?;
