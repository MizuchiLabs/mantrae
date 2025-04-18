// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: profiles.sql

package db

import (
	"context"
)

const createProfile = `-- name: CreateProfile :one
INSERT INTO
  profiles (name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?) RETURNING id
`

type CreateProfileParams struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Tls      bool    `json:"tls"`
}

func (q *Queries) CreateProfile(ctx context.Context, arg CreateProfileParams) (int64, error) {
	row := q.queryRow(ctx, q.createProfileStmt, createProfile,
		arg.Name,
		arg.Url,
		arg.Username,
		arg.Password,
		arg.Tls,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteProfile = `-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE
  id = ?
`

func (q *Queries) DeleteProfile(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProfileStmt, deleteProfile, id)
	return err
}

const getProfile = `-- name: GetProfile :one
SELECT
  id, name, url, username, password, tls, created_at, updated_at
FROM
  profiles
WHERE
  id = ?
`

func (q *Queries) GetProfile(ctx context.Context, id int64) (Profile, error) {
	row := q.queryRow(ctx, q.getProfileStmt, getProfile, id)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Username,
		&i.Password,
		&i.Tls,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProfileByName = `-- name: GetProfileByName :one
SELECT
  id, name, url, username, password, tls, created_at, updated_at
FROM
  profiles
WHERE
  name = ?
`

func (q *Queries) GetProfileByName(ctx context.Context, name string) (Profile, error) {
	row := q.queryRow(ctx, q.getProfileByNameStmt, getProfileByName, name)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Username,
		&i.Password,
		&i.Tls,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProfiles = `-- name: ListProfiles :many
SELECT
  id, name, url, username, password, tls, created_at, updated_at
FROM
  profiles
ORDER BY
  name
`

func (q *Queries) ListProfiles(ctx context.Context) ([]Profile, error) {
	rows, err := q.query(ctx, q.listProfilesStmt, listProfiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Profile
	for rows.Next() {
		var i Profile
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Username,
			&i.Password,
			&i.Tls,
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

const updateProfile = `-- name: UpdateProfile :exec
UPDATE profiles
SET
  name = ?,
  url = ?,
  username = ?,
  password = ?,
  tls = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ?
`

type UpdateProfileParams struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Tls      bool    `json:"tls"`
	ID       int64   `json:"id"`
}

func (q *Queries) UpdateProfile(ctx context.Context, arg UpdateProfileParams) error {
	_, err := q.exec(ctx, q.updateProfileStmt, updateProfile,
		arg.Name,
		arg.Url,
		arg.Username,
		arg.Password,
		arg.Tls,
		arg.ID,
	)
	return err
}
