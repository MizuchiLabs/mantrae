-- name: CreateTcpRouterDNSProvider :exec
INSERT
OR IGNORE INTO tcp_router_dns_providers (tcp_router_id, dns_provider_id)
VALUES
  (?, ?);

-- name: GetTcpRouterDNSProvider :one
SELECT
  *
FROM
  tcp_router_dns_providers
WHERE
  tcp_router_id = ?
  AND dns_provider_id = ?;

-- name: ListTcpRouterDNSProviders :many
SELECT
  *
FROM
  tcp_router_dns_providers
WHERE
  tcp_router_id = ?;

-- name: DeleteTcpRouterDNSProvider :exec
DELETE FROM tcp_router_dns_providers
WHERE
  tcp_router_id = ?
  AND dns_provider_id = ?;
