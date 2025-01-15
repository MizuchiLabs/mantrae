-- name: UpsertSetting :exec
INSERT INTO
    settings (key, value)
VALUES
    (?, ?) ON CONFLICT (key) DO
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
