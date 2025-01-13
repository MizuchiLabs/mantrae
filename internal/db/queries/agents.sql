-- name: CreateAgent :exec
INSERT INTO
    agents (
        id,
        profile_id,
        hostname,
        public_ip,
        private_ips,
        containers,
        active_ip,
        token
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?);

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

-- name: UpdateAgent :exec
UPDATE agents
SET
    hostname = ?,
    public_ip = ?,
    private_ips = ?,
    containers = ?,
    active_ip = ?,
    last_seen = CURRENT_TIMESTAMP
WHERE
    id = ?;

-- name: UpdateAgentLastSeen :exec
UPDATE agents
SET
    last_seen = CURRENT_TIMESTAMP
WHERE
    id = ?;

-- name: DeleteAgent :exec
DELETE FROM agents
WHERE
    id = ?;
