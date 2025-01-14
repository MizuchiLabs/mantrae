-- name: GetHTTPMiddlewaresBySource :one
SELECT
    json_extract (config, '$.middlewares') as middlewares
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPMiddlewaresBySource :one
SELECT
    json_extract (config, '$.tcpMiddlewares') as middlewares
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetHTTPMiddlewareByName :one
SELECT
    json_extract (config, '$.middlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: GetTCPMiddlewareByName :one
SELECT
    json_extract (config, '$.tcpMiddlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertHTTPMiddleware :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.middlewares.' || sqlc.arg (name),
        sqlc.arg (middleware_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: UpsertTCPMiddleware :exec
UPDATE traefik_config
SET
    config = json_set (
        config,
        '$.tcpMiddlewares.' || sqlc.arg (name),
        sqlc.arg (middleware_json)
    ),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteHTTPMiddleware :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.middlewares.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;

-- name: DeleteTCPMiddleware :exec
UPDATE traefik_config
SET
    config = json_remove (config, '$.tcpMiddlewares.' || sqlc.arg (name)),
    last_sync = CURRENT_TIMESTAMP
WHERE
    profile_id = ?
    AND source = ?;
