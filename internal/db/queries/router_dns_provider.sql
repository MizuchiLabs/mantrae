-- name: GetRouterDNSProviders :many
SELECT
  rdp.*,
  dp.name as provider_name,
  dp.type as provider_type
FROM
  router_dns_provider rdp
  JOIN dns_providers dp ON dp.id = rdp.provider_id
WHERE
  rdp.traefik_id = ?
  AND rdp.router_name = ?;

-- name: GetRouterDNSProviderByID :one
SELECT
  rdp.*,
  dp.name as provider_name,
  dp.type as provider_type
FROM
  router_dns_provider rdp
  JOIN dns_providers dp ON dp.id = rdp.provider_id
WHERE
  rdp.traefik_id = ?
  AND rdp.router_name = ?
  AND rdp.provider_id = ?;

-- name: ListRouterDNSProvidersByTraefikID :many
SELECT
  rdp.*,
  dp.name as provider_name,
  dp.type as provider_type
FROM
  router_dns_provider rdp
  JOIN dns_providers dp ON dp.id = rdp.provider_id
WHERE
  rdp.traefik_id = ?;

-- name: AddRouterDNSProvider :exec
INSERT INTO
  router_dns_provider (traefik_id, provider_id, router_name)
VALUES
  (?, ?, ?);

-- name: DeleteRouterDNSProvider :exec
DELETE FROM router_dns_provider
WHERE
  traefik_id = ?
  AND router_name = ?
  AND provider_id = ?;
