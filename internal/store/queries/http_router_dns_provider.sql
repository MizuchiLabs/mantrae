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
