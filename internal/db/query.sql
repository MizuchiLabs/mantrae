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
  name LIKE ?
LIMIT
  1;

-- name: ListProfiles :many
SELECT
  *
FROM
  profiles;

-- name: CreateProfile :one
INSERT INTO
  profiles (name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
  name = ?,
  url = ?,
  username = ?,
  password = ?,
  tls = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertProfile :one
INSERT INTO
  profiles (id, name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  name = EXCLUDED.name,
  url = EXCLUDED.url,
  username = EXCLUDED.username,
  password = EXCLUDED.password,
  tls = EXCLUDED.tls RETURNING *;

-- name: DeleteProfileByID :exec
DELETE FROM profiles
WHERE
  id = ?;

-- name: DeleteProfileByName :exec
DELETE FROM profiles
WHERE
  name = ?;

-- name: GetEntryPointByName :one
SELECT
  *
FROM
  entrypoints
WHERE
  name = ?
  AND profile_id = ?
LIMIT
  1;

-- name: ListEntryPoints :many
SELECT
  *
FROM
  entrypoints;

-- name: ListEntryPointsByProfileID :many
SELECT
  *
FROM
  entrypoints
WHERE
  profile_id = ?;

-- name: UpsertEntryPoint :one
INSERT INTO
  entrypoints (profile_id, name, address, as_default, http)
VALUES
  (?, ?, ?, ?, ?) ON CONFLICT (profile_id, name) DO
UPDATE
SET
  address = EXCLUDED.address,
  as_default = EXCLUDED.as_default,
  http = EXCLUDED.http RETURNING *;

-- name: DeleteEntryPointByName :exec
DELETE FROM entrypoints
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      profiles.name = ?
  )
  AND entrypoints.name = ?;

-- name: GetRouterByID :one
SELECT
  *
FROM
  routers
WHERE
  id = ?
LIMIT
  1;

-- name: GetRouterByName :one
SELECT
  *
FROM
  routers
WHERE
  name LIKE ?
  AND profile_id = ?
LIMIT
  1;

-- name: ListRouters :many
SELECT
  *
FROM
  routers;

-- name: ListRoutersByProfileID :many
SELECT
  *
FROM
  routers
WHERE
  profile_id = ?;

-- name: ListRoutersByProvider :many
SELECT
  *
FROM
  routers
WHERE
  provider = ?;

-- name: ListRoutersByAgentID :many
SELECT
  *
FROM
  routers
WHERE
  agent_id = ?;

-- name: UpsertRouter :one
INSERT INTO
  routers (
    id,
    profile_id,
    name,
    provider,
    protocol,
    status,
    agent_id,
    entry_points,
    middlewares,
    rule,
    rule_syntax,
    service,
    priority,
    tls,
    dns_provider,
    errors
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (profile_id, name) DO
UPDATE
SET
  id = COALESCE(NULLIF(EXCLUDED.id, ''), routers.id),
  provider = COALESCE(NULLIF(EXCLUDED.provider, ''), routers.provider),
  protocol = COALESCE(NULLIF(EXCLUDED.protocol, ''), routers.protocol),
  status = COALESCE(NULLIF(EXCLUDED.status, ''), routers.status),
  entry_points = COALESCE(
    NULLIF(EXCLUDED.entry_points, ''),
    routers.entry_points
  ),
  middlewares = COALESCE(
    NULLIF(EXCLUDED.middlewares, ''),
    routers.middlewares
  ),
  rule = COALESCE(NULLIF(EXCLUDED.rule, ''), routers.rule),
  rule_syntax = COALESCE(
    NULLIF(EXCLUDED.rule_syntax, ''),
    routers.rule_syntax
  ),
  service = COALESCE(NULLIF(EXCLUDED.service, ''), routers.service),
  priority = COALESCE(NULLIF(EXCLUDED.priority, ''), routers.priority),
  tls = COALESCE(NULLIF(EXCLUDED.tls, ''), routers.tls),
  dns_provider = CASE
    WHEN EXCLUDED.dns_provider = 0 THEN NULL
    WHEN EXCLUDED.dns_provider IS NOT NULL THEN EXCLUDED.dns_provider
    ELSE routers.dns_provider
  END,
  agent_id = COALESCE(NULLIF(EXCLUDED.agent_id, ''), routers.agent_id),
  errors = COALESCE(NULLIF(EXCLUDED.errors, ''), routers.errors) RETURNING *;

-- name: DeleteRouterByID :exec
DELETE FROM routers
WHERE
  id = ?;

-- name: DeleteRouterByName :exec
DELETE FROM routers
WHERE
  name = ?;

-- name: GetServiceByID :one
SELECT
  *
FROM
  services
WHERE
  id = ?
LIMIT
  1;

-- name: GetServiceByName :one
SELECT
  *
FROM
  services
WHERE
  name LIKE ?
  AND profile_id = ?
LIMIT
  1;

-- name: ListServices :many
SELECT
  *
FROM
  services;

-- name: ListServicesByProfileID :many
SELECT
  *
FROM
  services
WHERE
  profile_id = ?;

-- name: ListServicesByProvider :many
SELECT
  *
FROM
  services
WHERE
  provider = ?;

-- name: UpsertService :one
INSERT INTO
  services (
    id,
    profile_id,
    name,
    provider,
    type,
    protocol,
    agent_id,
    status,
    server_status,
    load_balancer,
    weighted,
    mirroring,
    failover
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (name, profile_id) DO
UPDATE
SET
  provider = COALESCE(NULLIF(EXCLUDED.provider, ''), services.provider),
  type = COALESCE(NULLIF(EXCLUDED.type, ''), services.type),
  protocol = COALESCE(NULLIF(EXCLUDED.protocol, ''), services.protocol),
  status = COALESCE(NULLIF(EXCLUDED.status, ''), services.status),
  server_status = COALESCE(
    NULLIF(EXCLUDED.server_status, ''),
    services.server_status
  ),
  load_balancer = COALESCE(
    NULLIF(EXCLUDED.load_balancer, ''),
    services.load_balancer
  ),
  weighted = COALESCE(NULLIF(EXCLUDED.weighted, ''), services.weighted),
  mirroring = COALESCE(
    NULLIF(EXCLUDED.mirroring, ''),
    services.mirroring
  ),
  failover = COALESCE(NULLIF(EXCLUDED.failover, ''), services.failover),
  agent_id = COALESCE(NULLIF(EXCLUDED.agent_id, ''), services.agent_id) RETURNING *;

-- name: DeleteServiceByID :exec
DELETE FROM services
WHERE
  id = ?;

-- name: DeleteServiceByName :exec
DELETE FROM services
WHERE
  name = ?;

-- name: GetMiddlewareByID :one
SELECT
  *
FROM
  middlewares
WHERE
  id = ?
LIMIT
  1;

-- name: GetMiddlewareByName :one
SELECT
  *
FROM
  middlewares
WHERE
  name LIKE ?
  AND profile_id = ?
LIMIT
  1;

-- name: ListMiddlewares :many
SELECT
  *
FROM
  middlewares;

-- name: ListMiddlewaresByProfileID :many
SELECT
  *
FROM
  middlewares
WHERE
  profile_id = ?;

-- name: ListMiddlewaresByProvider :many
SELECT
  *
FROM
  middlewares
WHERE
  provider = ?;

-- name: UpsertMiddleware :one
INSERT INTO
  middlewares (
    id,
    profile_id,
    name,
    provider,
    type,
    protocol,
    agent_id,
    content
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (name, profile_id) DO
UPDATE
SET
  provider = COALESCE(
    NULLIF(EXCLUDED.provider, ''),
    middlewares.provider
  ),
  type = COALESCE(NULLIF(EXCLUDED.type, ''), middlewares.type),
  protocol = COALESCE(
    NULLIF(EXCLUDED.protocol, ''),
    middlewares.protocol
  ),
  agent_id = COALESCE(
    NULLIF(EXCLUDED.agent_id, ''),
    middlewares.agent_id
  ),
  content = COALESCE(NULLIF(EXCLUDED.content, ''), middlewares.content) RETURNING *;

-- name: DeleteMiddlewareByID :exec
DELETE FROM middlewares
WHERE
  id = ?;

-- name: DeleteMiddlewareByName :exec
DELETE FROM middlewares
WHERE
  name = ?;

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

-- name: CreateProvider :one
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
  (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProvider :one
UPDATE providers
SET
  name = ?,
  type = ?,
  external_ip = ?,
  api_key = ?,
  api_url = ?,
  zone_type = ?,
  proxied = ?,
  is_active = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertProvider :one
INSERT INTO
  providers (
    id,
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
  (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  name = EXCLUDED.name,
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

-- name: DeleteProviderByName :exec
DELETE FROM providers
WHERE
  name = ?;

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

-- name: CreateUser :one
INSERT INTO
  users (username, password, email, is_admin)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  username = ?,
  password = ?,
  email = ?,
  is_admin = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertUser :one
INSERT INTO
  users (id, username, password, email, is_admin)
VALUES
  (?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  username = COALESCE(NULLIF(EXCLUDED.username, ''), users.username),
  password = COALESCE(NULLIF(EXCLUDED.password, ''), users.password),
  email = COALESCE(NULLIF(EXCLUDED.email, ''), users.email),
  is_admin = COALESCE(EXCLUDED.is_admin, users.is_admin) RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE
  id = ?;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE
  username = ?;

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

-- name: UpsertSetting :one
INSERT INTO
  settings (id, key, value)
VALUES
  (?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  key = EXCLUDED.key,
  value = EXCLUDED.value RETURNING *;

-- name: DeleteSettingByID :exec
DELETE FROM settings
WHERE
  id = ?;

-- name: DeleteSettingByKey :exec
DELETE FROM settings
WHERE
  key = ?;

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
    deleted,
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
  deleted = COALESCE(NULLIF(EXCLUDED.deleted, FALSE), agents.deleted),
  last_seen = COALESCE(NULLIF(EXCLUDED.last_seen, ''), agents.last_seen) RETURNING *;

-- name: DeleteAgentByID :exec
DELETE FROM agents
WHERE
  id = ?;

-- name: DeleteAgentByHostname :exec
DELETE FROM agents
WHERE
  hostname = ?;
