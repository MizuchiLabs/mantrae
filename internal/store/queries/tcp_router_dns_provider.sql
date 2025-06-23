-- name: CreateTcpRouterDNSProvider :exec
INSERT
OR IGNORE INTO tcp_router_dns_providers (tcp_router_id, dns_provider_id)
VALUES
  (?, ?);

-- name: GetDnsProvidersByTcpRouter :many
SELECT
  dp.*
FROM
  dns_providers dp
  JOIN tcp_router_dns_providers trdp ON dp.id = trdp.dns_provider_id
WHERE
  trdp.tcp_router_id = ?;

-- name: DeleteTcpRouterDNSProvider :exec
DELETE FROM tcp_router_dns_providers
WHERE
  tcp_router_id = ?
  AND dns_provider_id = ?;

-- name: GetTcpRouterDomains :many
SELECT
  tr.id AS router_id,
  tr.name AS router_name,
  tr.config AS config_json,
  p.name AS profile_name,
  dp.id AS dns_provider_id,
  dp.name AS dns_provider_name
FROM
  tcp_routers tr
  JOIN profiles p ON tr.profile_id = p.id
  LEFT JOIN tcp_router_dns_providers link ON link.tcp_router_id = tr.id
  LEFT JOIN dns_providers dp ON link.dns_provider_id = dp.id;
