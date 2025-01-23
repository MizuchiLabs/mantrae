-- name: UpsertSetting :exec
INSERT INTO
  settings (key, value, description)
VALUES
  (?, ?, ?) ON CONFLICT (key) DO
UPDATE
SET
  value = excluded.value,
  description = excluded.description,
  updated_at = CURRENT_TIMESTAMP;

-- name: GetSetting :one
SELECT
  *
FROM
  settings
WHERE
  key = ?
LIMIT
  1;

-- name: ListSettings :many
SELECT
  *
FROM
  settings;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE
  key = ?;
