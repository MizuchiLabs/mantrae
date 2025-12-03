-- name: CreateEntryPoint :one
INSERT INTO
  entry_points (id, profile_id, name, address, is_default)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

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
  is_default = TRUE
  AND profile_id = ?;

-- name: ListEntryPoints :many
SELECT
  *
FROM
  entry_points
WHERE
  profile_id = sqlc.arg ('profile_id')
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountEntryPoints :one
SELECT
  COUNT(*)
FROM
  entry_points
WHERE
  profile_id = ?;

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
