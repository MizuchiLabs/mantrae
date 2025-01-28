-- name: GetTraefikConfigByID :one
SELECT
  *
FROM
  traefik
WHERE
  id = ?;

-- name: GetLocalTraefikConfig :one
SELECT
  *
FROM
  traefik
WHERE
  profile_id = ?
  AND source = 'local';

-- name: GetAPITraefikConfig :one
SELECT
  *
FROM
  traefik
WHERE
  profile_id = ?
  AND source = 'api';

-- name: GetAgentTraefikConfigs :many
SELECT
  *
FROM
  traefik
WHERE
  profile_id = ?
  AND source = 'agent';

-- name: ListTraefikIDs :many
SELECT
  id
FROM
  traefik;

-- name: GetTraefikConfigBySource :many
SELECT
  *
FROM
  traefik
WHERE
  profile_id = ?
  AND source = ?;

-- name: UpsertTraefikConfig :exec
INSERT INTO
  traefik (
    profile_id,
    source,
    entrypoints,
    overview,
    config,
    version
  )
VALUES
  (?, ?, ?, ?, ?, ?) ON CONFLICT (profile_id, source)
WHERE
  source IN ('local', 'api') DO
UPDATE
SET
  entrypoints = excluded.entrypoints,
  overview = excluded.overview,
  config = excluded.config,
  version = excluded.version,
  updated_at = CURRENT_TIMESTAMP;

-- name: UpsertTraefikAgentConfig :exec
INSERT INTO
  traefik (profile_id, agent_id, source, config)
VALUES
  (?, ?, 'agent', ?) ON CONFLICT (profile_id, agent_id)
WHERE
  agent_id IS NOT NULL
  AND source = 'agent' DO
UPDATE
SET
  config = excluded.config,
  updated_at = CURRENT_TIMESTAMP;

-- name: DeleteTraefikConfig :exec
DELETE FROM traefik
WHERE
  id = ?;

-- name: DeleteTraefikConfigByAgent :exec
DELETE FROM traefik
WHERE
  agent_id = ?;
