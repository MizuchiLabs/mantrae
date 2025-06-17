-- name: CreateUdpRouterDNSProvider :exec
INSERT
OR IGNORE INTO udp_router_dns_providers (udp_router_id, dns_provider_id)
VALUES
  (?, ?);

-- name: GetUdpRouterDNSProvider :one
SELECT
  *
FROM
  udp_router_dns_providers
WHERE
  udp_router_id = ?
  AND dns_provider_id = ?;

-- name: ListUdpRouterDNSProvider :many
SELECT
  *
FROM
  udp_router_dns_providers
WHERE
  udp_router_id = ?;

-- name: DeleteUdpRouterDNSProvider :exec
DELETE FROM udp_router_dns_providers
WHERE
  udp_router_id = ?
  AND dns_provider_id = ?;
