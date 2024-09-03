// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const createConfig = `-- name: CreateConfig :one
INSERT INTO
  config (
    profile_id,
    entrypoints,
    routers,
    services,
    middlewares,
    version
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING profile_id, entrypoints, routers, services, middlewares, version
`

type CreateConfigParams struct {
	ProfileID   int64       `json:"profile_id"`
	Entrypoints interface{} `json:"entrypoints"`
	Routers     interface{} `json:"routers"`
	Services    interface{} `json:"services"`
	Middlewares interface{} `json:"middlewares"`
	Version     *string     `json:"version"`
}

func (q *Queries) CreateConfig(ctx context.Context, arg CreateConfigParams) (Config, error) {
	row := q.queryRow(ctx, q.createConfigStmt, createConfig,
		arg.ProfileID,
		arg.Entrypoints,
		arg.Routers,
		arg.Services,
		arg.Middlewares,
		arg.Version,
	)
	var i Config
	err := row.Scan(
		&i.ProfileID,
		&i.Entrypoints,
		&i.Routers,
		&i.Services,
		&i.Middlewares,
		&i.Version,
	)
	return i, err
}

const createCredential = `-- name: CreateCredential :exec
INSERT INTO
  credentials (username, password)
VALUES
  (?, ?)
`

type CreateCredentialParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateCredential(ctx context.Context, arg CreateCredentialParams) error {
	_, err := q.exec(ctx, q.createCredentialStmt, createCredential, arg.Username, arg.Password)
	return err
}

const createProfile = `-- name: CreateProfile :one
INSERT INTO
  profiles (name, url, username, password, tls)
VALUES
  (?, ?, ?, ?, ?) RETURNING id, name, url, username, password, tls
`

type CreateProfileParams struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Tls      bool    `json:"tls"`
}

func (q *Queries) CreateProfile(ctx context.Context, arg CreateProfileParams) (Profile, error) {
	row := q.queryRow(ctx, q.createProfileStmt, createProfile,
		arg.Name,
		arg.Url,
		arg.Username,
		arg.Password,
		arg.Tls,
	)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Username,
		&i.Password,
		&i.Tls,
	)
	return i, err
}

const createProvider = `-- name: CreateProvider :one
INSERT INTO
  providers (
    name,
    type,
    external_ip,
    api_key,
    api_url,
    is_active
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING id, name, type, external_ip, api_key, api_url, is_active
`

type CreateProviderParams struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	ExternalIp string  `json:"external_ip"`
	ApiKey     string  `json:"api_key"`
	ApiUrl     *string `json:"api_url"`
	IsActive   bool    `json:"is_active"`
}

func (q *Queries) CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error) {
	row := q.queryRow(ctx, q.createProviderStmt, createProvider,
		arg.Name,
		arg.Type,
		arg.ExternalIp,
		arg.ApiKey,
		arg.ApiUrl,
		arg.IsActive,
	)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.ExternalIp,
		&i.ApiKey,
		&i.ApiUrl,
		&i.IsActive,
	)
	return i, err
}

const deleteConfigByProfileID = `-- name: DeleteConfigByProfileID :exec
DELETE FROM config
WHERE
  profile_id = ?
`

func (q *Queries) DeleteConfigByProfileID(ctx context.Context, profileID int64) error {
	_, err := q.exec(ctx, q.deleteConfigByProfileIDStmt, deleteConfigByProfileID, profileID)
	return err
}

const deleteConfigByProfileName = `-- name: DeleteConfigByProfileName :exec
DELETE FROM config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  )
`

func (q *Queries) DeleteConfigByProfileName(ctx context.Context, name string) error {
	_, err := q.exec(ctx, q.deleteConfigByProfileNameStmt, deleteConfigByProfileName, name)
	return err
}

const deleteCredentialByID = `-- name: DeleteCredentialByID :exec
DELETE FROM credentials
WHERE
  id = ?
`

func (q *Queries) DeleteCredentialByID(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteCredentialByIDStmt, deleteCredentialByID, id)
	return err
}

const deleteCredentialByUsername = `-- name: DeleteCredentialByUsername :exec
DELETE FROM credentials
WHERE
  username = ?
`

func (q *Queries) DeleteCredentialByUsername(ctx context.Context, username string) error {
	_, err := q.exec(ctx, q.deleteCredentialByUsernameStmt, deleteCredentialByUsername, username)
	return err
}

const deleteProfileByID = `-- name: DeleteProfileByID :exec
DELETE FROM profiles
WHERE
  id = ?
`

func (q *Queries) DeleteProfileByID(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProfileByIDStmt, deleteProfileByID, id)
	return err
}

const deleteProfileByName = `-- name: DeleteProfileByName :exec
DELETE FROM profiles
WHERE
  name = ?
`

func (q *Queries) DeleteProfileByName(ctx context.Context, name string) error {
	_, err := q.exec(ctx, q.deleteProfileByNameStmt, deleteProfileByName, name)
	return err
}

const deleteProviderByID = `-- name: DeleteProviderByID :exec
DELETE FROM providers
WHERE
  id = ?
`

func (q *Queries) DeleteProviderByID(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteProviderByIDStmt, deleteProviderByID, id)
	return err
}

const deleteProviderByName = `-- name: DeleteProviderByName :exec
DELETE FROM providers
WHERE
  name = ?
`

func (q *Queries) DeleteProviderByName(ctx context.Context, name string) error {
	_, err := q.exec(ctx, q.deleteProviderByNameStmt, deleteProviderByName, name)
	return err
}

const getConfigByProfileID = `-- name: GetConfigByProfileID :one
SELECT
  profile_id, entrypoints, routers, services, middlewares, version
FROM
  config
WHERE
  profile_id = ?
LIMIT
  1
`

func (q *Queries) GetConfigByProfileID(ctx context.Context, profileID int64) (Config, error) {
	row := q.queryRow(ctx, q.getConfigByProfileIDStmt, getConfigByProfileID, profileID)
	var i Config
	err := row.Scan(
		&i.ProfileID,
		&i.Entrypoints,
		&i.Routers,
		&i.Services,
		&i.Middlewares,
		&i.Version,
	)
	return i, err
}

const getConfigByProfileName = `-- name: GetConfigByProfileName :one
SELECT
  profile_id, entrypoints, routers, services, middlewares, version
FROM
  config
WHERE
  profile_id = (
    SELECT
      id
    FROM
      profiles
    WHERE
      name = ?
  )
LIMIT
  1
`

func (q *Queries) GetConfigByProfileName(ctx context.Context, name string) (Config, error) {
	row := q.queryRow(ctx, q.getConfigByProfileNameStmt, getConfigByProfileName, name)
	var i Config
	err := row.Scan(
		&i.ProfileID,
		&i.Entrypoints,
		&i.Routers,
		&i.Services,
		&i.Middlewares,
		&i.Version,
	)
	return i, err
}

const getCredentialByID = `-- name: GetCredentialByID :one
SELECT
  id, username, password
FROM
  credentials
WHERE
  id = ?
LIMIT
  1
`

func (q *Queries) GetCredentialByID(ctx context.Context, id int64) (Credential, error) {
	row := q.queryRow(ctx, q.getCredentialByIDStmt, getCredentialByID, id)
	var i Credential
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}

const getCredentialByUsername = `-- name: GetCredentialByUsername :one
SELECT
  id, username, password
FROM
  credentials
WHERE
  username = ?
LIMIT
  1
`

func (q *Queries) GetCredentialByUsername(ctx context.Context, username string) (Credential, error) {
	row := q.queryRow(ctx, q.getCredentialByUsernameStmt, getCredentialByUsername, username)
	var i Credential
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}

const getProfileByID = `-- name: GetProfileByID :one
SELECT
  id, name, url, username, password, tls
FROM
  profiles
WHERE
  id = ?
LIMIT
  1
`

func (q *Queries) GetProfileByID(ctx context.Context, id int64) (Profile, error) {
	row := q.queryRow(ctx, q.getProfileByIDStmt, getProfileByID, id)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Username,
		&i.Password,
		&i.Tls,
	)
	return i, err
}

const getProfileByName = `-- name: GetProfileByName :one
SELECT
  id, name, url, username, password, tls
FROM
  profiles
WHERE
  name = ?
LIMIT
  1
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
	)
	return i, err
}

