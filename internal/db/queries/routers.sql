-- name: GetHTTPRoutersBySource :one
SELECT
    json_extract (config, '$.routers') as routers
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPRoutersBySource :one
SELECT
    json_extract (config, '$.tcpRouters') as routers
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetUDPRoutersBySource :one
SELECT
    json_extract (config, '$.udpRouters') as routers
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetHTTPRouterByName :one
SELECT
    json_extract (config, '$.routers.' || sqlc.arg (name)) as router
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPRouterByName :one
SELECT
    json_extract (config, '$.tcpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetUDPRouterByName :one
SELECT
    json_extract (config, '$.udpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertHTTPRouter :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.routers.' || sqlc.arg (name),
        sqlc.arg (router_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertTCPRouter :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.tcpRouters.' || sqlc.arg (name),
        sqlc.arg (router_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertUDPRouter :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.udpRouters.' || sqlc.arg (name),
        sqlc.arg (router_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteHTTPRouter :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.routers.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteTCPRouter :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.tcpRouters.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteUDPRouter :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.udpRouters.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;
