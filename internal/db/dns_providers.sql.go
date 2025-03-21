// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: dns_providers.sql

package db

import (
	"context"
)

const createDNSProvider = `-- name: CreateDNSProvider :exec
INSERT INTO
    dns_providers (name, type, config, is_active)
VALUES
    (?, ?, ?, ?)
`

type CreateDNSProviderParams struct {
	Name     string             `json:"name"`
	Type     string             `json:"type"`
	Config   *DNSProviderConfig `json:"config"`
	IsActive bool               `json:"isActive"`
}

func (q *Queries) CreateDNSProvider(ctx context.Context, arg CreateDNSProviderParams) error {
	_, err := q.exec(ctx, q.createDNSProviderStmt, createDNSProvider,
		arg.Name,
		arg.Type,
		arg.Config,
		arg.IsActive,
	)
	return err
}

const deleteDNSProvider = `-- name: DeleteDNSProvider :exec
DELETE FROM dns_providers
WHERE
    id = ?
`

func (q *Queries) DeleteDNSProvider(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteDNSProviderStmt, deleteDNSProvider, id)
	return err
}

const getActiveDNSProvider = `-- name: GetActiveDNSProvider :one
SELECT
    id, name, type, config, is_active, created_at, updated_at
FROM
    dns_providers
WHERE
    is_active = true
LIMIT
    1
`

func (q *Queries) GetActiveDNSProvider(ctx context.Context) (DnsProvider, error) {
	row := q.queryRow(ctx, q.getActiveDNSProviderStmt, getActiveDNSProvider)
	var i DnsProvider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDNSProvider = `-- name: GetDNSProvider :one
SELECT
    id, name, type, config, is_active, created_at, updated_at
FROM
    dns_providers
WHERE
    id = ?
`

func (q *Queries) GetDNSProvider(ctx context.Context, id int64) (DnsProvider, error) {
	row := q.queryRow(ctx, q.getDNSProviderStmt, getDNSProvider, id)
	var i DnsProvider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Config,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listDNSProviders = `-- name: ListDNSProviders :many
SELECT
    id, name, type, config, is_active, created_at, updated_at
FROM
    dns_providers
ORDER BY
    name
`

func (q *Queries) ListDNSProviders(ctx context.Context) ([]DnsProvider, error) {
	rows, err := q.query(ctx, q.listDNSProvidersStmt, listDNSProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DnsProvider
	for rows.Next() {
		var i DnsProvider
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Config,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const updateDNSProvider = `-- name: UpdateDNSProvider :exec
UPDATE dns_providers
SET
    name = ?,
    type = ?,
    config = ?,
    is_active = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?
`

type UpdateDNSProviderParams struct {
	Name     string             `json:"name"`
	Type     string             `json:"type"`
	Config   *DNSProviderConfig `json:"config"`
	IsActive bool               `json:"isActive"`
	ID       int64              `json:"id"`
}

func (q *Queries) UpdateDNSProvider(ctx context.Context, arg UpdateDNSProviderParams) error {
	_, err := q.exec(ctx, q.updateDNSProviderStmt, updateDNSProvider,
		arg.Name,
		arg.Type,
		arg.Config,
		arg.IsActive,
		arg.ID,
	)
	return err
}
