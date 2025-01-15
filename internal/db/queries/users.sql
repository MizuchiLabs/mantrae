-- name: CreateUser :exec
INSERT INTO
    users (username, password, email, is_admin)
VALUES
    (?, ?, ?, ?);

-- name: GetUser :one
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

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    username;

-- name: UpdateUser :exec
UPDATE users
SET
    username = ?,
    password = ?,
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
