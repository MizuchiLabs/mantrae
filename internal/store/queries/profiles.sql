-- name: CreateProfile :one
INSERT INTO
  profiles (name, description, token, created_at, updated_at)
VALUES
  (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetProfile :one
SELECT
  *
FROM
  profiles
WHERE
  id = ?;

-- name: GetProfileByName :one
SELECT
  *
FROM
  profiles
WHERE
  name = ?;

-- name: ListProfiles :many
SELECT
  *
FROM
  profiles
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountProfiles :one
SELECT
  COUNT(*)
FROM
  profiles;

-- name: UpdateProfile :one
UPDATE profiles
SET
  name = ?,
  description = ?,
  token = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE
  id = ?;
