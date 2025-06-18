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

-- name: ListEntryPoints :many
SELECT
  *
FROM
  entry_points
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

-- name: DeleteEntryPoint :exec
DELETE FROM entry_points
WHERE
  id = ?;
