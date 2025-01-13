-- name: GetProfileByID :one
SELECT
    *
FROM
    profiles
WHERE
    id = ?
LIMIT
    1;

-- name: GetProfileByName :one
SELECT
    *
FROM
    profiles
WHERE
    name = ?
LIMIT
    1;

-- name: ListProfiles :many
SELECT
    *
FROM
    profiles;

-- name: UpsertProfile :one
INSERT INTO
    profiles (name, url, username, password, tls)
VALUES
    (?, ?, ?, ?, ?) ON CONFLICT (name) DO
UPDATE
SET
    url = EXCLUDED.url,
    username = EXCLUDED.username,
    password = EXCLUDED.password,
    tls = EXCLUDED.tls RETURNING *;

-- name: DeleteProfileByID :exec
DELETE FROM profiles
WHERE
    id = ?;

-- name: GetTraefikConfigByProfileID :one
SELECT
    *
FROM
    traefik_api
WHERE
    profile_id = ?
LIMIT
    1;

-- name: UpsertTraefikConfig :one
INSERT INTO
    traefik_api (
        profile_id,
        entrypoints,
        overview,
        external,
        internal
    )
VALUES
    (?, ?, ?, ?, ?) ON CONFLICT (profile_id) DO
UPDATE
SET
    entrypoints = excluded.entrypoints,
    overview = excluded.overview,
    external = excluded.external,
    internal = excluded.internal RETURNING *;

-- name: GetTraefikEntrypointsByProfileID :one
SELECT
    entrypoints
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetTraefikOverviewByProfileID :one
SELECT
    overview
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTraefikConfigByProfileID :one
SELECT
    external
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTraefikConfigByProfileID :one
SELECT
    internal
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPRoutersByProfileID :one
SELECT
    json_extract (external, '$.routers') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPRoutersByProfileID :one
SELECT
    json_extract (external, '$.tcpRouters') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalUDPRoutersByProfileID :one
SELECT
    json_extract (external, '$.udpRouters') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPRouterByName :one
