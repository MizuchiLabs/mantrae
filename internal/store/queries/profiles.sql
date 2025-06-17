-- name: CreateProfile :one
INSERT INTO
  profiles (name, description, created_at, updated_at)
VALUES
  (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING *;

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
  name
LIMIT
  ?
OFFSET
  ?;

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
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE
  id = ?;
