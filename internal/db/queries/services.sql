-- name: GetHTTPServicesBySource :one
SELECT
    json_extract (config, '$.services') as services
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPServicesBySource :one
SELECT
    json_extract (config, '$.tcpServices') as services
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetUDPServicesBySource :one
SELECT
    json_extract (config, '$.udpServices') as services
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetHTTPServiceByName :one
SELECT
    json_extract (config, '$.services.' || sqlc.arg (name)) as service
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPServiceByName :one
SELECT
    json_extract (config, '$.tcpServices.' || sqlc.arg (name)) as service
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetUDPServiceByName :one
SELECT
    json_extract (config, '$.udpServices.' || sqlc.arg (name)) as service
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertHTTPService :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.services.' || sqlc.arg (name),
        sqlc.arg (service_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertTCPService :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.tcpServices.' || sqlc.arg (name),
        sqlc.arg (service_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertUDPService :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.udpServices.' || sqlc.arg (name),
        sqlc.arg (service_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteHTTPService :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.services.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteTCPService :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.tcpServices.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteUDPService :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.udpServices.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;
