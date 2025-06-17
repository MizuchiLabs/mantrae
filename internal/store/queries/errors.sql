-- name: LogError :exec
INSERT INTO
  errors (profile_id, category, message, details)
VALUES
  (?, ?, ?, ?) ON CONFLICT (profile_id, category, details) DO
UPDATE
SET
  created_at = CURRENT_TIMESTAMP;

-- name: ListErrors :many
SELECT
  *
FROM
  errors
ORDER BY
  created_at DESC;

-- name: GetErrorsByProfile :many
SELECT
  *
FROM
  errors
WHERE
  profile_id = ?
ORDER BY
  created_at DESC;

-- name: DeleteErrorById :exec
DELETE FROM errors
WHERE
  id = ?;

-- name: DeleteErrorsByProfile :exec
DELETE FROM errors
WHERE
  profile_id = ?;

-- name: DeleteErrorsByProfileCategory :exec
DELETE FROM errors
WHERE
  profile_id = ?
  AND category = ?;
