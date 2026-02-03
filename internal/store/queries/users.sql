-- name: CreateUser :one
INSERT INTO
  users (id, username, password, email)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: GetUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = ?;

-- name: GetUserByUsername :one
SELECT
  *
FROM
  users
WHERE
  username = ?;

-- name: GetUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = ?;

-- name: ListUsers :many
SELECT
  *
FROM
  users
ORDER BY
  created_at DESC
LIMIT
  COALESCE(CAST(sqlc.narg ('limit') AS INTEGER), -1)
OFFSET
  COALESCE(CAST(sqlc.narg ('offset') AS INTEGER), 0);

-- name: CountUsers :one
SELECT
  COUNT(*)
FROM
  users;

-- name: UpdateUser :one
UPDATE users
SET
  username = ?,
  email = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
  last_login = CURRENT_TIMESTAMP
WHERE
  id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = ?;

-- name: UpdateUserPassword :exec
UPDATE users
SET
  password = ?
WHERE
  id = ?;
