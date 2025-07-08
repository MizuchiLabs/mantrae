-- name: CreateHttpRouter :one
INSERT INTO
  http_routers (
    profile_id,
    agent_id,
    name,
    config,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

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
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpRoutersByAgent :many
SELECT
  *
FROM
  http_routers
WHERE
  agent_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: ListHttpRoutersEnabled :many
SELECT
  *
FROM
  http_routers
WHERE
  profile_id = ?
  AND enabled = TRUE
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountHttpRouters :one
SELECT
  COUNT(*)
FROM
  http_routers;

-- name: CountHttpRoutersByProfile :one
SELECT
  COUNT(*)
FROM
  http_routers
WHERE
  profile_id = ?;

-- name: CountHttpRoutersByAgent :one
SELECT
  COUNT(*)
FROM
  http_routers
WHERE
  agent_id = ?;

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
