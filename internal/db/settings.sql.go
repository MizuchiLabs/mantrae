// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: settings.sql

package db

import (
	"context"
)

const deleteSetting = `-- name: DeleteSetting :exec
DELETE FROM settings
WHERE
    key = ?
`

func (q *Queries) DeleteSetting(ctx context.Context, key string) error {
	_, err := q.exec(ctx, q.deleteSettingStmt, deleteSetting, key)
	return err
}

const getSetting = `-- name: GetSetting :one
SELECT
    "key", value, updated_at
FROM
    settings
WHERE
    key = ?
`

func (q *Queries) GetSetting(ctx context.Context, key string) (Setting, error) {
	row := q.queryRow(ctx, q.getSettingStmt, getSetting, key)
	var i Setting
	err := row.Scan(&i.Key, &i.Value, &i.UpdatedAt)
	return i, err
}

const listSettings = `-- name: ListSettings :many
SELECT
    "key", value, updated_at
FROM
    settings
`

func (q *Queries) ListSettings(ctx context.Context) ([]Setting, error) {
	rows, err := q.query(ctx, q.listSettingsStmt, listSettings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Setting
	for rows.Next() {
		var i Setting
		if err := rows.Scan(&i.Key, &i.Value, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertSetting = `-- name: UpsertSetting :exec
INSERT INTO
    settings (key, value)
VALUES
    (?, ?) ON CONFLICT (key) DO
UPDATE
SET
    value = excluded.value
`

type UpsertSettingParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (q *Queries) UpsertSetting(ctx context.Context, arg UpsertSettingParams) error {
	_, err := q.exec(ctx, q.upsertSettingStmt, upsertSetting, arg.Key, arg.Value)
	return err
}
