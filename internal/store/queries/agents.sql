-- name: CreateAgent :one
INSERT INTO
  agents (id, profile_id, token, created_at)
VALUES
  (?, ?, ?, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetAgent :one
SELECT
  *
FROM
  agents
WHERE
  id = ?;

-- name: ListAgents :many
SELECT
  *
FROM
  agents
WHERE
  profile_id = sqlc.arg ('profile_id')
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountAgents :one
SELECT
  COUNT(*)
FROM
  agents
WHERE
  profile_id = ?;

-- name: UpdateAgent :one
UPDATE agents
SET
  hostname = COALESCE(?, hostname),
  public_ip = COALESCE(?, public_ip),
  private_ip = COALESCE(?, private_ip),
  active_ip = COALESCE(?, active_ip),
  containers = COALESCE(?, containers),
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: UpdateAgentIP :exec
UPDATE agents
SET
  active_ip = ?
WHERE
  id = ?;

-- name: UpdateAgentToken :exec
UPDATE agents
SET
  token = ?
WHERE
  id = ?;

-- name: DeleteAgent :exec
DELETE FROM agents
WHERE
  id = ?;
