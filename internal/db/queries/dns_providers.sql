-- name: CreateDNSProvider :exec
INSERT INTO
    dns_providers (name, type, config, is_active)
VALUES
    (?, ?, ?, ?);

-- name: GetDNSProvider :one
SELECT
    *
FROM
    dns_providers
WHERE
    id = ?;

-- name: GetActiveDNSProvider :one
SELECT
    *
FROM
    dns_providers
WHERE
    is_active = true
LIMIT
    1;

-- name: ListDNSProviders :many
SELECT
    *
FROM
    dns_providers
ORDER BY
    name;

-- name: UpdateDNSProvider :exec
UPDATE dns_providers
SET
    name = ?,
    type = ?,
    config = ?,
    is_active = ?
WHERE
    id = ?;

-- name: DeleteDNSProvider :exec
DELETE FROM dns_providers
WHERE
    id = ?;
