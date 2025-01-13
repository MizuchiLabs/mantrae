// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: traefik_config.sql

package db

import (
	"context"
	"time"

	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

const createTraefikConfig = `-- name: CreateTraefikConfig :exec
INSERT INTO
    traefik_config (
        profile_id,
        source,
        entrypoints,
        overview,
        config,
        last_sync
    )
VALUES
    (?, ?, ?, ?, ?, ?)
`

type CreateTraefikConfigParams struct {
	ProfileID   int64                  `json:"profileId"`
	Source      string                 `json:"source"`
	Entrypoints *TraefikEntryPoints    `json:"entrypoints"`
	Overview    *TraefikOverview       `json:"overview"`
	Config      *runtime.Configuration `json:"config"`
	LastSync    *time.Time             `json:"lastSync"`
}

func (q *Queries) CreateTraefikConfig(ctx context.Context, arg CreateTraefikConfigParams) error {
	_, err := q.exec(ctx, q.createTraefikConfigStmt, createTraefikConfig,
		arg.ProfileID,
		arg.Source,
		arg.Entrypoints,
		arg.Overview,
		arg.Config,
		arg.LastSync,
	)
	return err
}

const deleteTraefikConfig = `-- name: DeleteTraefikConfig :exec
DELETE FROM traefik_config
WHERE
    id = ?
`

func (q *Queries) DeleteTraefikConfig(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteTraefikConfigStmt, deleteTraefikConfig, id)
	return err
}

const getTraefikConfig = `-- name: GetTraefikConfig :one
SELECT
    id, profile_id, source, entrypoints, overview, config, last_sync
FROM
    traefik_config
WHERE
    id = ?
`

func (q *Queries) GetTraefikConfig(ctx context.Context, id int64) (TraefikConfig, error) {
	row := q.queryRow(ctx, q.getTraefikConfigStmt, getTraefikConfig, id)
	var i TraefikConfig
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.Source,
		&i.Entrypoints,
		&i.Overview,
		&i.Config,
		&i.LastSync,
	)
	return i, err
}

const getTraefikConfigBySource = `-- name: GetTraefikConfigBySource :one
SELECT
    id, profile_id, source, entrypoints, overview, config, last_sync
FROM
    traefik_config
WHERE
    profile_id = ?
    AND source = ?
`

type GetTraefikConfigBySourceParams struct {
	ProfileID int64  `json:"profileId"`
	Source    string `json:"source"`
}

func (q *Queries) GetTraefikConfigBySource(ctx context.Context, arg GetTraefikConfigBySourceParams) (TraefikConfig, error) {
	row := q.queryRow(ctx, q.getTraefikConfigBySourceStmt, getTraefikConfigBySource, arg.ProfileID, arg.Source)
	var i TraefikConfig
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.Source,
		&i.Entrypoints,
		&i.Overview,
		&i.Config,
		&i.LastSync,
	)
	return i, err
}

const getTraefikConfigLastSync = `-- name: GetTraefikConfigLastSync :one
SELECT
    last_sync
FROM
    traefik_config
WHERE
    id = ?
`

func (q *Queries) GetTraefikConfigLastSync(ctx context.Context, id int64) (*time.Time, error) {
	row := q.queryRow(ctx, q.getTraefikConfigLastSyncStmt, getTraefikConfigLastSync, id)
	var last_sync *time.Time
	err := row.Scan(&last_sync)
	return last_sync, err
}

const updateTraefikConfig = `-- name: UpdateTraefikConfig :exec
UPDATE traefik_config
SET
    source = ?,
    entrypoints = ?,
    overview = ?,
    config = ?,
    last_sync = ?
WHERE
    id = ?
`

type UpdateTraefikConfigParams struct {
	Source      string                 `json:"source"`
	Entrypoints *TraefikEntryPoints    `json:"entrypoints"`
	Overview    *TraefikOverview       `json:"overview"`
	Config      *runtime.Configuration `json:"config"`
	LastSync    *time.Time             `json:"lastSync"`
	ID          int64                  `json:"id"`
}

func (q *Queries) UpdateTraefikConfig(ctx context.Context, arg UpdateTraefikConfigParams) error {
	_, err := q.exec(ctx, q.updateTraefikConfigStmt, updateTraefikConfig,
		arg.Source,
		arg.Entrypoints,
		arg.Overview,
		arg.Config,
		arg.LastSync,
		arg.ID,
	)
	return err
}

const updateTraefikConfigLastSync = `-- name: UpdateTraefikConfigLastSync :exec
UPDATE traefik_config
SET
    last_sync = CURRENT_TIMESTAMP
WHERE
    id = ?
`

func (q *Queries) UpdateTraefikConfigLastSync(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.updateTraefikConfigLastSyncStmt, updateTraefikConfigLastSync, id)
	return err
}
