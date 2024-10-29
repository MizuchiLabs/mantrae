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

-- name: GetConfigByProfileID :one
SELECT
  *
FROM
  config
WHERE
  profile_id = ?
LIMIT
  1;

-- name: GetConfigByProfileName :one
SELECT
  *
FROM
  config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  )
LIMIT
  1;

-- name: ListConfigs :many
SELECT
  *
FROM
  config;

-- name: CreateConfig :one
INSERT INTO
  config (
    profile_id,
    overview,
    entrypoints,
    routers,
    services,
    middlewares,
    tls,
    version
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateConfig :one
UPDATE config
SET
  overview = ?,
  entrypoints = ?,
  routers = ?,
  services = ?,
  middlewares = ?,
  tls = ?,
  version = ?
WHERE
  profile_id = ? RETURNING *;

-- name: DeleteConfigByProfileID :exec
DELETE FROM config
WHERE
  profile_id = ?;

-- name: DeleteConfigByProfileName :exec
DELETE FROM config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  );

-- name: GetRouterByID :one
SELECT
  *
FROM
  routers
WHERE
  id = ?
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

-- name: CreateRouter :one
INSERT INTO
  routers (
    profile_id,
    name,
    provider,
    protocol,
    status,
    entry_points,
    middlewares,
    rule,
    rule_syntax,
    service,
    priority,
    tls,
    dns_provider,
    agent_id,
    errors
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateRouter :one
UPDATE routers
SET
  name = ?,
  provider = ?,
  protocol = ?,
  status = ?,
  entry_points = ?,
  middlewares = ?,
  rule = ?,
  rule_syntax = ?,
  service = ?,
  priority = ?,
  tls = ?,
  dns_provider = ?,
  agent_id = ?,
  errors = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertRouter :one
INSERT INTO
  routers (
    id,
    profile_id,
    name,
    provider,
    protocol,
    status,
    entry_points,
    middlewares,
    rule,
    rule_syntax,
    service,
    priority,
    tls,
    dns_provider,
    agent_id,
    errors
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  profile_id = EXCLUDED.profile_id,
  name = EXCLUDED.name,
  provider = EXCLUDED.provider,
  protocol = EXCLUDED.protocol,
  status = EXCLUDED.status,
  entry_points = EXCLUDED.entry_points,
  middlewares = EXCLUDED.middlewares,
  rule = EXCLUDED.rule,
  rule_syntax = EXCLUDED.rule_syntax,
  service = EXCLUDED.service,
  priority = EXCLUDED.priority,
  tls = EXCLUDED.tls,
  dns_provider = EXCLUDED.dns_provider,
  agent_id = EXCLUDED.agent_id,
  errors = EXCLUDED.errors RETURNING *;

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

-- name: CreateService :one
INSERT INTO
  services (
    profile_id,
    name,
    provider,
    type,
    protocol,
    status,
    server_status,
    load_balancer,
    weighted,
    mirroring,
    failover,
    agent_id
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateService :one
UPDATE services
SET
  name = ?,
  provider = ?,
  protocol = ?,
  status = ?,
  server_status = ?,
  load_balancer = ?,
  weighted = ?,
  mirroring = ?,
  failover = ?,
  agent_id = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertService :one
INSERT INTO
  services (
    id,
    profile_id,
    name,
    provider,
    type,
    protocol,
    status,
    server_status,
    load_balancer,
    weighted,
    mirroring,
    failover,
    agent_id
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  profile_id = EXCLUDED.profile_id,
  name = EXCLUDED.name,
  provider = EXCLUDED.provider,
  type = EXCLUDED.type,
  protocol = EXCLUDED.protocol,
  status = EXCLUDED.status,
  server_status = EXCLUDED.server_status,
  load_balancer = EXCLUDED.load_balancer,
  weighted = EXCLUDED.weighted,
  mirroring = EXCLUDED.mirroring,
  failover = EXCLUDED.failover,
  agent_id = EXCLUDED.agent_id RETURNING *;

-- name: DeleteServiceByID :exec
DELETE FROM services
WHERE
  id = ?;

-- name: DeleteServiceByName :exec
DELETE FROM services
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

-- name: GetProviderByName :one
SELECT
  *
FROM
  providers
WHERE
  name = ?
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
  users (username, password, email, type)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  username = ?,
  password = ?,
  email = ?,
  type = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertUser :one
INSERT INTO
  users (id, username, password, email, type)
VALUES
  (?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  username = EXCLUDED.username,
  password = EXCLUDED.password,
  email = EXCLUDED.email,
  type = EXCLUDED.type RETURNING *;

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
LIMIT
  1;

-- name: ListAgents :many
SELECT
  *
FROM
  agents;

-- name: CreateAgent :one
INSERT INTO
  agents (
    id,
    hostname,
    public_ip,
    private_ips,
    containers,
    last_seen
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateAgent :one
UPDATE agents
SET
  hostname = ?,
  public_ip = ?,
  private_ips = ?,
  containers = ?,
  last_seen = ?
WHERE
  id = ? RETURNING *;

-- name: UpsertAgent :one
INSERT INTO
  agents (
    id,
    hostname,
    public_ip,
    private_ips,
    containers,
    last_seen
  )
VALUES
  (?, ?, ?, ?, ?, ?) ON CONFLICT (id) DO
UPDATE
SET
  hostname = EXCLUDED.hostname,
  public_ip = EXCLUDED.public_ip,
  private_ips = EXCLUDED.private_ips,
  containers = EXCLUDED.containers,
  last_seen = EXCLUDED.last_seen RETURNING *;

-- name: DeleteAgentByID :exec
DELETE FROM agents
WHERE
  id = ?;

-- name: DeleteAgentByHostname :exec
DELETE FROM agents
WHERE
  hostname = ?;
