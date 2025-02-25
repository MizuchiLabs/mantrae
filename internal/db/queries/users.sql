-- name: CreateUser :exec
INSERT INTO
  users (username, password, email, is_admin)
VALUES
  (?, ?, ?, ?);

-- name: GetUser :one
SELECT
  id,
  username,
  email,
  is_admin,
  otp,
  otp_expiry,
  last_login,
  created_at,
  updated_at
FROM
  users
WHERE
  id = ?;

-- name: GetUserByUsername :one
SELECT
  id,
  username,
  email,
  is_admin,
  otp,
  otp_expiry,
  last_login,
  created_at,
  updated_at
FROM
  users
WHERE
  username = ?;

-- name: GetUserPassword :one
SELECT
  password
FROM
  users
WHERE
  id = ?;

-- name: ListUsers :many
SELECT
  id,
  username,
  email,
  is_admin,
  otp,
  otp_expiry,
  last_login,
  created_at,
  updated_at
FROM
  users
ORDER BY
  username;

-- name: UpdateUser :exec
UPDATE users
SET
  username = ?,
  email = ?,
  is_admin = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ?;

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
