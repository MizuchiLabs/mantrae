// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateProfile(ctx context.Context, arg CreateProfileParams) (Profile, error)
	CreateProvider(ctx context.Context, arg CreateProviderParams) (Provider, error)
	CreateSetting(ctx context.Context, arg CreateSettingParams) (Setting, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAgentByHostname(ctx context.Context, hostname string) error
	DeleteAgentByID(ctx context.Context, id string) error
	DeleteEntryPointByName(ctx context.Context, arg DeleteEntryPointByNameParams) error
	DeleteMiddlewareByID(ctx context.Context, id string) error
	DeleteMiddlewareByName(ctx context.Context, name string) error
	DeleteProfileByID(ctx context.Context, id int64) error
	DeleteProfileByName(ctx context.Context, name string) error
	DeleteProviderByID(ctx context.Context, id int64) error
	DeleteProviderByName(ctx context.Context, name string) error
	DeleteRouterByID(ctx context.Context, id string) error
	DeleteRouterByName(ctx context.Context, name string) error
	DeleteServiceByID(ctx context.Context, id string) error
	DeleteServiceByName(ctx context.Context, name string) error
	DeleteSettingByID(ctx context.Context, id int64) error
	DeleteSettingByKey(ctx context.Context, key string) error
	DeleteUserByID(ctx context.Context, id int64) error
	DeleteUserByUsername(ctx context.Context, username string) error
	GetAgentByHostname(ctx context.Context, arg GetAgentByHostnameParams) (Agent, error)
	GetAgentByID(ctx context.Context, id string) (Agent, error)
	GetDefaultProvider(ctx context.Context) (Provider, error)
	GetMiddlewareByID(ctx context.Context, id string) (Middleware, error)
	GetMiddlewareByName(ctx context.Context, arg GetMiddlewareByNameParams) (Middleware, error)
	GetProfileByID(ctx context.Context, id int64) (Profile, error)
	GetProfileByName(ctx context.Context, name string) (Profile, error)
	GetProviderByID(ctx context.Context, id int64) (Provider, error)
	GetRouterByID(ctx context.Context, id string) (Router, error)
	GetRouterByName(ctx context.Context, arg GetRouterByNameParams) (Router, error)
	GetServiceByID(ctx context.Context, id string) (Service, error)
	GetServiceByName(ctx context.Context, arg GetServiceByNameParams) (Service, error)
	GetSettingByKey(ctx context.Context, key string) (Setting, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAgents(ctx context.Context) ([]Agent, error)
	ListAgentsByProfileID(ctx context.Context, profileID int64) ([]Agent, error)
	ListEntryPoints(ctx context.Context) ([]Entrypoint, error)
	ListEntryPointsByProfileID(ctx context.Context, profileID int64) ([]Entrypoint, error)
	ListMiddlewares(ctx context.Context) ([]Middleware, error)
	ListMiddlewaresByProfileID(ctx context.Context, profileID int64) ([]Middleware, error)
	ListMiddlewaresByProvider(ctx context.Context, provider string) ([]Middleware, error)
	ListProfiles(ctx context.Context) ([]Profile, error)
	ListProviders(ctx context.Context) ([]Provider, error)
	ListRouters(ctx context.Context) ([]Router, error)
	ListRoutersByProfileID(ctx context.Context, profileID int64) ([]Router, error)
	ListRoutersByProvider(ctx context.Context, provider string) ([]Router, error)
	ListServices(ctx context.Context) ([]Service, error)
	ListServicesByProfileID(ctx context.Context, profileID int64) ([]Service, error)
	ListServicesByProvider(ctx context.Context, provider string) ([]Service, error)
	ListSettings(ctx context.Context) ([]Setting, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateProfile(ctx context.Context, arg UpdateProfileParams) (Profile, error)
	UpdateProvider(ctx context.Context, arg UpdateProviderParams) (Provider, error)
	UpdateSetting(ctx context.Context, arg UpdateSettingParams) (Setting, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpsertAgent(ctx context.Context, arg UpsertAgentParams) (Agent, error)
	UpsertEntryPoint(ctx context.Context, arg UpsertEntryPointParams) (Entrypoint, error)
	UpsertMiddleware(ctx context.Context, arg UpsertMiddlewareParams) (Middleware, error)
	UpsertProfile(ctx context.Context, arg UpsertProfileParams) (Profile, error)
	UpsertProvider(ctx context.Context, arg UpsertProviderParams) (Provider, error)
	UpsertRouter(ctx context.Context, arg UpsertRouterParams) (Router, error)
	UpsertService(ctx context.Context, arg UpsertServiceParams) (Service, error)
	UpsertSetting(ctx context.Context, arg UpsertSettingParams) (Setting, error)
	UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
