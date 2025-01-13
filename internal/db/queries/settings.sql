-- name: UpsertSetting :exec
INSERT INTO
    settings (key, value)
VALUES
    (?, ?) ON CONFLICT (key) DO
UPDATE
SET
    value = excluded.value;

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
