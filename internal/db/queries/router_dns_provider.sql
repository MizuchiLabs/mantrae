-- name: CreateRouterDNSProvider :exec
INSERT INTO
    router_dns_provider (traefik_id, provider_id, router_name)
VALUES
    (?, ?, ?);

-- name: GetRouterDNSProviders :many
SELECT
    *
FROM
    router_dns_provider
WHERE
    traefik_id = ?;

-- name: DeleteRouterDNSProvider :exec
DELETE FROM router_dns_provider
WHERE
    traefik_id = ?
    AND provider_id = ?
    AND router_name = ?;

-- name: DeleteRouterDNSProvidersByTraefik :exec
DELETE FROM router_dns_provider
WHERE
    traefik_id = ?;
