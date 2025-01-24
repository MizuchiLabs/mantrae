-- name: GetRouterDNSProvider :one
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

-- name: UpsertRouterDNSProvider :exec
INSERT INTO
  router_dns_provider (traefik_id, provider_id, router_name)
VALUES
  (?, ?, ?) ON CONFLICT (traefik_id, router_name) DO
UPDATE
SET
  provider_id = EXCLUDED.provider_id;

-- name: DeleteRouterDNSProvider :exec
DELETE FROM router_dns_provider
WHERE
  traefik_id = ?
  AND router_name = ?;