SELECT
    json_extract (external, '$.routers.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPRouterByName :one
SELECT
    json_extract (external, '$.tcpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalUDPRouterByName :one
SELECT
    json_extract (external, '$.udpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPServicesByProfileID :one
SELECT
    json_extract (external, '$.services') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPServicesByProfileID :one
SELECT
    json_extract (external, '$.tcpServices') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalUDPServicesByProfileID :one
SELECT
    json_extract (external, '$.udpServices') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPServiceByName :one
SELECT
    json_extract (external, '$.services.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPServiceByName :one
SELECT
    json_extract (external, '$.tcpServices.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalUDPServiceByName :one
SELECT
    json_extract (external, '$.udpServices.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPMiddlewaresByProfileID :one
SELECT
    json_extract (external, '$.middlewares') as middlewares
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPMiddlewaresByProfileID :one
SELECT
    json_extract (external, '$.tcpMiddlewares') as middlewares
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalHTTPMiddlewareByName :one
SELECT
    json_extract (external, '$.middlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetExternalTCPMiddlewareByName :one
SELECT
    json_extract (external, '$.tcpMiddlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalHTTPRoutersByProfileID :one
SELECT
    json_extract (internal, '$.routers') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPRoutersByProfileID :one
SELECT
    json_extract (internal, '$.tcpRouters') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalUDPRoutersByProfileID :one
SELECT
    json_extract (internal, '$.udpRouters') as routers
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalHTTPRouterByName :one
SELECT
    json_extract (internal, '$.routers.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPRouterByName :one
SELECT
    json_extract (internal, '$.tcpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalUDPRouterByName :one
SELECT
    json_extract (internal, '$.udpRouters.' || sqlc.arg (name)) as router
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: UpsertInternalHTTPRouter :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.routers') IS NULL THEN json_set (internal, '$.routers', json ('{}'))
            ELSE internal
        END,
        '$.routers.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: UpsertInternalTCPRouter :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.tcpRouters') IS NULL THEN json_set (internal, '$.tcpRouters', json ('{}'))
            ELSE internal
        END,
        '$.tcpRouters.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: UpsertInternalUDPRouter :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.udpRouters') IS NULL THEN json_set (internal, '$.udpRouters', json ('{}'))
            ELSE internal
        END,
        '$.udpRouters.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: DeleteInternalHTTPRouter :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.routers.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: DeleteInternalTCPRouter :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.tcpRouters.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: DeleteInternalUDPRouter :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.udpRouters.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: GetInternalHTTPServicesByProfileID :one
SELECT
    json_extract (internal, '$.services') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPServicesByProfileID :one
SELECT
    json_extract (internal, '$.tcpServices') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalUDPServicesByProfileID :one
SELECT
    json_extract (internal, '$.udpServices') as services
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalHTTPServiceByName :one
SELECT
    json_extract (internal, '$.services.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPServiceByName :one
SELECT
    json_extract (internal, '$.tcpServices.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalUDPServiceByName :one
SELECT
    json_extract (internal, '$.udpServices.' || sqlc.arg (name)) as service
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: UpsertInternalHTTPService :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.services') IS NULL THEN json_set (internal, '$.services', json ('{}'))
            ELSE internal
        END,
        '$.services.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: UpsertInternalTCPService :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.tcpServices') IS NULL THEN json_set (internal, '$.tcpServices', json ('{}'))
            ELSE internal
        END,
        '$.tcpServices.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: UpsertInternalUDPService :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.udpServices') IS NULL THEN json_set (internal, '$.udpServices', json ('{}'))
            ELSE internal
        END,
        '$.udpServices.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: DeleteInternalHTTPService :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.services.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: DeleteInternalTCPService :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.tcpServices.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: DeleteInternalUDPService :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.udpServices.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: GetInternalHTTPMiddlewaresByProfileID :one
SELECT
    json_extract (internal, '$.middlewares') as middlewares
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPMiddlewaresByProfileID :one
SELECT
    json_extract (internal, '$.tcpMiddlewares') as middlewares
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalHTTPMiddlewareByName :one
SELECT
    json_extract (internal, '$.middlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: GetInternalTCPMiddlewareByName :one
SELECT
    json_extract (internal, '$.tcpMiddlewares.' || sqlc.arg (name)) as middleware
FROM
    traefik_api
WHERE
    profile_id = ?;

-- name: UpsertInternalHTTPMiddleware :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.middlewares') IS NULL THEN json_set (internal, '$.middlewares', json ('{}'))
            ELSE internal
        END,
        '$.middlewares.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: UpsertInternalTCPMiddleware :exec
UPDATE traefik_api
SET
    internal = json_set (
        CASE
            WHEN json_extract (internal, '$.tcpMiddlewares') IS NULL THEN json_set (internal, '$.tcpMiddlewares', json ('{}'))
            ELSE internal
        END,
        '$.tcpMiddlewares.' || sqlc.arg (name),
        json (sqlc.arg (body))
    )
WHERE
    profile_id = ?;

-- name: DeleteInternalHTTPMiddleware :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.middlewares.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: DeleteInternalTCPMiddleware :exec
UPDATE traefik_api
SET
    internal = json_remove (internal, '$.tcpMiddlewares.' || sqlc.arg (name))
WHERE
    profile_id = ?;

-- name: GetProviderByID :one
SELECT
    *
FROM
    providers
WHERE
    id = ?
LIMIT
    1;

-- name: GetDefaultProvider :one
SELECT
    *
FROM
    providers
WHERE
    is_active = true
LIMIT
    1;

-- name: ListProviders :many
SELECT
    *
FROM
    providers;

-- name: UpsertProvider :one
INSERT INTO
    providers (
        name,
        type,
        external_ip,
        api_key,
        api_url,
        zone_type,
        proxied,
        is_active
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (name) DO
UPDATE
SET
    type = EXCLUDED.type,
    external_ip = EXCLUDED.external_ip,
    api_key = EXCLUDED.api_key,
    api_url = EXCLUDED.api_url,
    zone_type = EXCLUDED.zone_type,
    proxied = EXCLUDED.proxied,
    is_active = EXCLUDED.is_active RETURNING *;

-- name: DeleteProviderByID :exec
DELETE FROM providers
WHERE
    id = ?;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = ?
LIMIT
    1;

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = ?
LIMIT
    1;

-- name: ListUsers :many
SELECT
    *
FROM
    users;

-- name: UpsertUser :one
INSERT INTO
    users (username, password, email, is_admin)
VALUES
    (?, ?, ?, ?) ON CONFLICT (username) DO
UPDATE
SET
    password = CASE
        WHEN EXCLUDED.password != '' THEN EXCLUDED.password
        ELSE users.password
    END,
    email = EXCLUDED.email,
    is_admin = EXCLUDED.is_admin RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE
    id = ?;

-- name: GetSettingByKey :one
SELECT
    *
FROM
    settings
WHERE
    key = ?
LIMIT
    1;

-- name: ListSettings :many
SELECT
    *
FROM
    settings;

-- name: CreateSetting :one
INSERT INTO
    settings (key, value)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateSetting :one
UPDATE settings
SET
    value = ?
WHERE
    key = ? RETURNING *;

-- name: DeleteSettingByID :exec
DELETE FROM settings
WHERE
    id = ?;

-- name: GetAgentByID :one
SELECT
    *
FROM
    agents
WHERE
    id = ?
LIMIT
    1;

-- name: GetAgentByHostname :one
SELECT
    *
FROM
    agents
WHERE
    hostname = ?
    AND profile_id = ?
LIMIT
    1;

-- name: ListAgents :many
SELECT
    *
FROM
    agents;

-- name: ListAgentsByProfileID :many
SELECT
    *
FROM
    agents
WHERE
    profile_id = ?;

-- name: UpsertAgent :one
INSERT INTO
    agents (
        id,
        profile_id,
        hostname,
        public_ip,
        private_ips,
        containers,
        active_ip,
        token,
        last_seen
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
    profile_id = COALESCE(NULLIF(EXCLUDED.profile_id, 0), agents.profile_id),
    hostname = COALESCE(NULLIF(EXCLUDED.hostname, ''), agents.hostname),
    public_ip = COALESCE(NULLIF(EXCLUDED.public_ip, ''), agents.public_ip),
    private_ips = COALESCE(
        NULLIF(EXCLUDED.private_ips, ''),
        agents.private_ips
    ),
    containers = COALESCE(
        NULLIF(EXCLUDED.containers, ''),
        agents.containers
    ),
    active_ip = COALESCE(NULLIF(EXCLUDED.active_ip, ''), agents.active_ip),
    token = COALESCE(NULLIF(EXCLUDED.token, ''), agents.token),
    last_seen = COALESCE(NULLIF(EXCLUDED.last_seen, ''), agents.last_seen) RETURNING *;

-- name: DeleteAgentByID :exec
DELETE FROM agents
WHERE
    id = ?;
