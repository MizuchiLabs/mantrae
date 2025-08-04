-- name: CreateUdpService :one
INSERT INTO
  udp_services (
    id,
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ) RETURNING *;

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
  profile_id = ?
  AND name = ?;

-- name: ListUdpServices :many
SELECT
  *
FROM
  udp_services
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  )
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: ListUdpServicesEnabled :many
SELECT
  *
FROM
  udp_services
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountUdpServices :one
SELECT
  COUNT(*)
FROM
  udp_services
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
