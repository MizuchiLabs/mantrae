-- name: CreateEntryPoint :one
INSERT INTO
  entry_points (
    profile_id,
    name,
    address,
    is_default,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetEntryPoint :one
SELECT
  *
FROM
  entry_points
WHERE
  id = ?;

-- name: GetDefaultEntryPoint :one
SELECT
  *
FROM
  entry_points
WHERE
  is_default = TRUE
LIMIT
  1;

-- name: UnsetDefaultEntryPoint :exec
UPDATE entry_points
SET
  is_default = FALSE
WHERE
  is_default = TRUE;

-- name: ListEntryPoints :many
SELECT
  *
FROM
  entry_points
WHERE
  profile_id = ?
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: CountEntryPoints :one
SELECT
  COUNT(*)
FROM
  entry_points;

-- name: UpdateEntryPoint :one
UPDATE entry_points
SET
  name = ?,
  address = ?,
  is_default = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteEntryPointByID :exec
DELETE FROM entry_points
WHERE
  id = ?;