const getProviderByID = `-- name: GetProviderByID :one
SELECT
  id, name, type, external_ip, api_key, api_url, is_active
FROM
  providers
WHERE
  id = ?
LIMIT
  1
`

func (q *Queries) GetProviderByID(ctx context.Context, id int64) (Provider, error) {
	row := q.queryRow(ctx, q.getProviderByIDStmt, getProviderByID, id)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.ExternalIp,
		&i.ApiKey,
		&i.ApiUrl,
		&i.IsActive,
	)
	return i, err
}

const getProviderByName = `-- name: GetProviderByName :one
SELECT
  id, name, type, external_ip, api_key, api_url, is_active
FROM
  providers
WHERE
  name = ?
LIMIT
  1
`

func (q *Queries) GetProviderByName(ctx context.Context, name string) (Provider, error) {
	row := q.queryRow(ctx, q.getProviderByNameStmt, getProviderByName, name)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.ExternalIp,
		&i.ApiKey,
		&i.ApiUrl,
		&i.IsActive,
	)
	return i, err
}

const listConfigs = `-- name: ListConfigs :many
SELECT
  profile_id, entrypoints, routers, services, middlewares, version
FROM
  config
`

func (q *Queries) ListConfigs(ctx context.Context) ([]Config, error) {
	rows, err := q.query(ctx, q.listConfigsStmt, listConfigs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Config
	for rows.Next() {
		var i Config
		if err := rows.Scan(
			&i.ProfileID,
			&i.Entrypoints,
			&i.Routers,
			&i.Services,
			&i.Middlewares,
			&i.Version,
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

const listCredentials = `-- name: ListCredentials :many
SELECT
  id, username, password
FROM
  credentials
`

func (q *Queries) ListCredentials(ctx context.Context) ([]Credential, error) {
	rows, err := q.query(ctx, q.listCredentialsStmt, listCredentials)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Credential
	for rows.Next() {
		var i Credential
		if err := rows.Scan(&i.ID, &i.Username, &i.Password); err != nil {
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

const listProfiles = `-- name: ListProfiles :many
SELECT
  id, name, url, username, password, tls
FROM
  profiles
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

const listProviders = `-- name: ListProviders :many
SELECT
  id, name, type, external_ip, api_key, api_url, is_active
FROM
  providers
`

func (q *Queries) ListProviders(ctx context.Context) ([]Provider, error) {
	rows, err := q.query(ctx, q.listProvidersStmt, listProviders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Provider
	for rows.Next() {
		var i Provider
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.ExternalIp,
			&i.ApiKey,
			&i.ApiUrl,
			&i.IsActive,
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

const updateConfig = `-- name: UpdateConfig :one
UPDATE config
SET
  entrypoints = ?,
  routers = ?,
  services = ?,
  middlewares = ?,
  version = ?
WHERE
  profile_id = ? RETURNING profile_id, entrypoints, routers, services, middlewares, version
`

type UpdateConfigParams struct {
	Entrypoints interface{} `json:"entrypoints"`
	Routers     interface{} `json:"routers"`
	Services    interface{} `json:"services"`
	Middlewares interface{} `json:"middlewares"`
	Version     *string     `json:"version"`
	ProfileID   int64       `json:"profile_id"`
}

func (q *Queries) UpdateConfig(ctx context.Context, arg UpdateConfigParams) (Config, error) {
	row := q.queryRow(ctx, q.updateConfigStmt, updateConfig,
		arg.Entrypoints,
		arg.Routers,
		arg.Services,
		arg.Middlewares,
		arg.Version,
		arg.ProfileID,
	)
	var i Config
	err := row.Scan(
		&i.ProfileID,
		&i.Entrypoints,
		&i.Routers,
		&i.Services,
		&i.Middlewares,
		&i.Version,
	)
	return i, err
}

const updateCredential = `-- name: UpdateCredential :exec
UPDATE credentials
SET
  username = ?,
  password = ?
WHERE
  id = ?
`

type UpdateCredentialParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateCredential(ctx context.Context, arg UpdateCredentialParams) error {
	_, err := q.exec(ctx, q.updateCredentialStmt, updateCredential, arg.Username, arg.Password, arg.ID)
	return err
}

const updateProfile = `-- name: UpdateProfile :one
UPDATE profiles
SET
  name = ?,
  url = ?,
  username = ?,
  password = ?,
  tls = ?
WHERE
  id = ? RETURNING id, name, url, username, password, tls
`

type UpdateProfileParams struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Tls      bool    `json:"tls"`
	ID       int64   `json:"id"`
}

func (q *Queries) UpdateProfile(ctx context.Context, arg UpdateProfileParams) (Profile, error) {
	row := q.queryRow(ctx, q.updateProfileStmt, updateProfile,
		arg.Name,
		arg.Url,
		arg.Username,
		arg.Password,
		arg.Tls,
		arg.ID,
	)
	var i Profile
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Username,
		&i.Password,
		&i.Tls,
	)
	return i, err
}

const updateProvider = `-- name: UpdateProvider :one
UPDATE providers
SET
  name = ?,
  type = ?,
  external_ip = ?,
  api_key = ?,
  api_url = ?,
  is_active = ?
WHERE
  id = ? RETURNING id, name, type, external_ip, api_key, api_url, is_active
`

type UpdateProviderParams struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	ExternalIp string  `json:"external_ip"`
	ApiKey     string  `json:"api_key"`
	ApiUrl     *string `json:"api_url"`
	IsActive   bool    `json:"is_active"`
	ID         int64   `json:"id"`
}

func (q *Queries) UpdateProvider(ctx context.Context, arg UpdateProviderParams) (Provider, error) {
	row := q.queryRow(ctx, q.updateProviderStmt, updateProvider,
		arg.Name,
		arg.Type,
		arg.ExternalIp,
		arg.ApiKey,
		arg.ApiUrl,
		arg.IsActive,
		arg.ID,
	)
	var i Provider
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.ExternalIp,
		&i.ApiKey,
		&i.ApiUrl,
		&i.IsActive,
	)
	return i, err
}

const validateAuth = `-- name: ValidateAuth :one
SELECT
  id,
  username
FROM
  credentials
WHERE
  username = ?
  AND password = ?
`

type ValidateAuthParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ValidateAuthRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) ValidateAuth(ctx context.Context, arg ValidateAuthParams) (ValidateAuthRow, error) {
	row := q.queryRow(ctx, q.validateAuthStmt, validateAuth, arg.Username, arg.Password)
	var i ValidateAuthRow
	err := row.Scan(&i.ID, &i.Username)
	return i, err
}
