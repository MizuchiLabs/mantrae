-- name: UpsertSetting :exec
INSERT INTO
  settings (key, value, updated_at)
VALUES
  (?, ?, CURRENT_TIMESTAMP) ON CONFLICT (key) DO
UPDATE
SET
  value = excluded.value,
  updated_at = CURRENT_TIMESTAMP;

-- name: GetSetting :one
SELECT
  *
FROM
  settings
WHERE
  key = ?;

-- name: ListSettings :many
SELECT
  *
FROM
  settings;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE
  key = ?;
