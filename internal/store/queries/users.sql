-- name: CreateUser :one
INSERT INTO
  users (id, username, password, email, is_admin)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

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
  username
LIMIT
  ?
OFFSET
  ?;

-- name: ListAdminUsers :many
SELECT
  *
FROM
  users
WHERE
  is_admin = TRUE
ORDER BY
  username
LIMIT
  ?
OFFSET
  ?;

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
  is_admin = ?,
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

-- name: UpdateUserResetToken :exec
UPDATE users
SET
  otp = ?,
  otp_expiry = ?
WHERE
  id = ?;

-- name: UpdateUserPassword :exec
UPDATE users
SET
  password = ?,
  otp = '',
  otp_expiry = NULL
WHERE
  id = ?;
