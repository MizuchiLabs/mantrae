-- name: CreateAgent :exec
INSERT INTO
  agents (id, profile_id, token)
VALUES
  (?, ?, ?);

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
ORDER BY
  hostname;

-- name: ListAgentsByProfile :many
SELECT
  *
FROM
  agents
WHERE
  profile_id = ?;

-- name: UpdateAgent :one
UPDATE agents
SET
  hostname = COALESCE(?, hostname),
  public_ip = COALESCE(?, public_ip),
  private_ips = COALESCE(?, private_ips),
  containers = COALESCE(?, containers),
  active_ip = COALESCE(?, active_ip),
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
