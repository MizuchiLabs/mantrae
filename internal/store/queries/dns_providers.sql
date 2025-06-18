-- name: CreateDnsProvider :one
INSERT INTO
  dns_providers (
    name,
    type,
    config,
    is_active,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

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

-- name: ListDnsProviders :many
SELECT
  *
FROM
  dns_providers
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

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
  is_active = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteDnsProvider :exec
DELETE FROM dns_providers
WHERE
  id = ?;
