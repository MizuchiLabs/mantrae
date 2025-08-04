-- name: CreateDnsProvider :one
INSERT INTO
  dns_providers (
    id,
    name,
    type,
    config,
    is_default,
    created_at,
    updated_at
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ) RETURNING *;

-- name: GetDnsProvider :one
SELECT
  *
FROM
  dns_providers
WHERE
  id = ?;

-- name: GetDnsProviderByName :one
SELECT
  *
FROM
  dns_providers
WHERE
  name = ?;

-- name: GetDefaultDNSProvider :one
SELECT
  *
FROM
  dns_providers
WHERE
  is_default = TRUE
LIMIT
  1;

-- name: UnsetDefaultDNSProvider :exec
UPDATE dns_providers
SET
  is_default = FALSE
WHERE
  is_default = TRUE;

-- name: ListDnsProviders :many
SELECT
  *
FROM
  dns_providers
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountDnsProviders :one
SELECT
  COUNT(*)
FROM
  dns_providers;

-- name: UpdateDnsProvider :one
UPDATE dns_providers
SET
  name = ?,
  type = ?,
  config = ?,
  is_default = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteDnsProvider :exec
DELETE FROM dns_providers
WHERE
  id = ?;
