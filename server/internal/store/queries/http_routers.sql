-- name: CreateHttpRouter :one
INSERT INTO
  http_routers (
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

-- name: GetHttpRouter :one
SELECT
  *
FROM
  http_routers
WHERE
  id = ?;

-- name: GetHttpRoutersUsingEntryPoint :many
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
  http_routers r
  JOIN json_each (r.config, '$.entryPoints') je
  JOIN ep_name ep ON je.value = ep.name;

-- name: GetHttpRoutersUsingMiddleware :many
WITH
  mw_name AS (
    SELECT
      name
    FROM
      http_middlewares
    WHERE
      http_middlewares.id = ?
      AND http_middlewares.profile_id = ?
  )
SELECT
  r.id,
  r.name,
  r.config,
  r.enabled
FROM
  http_routers r
  JOIN json_each (r.config, '$.middlewares') je
  JOIN mw_name mw ON je.value = mw.name;

-- name: ListHttpRouters :many
SELECT
  *
FROM
  http_routers
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

-- name: ListHttpRoutersEnabled :many
SELECT
  *
FROM
  http_routers
WHERE
  profile_id = ?
  AND enabled = TRUE;

-- name: CountHttpRouters :one
SELECT
  COUNT(*)
FROM
  http_routers
WHERE
  profile_id = sqlc.arg ('profile_id')
  AND (
    CAST(sqlc.narg ('agent_id') AS TEXT) IS NULL
    OR agent_id = CAST(sqlc.narg ('agent_id') AS TEXT)
  );

-- name: UpdateHttpRouter :one
UPDATE http_routers
SET
  name = ?,
  config = ?,
  enabled = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteHttpRouter :exec
DELETE FROM http_routers
WHERE
  id = ?;
