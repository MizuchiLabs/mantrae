// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateAgent(ctx context.Context, arg CreateAgentParams) (Agent, error)
	CreateConfig(ctx context.Context, arg CreateConfigParams) (Config, error)
	CreateProfile(ctx context.Context, arg CreateProfileParams) (Profile, error)
	CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error)
	CreateSetting(ctx context.Context, arg CreateSettingParams) (Setting, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAgentByHostname(ctx context.Context, hostname string) error
	DeleteAgentByID(ctx context.Context, id string) error
	DeleteConfigByProfileID(ctx context.Context, profileID int64) error
	DeleteConfigByProfileName(ctx context.Context, name string) error
	DeleteProfileByID(ctx context.Context, id int64) error
	DeleteProfileByName(ctx context.Context, name string) error
	DeleteProviderByID(ctx context.Context, id int64) error
	DeleteProviderByName(ctx context.Context, name string) error
	DeleteSettingByID(ctx context.Context, id int64) error
	DeleteSettingByKey(ctx context.Context, key string) error
	DeleteUserByID(ctx context.Context, id int64) error
	DeleteUserByUsername(ctx context.Context, username string) error
	GetAgentByHostname(ctx context.Context, hostname string) (Agent, error)
	GetAgentByID(ctx context.Context, id string) (Agent, error)
	GetConfigByProfileID(ctx context.Context, profileID int64) (Config, error)
	GetConfigByProfileName(ctx context.Context, name string) (Config, error)
	GetProfileByID(ctx context.Context, id int64) (Profile, error)
	GetProfileByName(ctx context.Context, name string) (Profile, error)
	GetProviderByID(ctx context.Context, id int64) (Provider, error)
	GetProviderByName(ctx context.Context, name string) (Provider, error)
	GetSettingByKey(ctx context.Context, key string) (Setting, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAgents(ctx context.Context) ([]Agent, error)
	ListConfigs(ctx context.Context) ([]Config, error)
	ListProfiles(ctx context.Context) ([]Profile, error)
	ListProviders(ctx context.Context) ([]Provider, error)
	ListSettings(ctx context.Context) ([]Setting, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateAgent(ctx context.Context, arg UpdateAgentParams) (Agent, error)
	UpdateConfig(ctx context.Context, arg UpdateConfigParams) (Config, error)
	UpdateProfile(ctx context.Context, arg UpdateProfileParams) (Profile, error)
	UpdateProvider(ctx context.Context, arg UpdateProviderParams) (Provider, error)
	UpdateSetting(ctx context.Context, arg UpdateSettingParams) (Setting, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpsertAgent(ctx context.Context, arg UpsertAgentParams) (Agent, error)
	UpsertProfile(ctx context.Context, arg UpsertProfileParams) (Profile, error)
	UpsertProvider(ctx context.Context, arg UpsertProviderParams) (Provider, error)
	UpsertSetting(ctx context.Context, arg UpsertSettingParams) (Setting, error)
	UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
