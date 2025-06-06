-- name: CreateProfile :one
INSERT INTO
  profiles (name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?) RETURNING id;

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
  name;

-- name: UpdateProfile :exec
UPDATE profiles
SET
  name = ?,
  url = ?,
  username = ?,
  password = ?,
  tls = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ?;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE
  id = ?;
