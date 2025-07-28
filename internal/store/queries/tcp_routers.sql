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

-- name: ListTcpRoutersEnabled :many
SELECT
  *
FROM
  tcp_routers
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountTcpRouters :one
SELECT
  COUNT(*)
FROM
  tcp_routers
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

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
