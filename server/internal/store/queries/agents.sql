-- name: CreateAgent :one
INSERT INTO
  agents (id, profile_id, token)
VALUES
  (?, ?, ?) RETURNING *;

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
  hostname = COALESCE(sqlc.arg ('hostname'), hostname),
  public_ip = COALESCE(sqlc.arg ('publicIp'), public_ip),
  private_ip = COALESCE(sqlc.arg ('privateIp'), private_ip),
  active_ip = COALESCE(sqlc.arg ('activeIp'), active_ip),
  containers = COALESCE(sqlc.arg ('containers'), containers),
  token = COALESCE(sqlc.narg ('token'), token),
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = sqlc.arg ('id') RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agents
WHERE
  id = ?;
