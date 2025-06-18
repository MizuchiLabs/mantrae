-- name: CreateHttpRouterDNSProvider :exec
INSERT
OR IGNORE INTO http_router_dns_providers (http_router_id, dns_provider_id)
VALUES
  (?, ?);

-- name: GetHttpRouterDNSProvider :one
SELECT
  *
FROM
  http_router_dns_providers
WHERE
  http_router_id = ?
  AND dns_provider_id = ?;

-- name: ListHttpRouterDNSProviders :many
SELECT
  *
FROM
  http_router_dns_providers
WHERE
  http_router_id = ?;

-- name: DeleteHttpRouterDNSProvider :exec
DELETE FROM http_router_dns_providers
WHERE
  http_router_id = ?
  AND dns_provider_id = ?;

-- name: GetHttpRouterDomains :many
SELECT
  hr.id AS router_id,
  hr.name AS router_name,
  hr.config AS config_json,
  p.name AS profile_name,
  dp.id AS dns_provider_id,
  dp.name AS dns_provider_name
FROM
  http_routers hr
  JOIN profiles p ON hr.profile_id = p.id
  LEFT JOIN http_router_dns_providers link ON link.http_router_id = hr.id
  LEFT JOIN dns_providers dp ON link.dns_provider_id = dp.id;
