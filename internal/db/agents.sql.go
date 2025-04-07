// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: agents.sql

package db

import (
	"context"
)

const createAgent = `-- name: CreateAgent :exec
INSERT INTO
  agents (id, profile_id, token, created_at)
VALUES
  (?, ?, ?, CURRENT_TIMESTAMP)
`

type CreateAgentParams struct {
	ID        string `json:"id"`
	ProfileID int64  `json:"profileId"`
	Token     string `json:"token"`
}

func (q *Queries) CreateAgent(ctx context.Context, arg CreateAgentParams) error {
	_, err := q.exec(ctx, q.createAgentStmt, createAgent, arg.ID, arg.ProfileID, arg.Token)
	return err
}

const deleteAgent = `-- name: DeleteAgent :exec
DELETE FROM agents
WHERE
  id = ?
`

func (q *Queries) DeleteAgent(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteAgentStmt, deleteAgent, id)
	return err
}

const getAgent = `-- name: GetAgent :one
SELECT
  id, profile_id, hostname, public_ip, private_ips, containers, active_ip, token, created_at, updated_at
FROM
  agents
WHERE
  id = ?
`

func (q *Queries) GetAgent(ctx context.Context, id string) (Agent, error) {
	row := q.queryRow(ctx, q.getAgentStmt, getAgent, id)
	var i Agent
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.Hostname,
		&i.PublicIp,
		&i.PrivateIps,
		&i.Containers,
		&i.ActiveIp,
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAgents = `-- name: ListAgents :many
SELECT
  id, profile_id, hostname, public_ip, private_ips, containers, active_ip, token, created_at, updated_at
FROM
  agents
ORDER BY
  hostname
`

func (q *Queries) ListAgents(ctx context.Context) ([]Agent, error) {
	rows, err := q.query(ctx, q.listAgentsStmt, listAgents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Agent
	for rows.Next() {
		var i Agent
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.Hostname,
			&i.PublicIp,
			&i.PrivateIps,
			&i.Containers,
			&i.ActiveIp,
			&i.Token,
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

const listAgentsByProfile = `-- name: ListAgentsByProfile :many
SELECT
  id, profile_id, hostname, public_ip, private_ips, containers, active_ip, token, created_at, updated_at
FROM
  agents
WHERE
  profile_id = ?
`

func (q *Queries) ListAgentsByProfile(ctx context.Context, profileID int64) ([]Agent, error) {
	rows, err := q.query(ctx, q.listAgentsByProfileStmt, listAgentsByProfile, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Agent
	for rows.Next() {
		var i Agent
		if err := rows.Scan(
			&i.ID,
			&i.ProfileID,
			&i.Hostname,
			&i.PublicIp,
			&i.PrivateIps,
			&i.Containers,
			&i.ActiveIp,
			&i.Token,
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

const updateAgent = `-- name: UpdateAgent :one
UPDATE agents
SET
  hostname = COALESCE(?, hostname),
  public_ip = COALESCE(?, public_ip),
  private_ips = COALESCE(?, private_ips),
  containers = COALESCE(?, containers),
  active_ip = COALESCE(?, active_ip),
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING id, profile_id, hostname, public_ip, private_ips, containers, active_ip, token, created_at, updated_at
`

type UpdateAgentParams struct {
	Hostname   *string          `json:"hostname"`
	PublicIp   *string          `json:"publicIp"`
	PrivateIps *AgentPrivateIPs `json:"privateIps"`
	Containers *AgentContainers `json:"containers"`
	ActiveIp   *string          `json:"activeIp"`
	ID         string           `json:"id"`
}

func (q *Queries) UpdateAgent(ctx context.Context, arg UpdateAgentParams) (Agent, error) {
	row := q.queryRow(ctx, q.updateAgentStmt, updateAgent,
		arg.Hostname,
		arg.PublicIp,
		arg.PrivateIps,
		arg.Containers,
		arg.ActiveIp,
		arg.ID,
	)
	var i Agent
	err := row.Scan(
		&i.ID,
		&i.ProfileID,
		&i.Hostname,
		&i.PublicIp,
		&i.PrivateIps,
		&i.Containers,
		&i.ActiveIp,
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAgentIP = `-- name: UpdateAgentIP :exec
UPDATE agents
SET
  active_ip = ?
WHERE
  id = ?
`

type UpdateAgentIPParams struct {
	ActiveIp *string `json:"activeIp"`
	ID       string  `json:"id"`
}

func (q *Queries) UpdateAgentIP(ctx context.Context, arg UpdateAgentIPParams) error {
	_, err := q.exec(ctx, q.updateAgentIPStmt, updateAgentIP, arg.ActiveIp, arg.ID)
	return err
}

const updateAgentToken = `-- name: UpdateAgentToken :exec
UPDATE agents
SET
  token = ?
WHERE
  id = ?
`

type UpdateAgentTokenParams struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

func (q *Queries) UpdateAgentToken(ctx context.Context, arg UpdateAgentTokenParams) error {
	_, err := q.exec(ctx, q.updateAgentTokenStmt, updateAgentToken, arg.Token, arg.ID)
	return err
}
