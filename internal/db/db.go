// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createAgentStmt, err = db.PrepareContext(ctx, createAgent); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAgent: %w", err)
	}
	if q.createConfigStmt, err = db.PrepareContext(ctx, createConfig); err != nil {
		return nil, fmt.Errorf("error preparing query CreateConfig: %w", err)
	}
	if q.createProfileStmt, err = db.PrepareContext(ctx, createProfile); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProfile: %w", err)
	}
	if q.createProviderStmt, err = db.PrepareContext(ctx, createProvider); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProvider: %w", err)
	}
	if q.createRouterStmt, err = db.PrepareContext(ctx, createRouter); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRouter: %w", err)
	}
	if q.createServiceStmt, err = db.PrepareContext(ctx, createService); err != nil {
		return nil, fmt.Errorf("error preparing query CreateService: %w", err)
	}
	if q.createSettingStmt, err = db.PrepareContext(ctx, createSetting); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSetting: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteAgentByHostnameStmt, err = db.PrepareContext(ctx, deleteAgentByHostname); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAgentByHostname: %w", err)
	}
	if q.deleteAgentByIDStmt, err = db.PrepareContext(ctx, deleteAgentByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAgentByID: %w", err)
	}
	if q.deleteConfigByProfileIDStmt, err = db.PrepareContext(ctx, deleteConfigByProfileID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteConfigByProfileID: %w", err)
	}
	if q.deleteConfigByProfileNameStmt, err = db.PrepareContext(ctx, deleteConfigByProfileName); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteConfigByProfileName: %w", err)
	}
	if q.deleteProfileByIDStmt, err = db.PrepareContext(ctx, deleteProfileByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfileByID: %w", err)
	}
	if q.deleteProfileByNameStmt, err = db.PrepareContext(ctx, deleteProfileByName); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfileByName: %w", err)
	}
	if q.deleteProviderByIDStmt, err = db.PrepareContext(ctx, deleteProviderByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProviderByID: %w", err)
	}
	if q.deleteProviderByNameStmt, err = db.PrepareContext(ctx, deleteProviderByName); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProviderByName: %w", err)
	}
	if q.deleteRouterByIDStmt, err = db.PrepareContext(ctx, deleteRouterByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRouterByID: %w", err)
	}
	if q.deleteRouterByNameStmt, err = db.PrepareContext(ctx, deleteRouterByName); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRouterByName: %w", err)
	}
	if q.deleteServiceByIDStmt, err = db.PrepareContext(ctx, deleteServiceByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteServiceByID: %w", err)
	}
	if q.deleteServiceByNameStmt, err = db.PrepareContext(ctx, deleteServiceByName); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteServiceByName: %w", err)
	}
	if q.deleteSettingByIDStmt, err = db.PrepareContext(ctx, deleteSettingByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSettingByID: %w", err)
	}
	if q.deleteSettingByKeyStmt, err = db.PrepareContext(ctx, deleteSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSettingByKey: %w", err)
	}
	if q.deleteUserByIDStmt, err = db.PrepareContext(ctx, deleteUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserByID: %w", err)
	}
	if q.deleteUserByUsernameStmt, err = db.PrepareContext(ctx, deleteUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserByUsername: %w", err)
	}
	if q.getAgentByHostnameStmt, err = db.PrepareContext(ctx, getAgentByHostname); err != nil {
		return nil, fmt.Errorf("error preparing query GetAgentByHostname: %w", err)
	}
	if q.getAgentByIDStmt, err = db.PrepareContext(ctx, getAgentByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetAgentByID: %w", err)
	}
	if q.getConfigByProfileIDStmt, err = db.PrepareContext(ctx, getConfigByProfileID); err != nil {
		return nil, fmt.Errorf("error preparing query GetConfigByProfileID: %w", err)
	}
	if q.getConfigByProfileNameStmt, err = db.PrepareContext(ctx, getConfigByProfileName); err != nil {
		return nil, fmt.Errorf("error preparing query GetConfigByProfileName: %w", err)
	}
	if q.getProfileByIDStmt, err = db.PrepareContext(ctx, getProfileByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileByID: %w", err)
	}
	if q.getProfileByNameStmt, err = db.PrepareContext(ctx, getProfileByName); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileByName: %w", err)
	}
	if q.getProviderByIDStmt, err = db.PrepareContext(ctx, getProviderByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetProviderByID: %w", err)
	}
	if q.getProviderByNameStmt, err = db.PrepareContext(ctx, getProviderByName); err != nil {
		return nil, fmt.Errorf("error preparing query GetProviderByName: %w", err)
	}
	if q.getRouterByIDStmt, err = db.PrepareContext(ctx, getRouterByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetRouterByID: %w", err)
	}
	if q.getServiceByIDStmt, err = db.PrepareContext(ctx, getServiceByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetServiceByID: %w", err)
	}
	if q.getSettingByKeyStmt, err = db.PrepareContext(ctx, getSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettingByKey: %w", err)
	}
	if q.getUserByIDStmt, err = db.PrepareContext(ctx, getUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByID: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.listAgentsStmt, err = db.PrepareContext(ctx, listAgents); err != nil {
		return nil, fmt.Errorf("error preparing query ListAgents: %w", err)
	}
	if q.listConfigsStmt, err = db.PrepareContext(ctx, listConfigs); err != nil {
		return nil, fmt.Errorf("error preparing query ListConfigs: %w", err)
	}
	if q.listProfilesStmt, err = db.PrepareContext(ctx, listProfiles); err != nil {
		return nil, fmt.Errorf("error preparing query ListProfiles: %w", err)
	}
	if q.listProvidersStmt, err = db.PrepareContext(ctx, listProviders); err != nil {
		return nil, fmt.Errorf("error preparing query ListProviders: %w", err)
	}
	if q.listRoutersStmt, err = db.PrepareContext(ctx, listRouters); err != nil {
		return nil, fmt.Errorf("error preparing query ListRouters: %w", err)
	}
	if q.listRoutersByProfileIDStmt, err = db.PrepareContext(ctx, listRoutersByProfileID); err != nil {
		return nil, fmt.Errorf("error preparing query ListRoutersByProfileID: %w", err)
	}
	if q.listServicesStmt, err = db.PrepareContext(ctx, listServices); err != nil {
		return nil, fmt.Errorf("error preparing query ListServices: %w", err)
	}
	if q.listServicesByProfileIDStmt, err = db.PrepareContext(ctx, listServicesByProfileID); err != nil {
		return nil, fmt.Errorf("error preparing query ListServicesByProfileID: %w", err)
	}
	if q.listSettingsStmt, err = db.PrepareContext(ctx, listSettings); err != nil {
		return nil, fmt.Errorf("error preparing query ListSettings: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
	}
	if q.updateAgentStmt, err = db.PrepareContext(ctx, updateAgent); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAgent: %w", err)
	}
	if q.updateConfigStmt, err = db.PrepareContext(ctx, updateConfig); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateConfig: %w", err)
	}
	if q.updateProfileStmt, err = db.PrepareContext(ctx, updateProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfile: %w", err)
	}
	if q.updateProviderStmt, err = db.PrepareContext(ctx, updateProvider); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProvider: %w", err)
	}
	if q.updateRouterStmt, err = db.PrepareContext(ctx, updateRouter); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateRouter: %w", err)
	}
	if q.updateServiceStmt, err = db.PrepareContext(ctx, updateService); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateService: %w", err)
	}
	if q.updateSettingStmt, err = db.PrepareContext(ctx, updateSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSetting: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	if q.upsertAgentStmt, err = db.PrepareContext(ctx, upsertAgent); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertAgent: %w", err)
	}
	if q.upsertProfileStmt, err = db.PrepareContext(ctx, upsertProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertProfile: %w", err)
	}
	if q.upsertProviderStmt, err = db.PrepareContext(ctx, upsertProvider); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertProvider: %w", err)
	}
	if q.upsertRouterStmt, err = db.PrepareContext(ctx, upsertRouter); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertRouter: %w", err)
	}
	if q.upsertServiceStmt, err = db.PrepareContext(ctx, upsertService); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertService: %w", err)
	}
	if q.upsertSettingStmt, err = db.PrepareContext(ctx, upsertSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertSetting: %w", err)
	}
	if q.upsertUserStmt, err = db.PrepareContext(ctx, upsertUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertUser: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createAgentStmt != nil {
		if cerr := q.createAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAgentStmt: %w", cerr)
		}
	}
	if q.createConfigStmt != nil {
		if cerr := q.createConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createConfigStmt: %w", cerr)
		}
	}
	if q.createProfileStmt != nil {
		if cerr := q.createProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProfileStmt: %w", cerr)
		}
	}
	if q.createProviderStmt != nil {
		if cerr := q.createProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProviderStmt: %w", cerr)
		}
	}
	if q.createRouterStmt != nil {
		if cerr := q.createRouterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRouterStmt: %w", cerr)
		}
	}
	if q.createServiceStmt != nil {
		if cerr := q.createServiceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createServiceStmt: %w", cerr)
		}
	}
	if q.createSettingStmt != nil {
		if cerr := q.createSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSettingStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteAgentByHostnameStmt != nil {
		if cerr := q.deleteAgentByHostnameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAgentByHostnameStmt: %w", cerr)
		}
	}
	if q.deleteAgentByIDStmt != nil {
		if cerr := q.deleteAgentByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAgentByIDStmt: %w", cerr)
		}
	}
	if q.deleteConfigByProfileIDStmt != nil {
		if cerr := q.deleteConfigByProfileIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteConfigByProfileIDStmt: %w", cerr)
		}
	}
	if q.deleteConfigByProfileNameStmt != nil {
		if cerr := q.deleteConfigByProfileNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteConfigByProfileNameStmt: %w", cerr)
		}
	}
	if q.deleteProfileByIDStmt != nil {
		if cerr := q.deleteProfileByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileByIDStmt: %w", cerr)
		}
	}
	if q.deleteProfileByNameStmt != nil {
		if cerr := q.deleteProfileByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileByNameStmt: %w", cerr)
		}
	}
	if q.deleteProviderByIDStmt != nil {
		if cerr := q.deleteProviderByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProviderByIDStmt: %w", cerr)
		}
	}
	if q.deleteProviderByNameStmt != nil {
		if cerr := q.deleteProviderByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProviderByNameStmt: %w", cerr)
		}
	}
	if q.deleteRouterByIDStmt != nil {
		if cerr := q.deleteRouterByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRouterByIDStmt: %w", cerr)
		}
	}
	if q.deleteRouterByNameStmt != nil {
		if cerr := q.deleteRouterByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRouterByNameStmt: %w", cerr)
		}
	}
	if q.deleteServiceByIDStmt != nil {
		if cerr := q.deleteServiceByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteServiceByIDStmt: %w", cerr)
		}
	}
	if q.deleteServiceByNameStmt != nil {
		if cerr := q.deleteServiceByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteServiceByNameStmt: %w", cerr)
		}
	}
	if q.deleteSettingByIDStmt != nil {
		if cerr := q.deleteSettingByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSettingByIDStmt: %w", cerr)
		}
	}
	if q.deleteSettingByKeyStmt != nil {
		if cerr := q.deleteSettingByKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSettingByKeyStmt: %w", cerr)
		}
	}
	if q.deleteUserByIDStmt != nil {
		if cerr := q.deleteUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserByIDStmt: %w", cerr)
		}
	}
	if q.deleteUserByUsernameStmt != nil {
		if cerr := q.deleteUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserByUsernameStmt: %w", cerr)
		}
	}
	if q.getAgentByHostnameStmt != nil {
		if cerr := q.getAgentByHostnameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAgentByHostnameStmt: %w", cerr)
		}
	}
	if q.getAgentByIDStmt != nil {
		if cerr := q.getAgentByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAgentByIDStmt: %w", cerr)
		}
	}
	if q.getConfigByProfileIDStmt != nil {
		if cerr := q.getConfigByProfileIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getConfigByProfileIDStmt: %w", cerr)
		}
	}
	if q.getConfigByProfileNameStmt != nil {
		if cerr := q.getConfigByProfileNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getConfigByProfileNameStmt: %w", cerr)
		}
	}
	if q.getProfileByIDStmt != nil {
		if cerr := q.getProfileByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileByIDStmt: %w", cerr)
		}
	}
	if q.getProfileByNameStmt != nil {
		if cerr := q.getProfileByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileByNameStmt: %w", cerr)
		}
	}
	if q.getProviderByIDStmt != nil {
		if cerr := q.getProviderByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProviderByIDStmt: %w", cerr)
		}
	}
	if q.getProviderByNameStmt != nil {
		if cerr := q.getProviderByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProviderByNameStmt: %w", cerr)
		}
	}
	if q.getRouterByIDStmt != nil {
		if cerr := q.getRouterByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRouterByIDStmt: %w", cerr)
		}
	}
	if q.getServiceByIDStmt != nil {
		if cerr := q.getServiceByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getServiceByIDStmt: %w", cerr)
		}
	}
	if q.getSettingByKeyStmt != nil {
		if cerr := q.getSettingByKeyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingByKeyStmt: %w", cerr)
		}
	}
	if q.getUserByIDStmt != nil {
		if cerr := q.getUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByIDStmt: %w", cerr)
		}
	}
	if q.getUserByUsernameStmt != nil {
		if cerr := q.getUserByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByUsernameStmt: %w", cerr)
		}
	}
	if q.listAgentsStmt != nil {
		if cerr := q.listAgentsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAgentsStmt: %w", cerr)
		}
	}
	if q.listConfigsStmt != nil {
		if cerr := q.listConfigsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listConfigsStmt: %w", cerr)
		}
	}
	if q.listProfilesStmt != nil {
		if cerr := q.listProfilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listProfilesStmt: %w", cerr)
		}
	}
	if q.listProvidersStmt != nil {
		if cerr := q.listProvidersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listProvidersStmt: %w", cerr)
		}
	}
	if q.listRoutersStmt != nil {
		if cerr := q.listRoutersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRoutersStmt: %w", cerr)
		}
	}
	if q.listRoutersByProfileIDStmt != nil {
		if cerr := q.listRoutersByProfileIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRoutersByProfileIDStmt: %w", cerr)
		}
	}
	if q.listServicesStmt != nil {
		if cerr := q.listServicesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listServicesStmt: %w", cerr)
		}
	}
	if q.listServicesByProfileIDStmt != nil {
		if cerr := q.listServicesByProfileIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listServicesByProfileIDStmt: %w", cerr)
		}
	}
	if q.listSettingsStmt != nil {
		if cerr := q.listSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listSettingsStmt: %w", cerr)
		}
	}
	if q.listUsersStmt != nil {
		if cerr := q.listUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUsersStmt: %w", cerr)
		}
	}
	if q.updateAgentStmt != nil {
		if cerr := q.updateAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAgentStmt: %w", cerr)
		}
	}
	if q.updateConfigStmt != nil {
		if cerr := q.updateConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateConfigStmt: %w", cerr)
		}
	}
	if q.updateProfileStmt != nil {
		if cerr := q.updateProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileStmt: %w", cerr)
		}
	}
	if q.updateProviderStmt != nil {
		if cerr := q.updateProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProviderStmt: %w", cerr)
		}
	}
	if q.updateRouterStmt != nil {
		if cerr := q.updateRouterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateRouterStmt: %w", cerr)
		}
	}
	if q.updateServiceStmt != nil {
		if cerr := q.updateServiceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateServiceStmt: %w", cerr)
		}
	}
	if q.updateSettingStmt != nil {
		if cerr := q.updateSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSettingStmt: %w", cerr)
		}
	}
	if q.updateUserStmt != nil {
		if cerr := q.updateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStmt: %w", cerr)
		}
	}
	if q.upsertAgentStmt != nil {
		if cerr := q.upsertAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertAgentStmt: %w", cerr)
		}
	}
	if q.upsertProfileStmt != nil {
		if cerr := q.upsertProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertProfileStmt: %w", cerr)
		}
	}
	if q.upsertProviderStmt != nil {
		if cerr := q.upsertProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertProviderStmt: %w", cerr)
		}
	}
	if q.upsertRouterStmt != nil {
		if cerr := q.upsertRouterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertRouterStmt: %w", cerr)
		}
	}
	if q.upsertServiceStmt != nil {
		if cerr := q.upsertServiceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertServiceStmt: %w", cerr)
		}
	}
	if q.upsertSettingStmt != nil {
		if cerr := q.upsertSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertSettingStmt: %w", cerr)
		}
	}
	if q.upsertUserStmt != nil {
		if cerr := q.upsertUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertUserStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                            DBTX
	tx                            *sql.Tx
	createAgentStmt               *sql.Stmt
	createConfigStmt              *sql.Stmt
	createProfileStmt             *sql.Stmt
	createProviderStmt            *sql.Stmt
	createRouterStmt              *sql.Stmt
	createServiceStmt             *sql.Stmt
	createSettingStmt             *sql.Stmt
	createUserStmt                *sql.Stmt
	deleteAgentByHostnameStmt     *sql.Stmt
	deleteAgentByIDStmt           *sql.Stmt
	deleteConfigByProfileIDStmt   *sql.Stmt
	deleteConfigByProfileNameStmt *sql.Stmt
	deleteProfileByIDStmt         *sql.Stmt
	deleteProfileByNameStmt       *sql.Stmt
	deleteProviderByIDStmt        *sql.Stmt
	deleteProviderByNameStmt      *sql.Stmt
	deleteRouterByIDStmt          *sql.Stmt
	deleteRouterByNameStmt        *sql.Stmt
	deleteServiceByIDStmt         *sql.Stmt
	deleteServiceByNameStmt       *sql.Stmt
	deleteSettingByIDStmt         *sql.Stmt
	deleteSettingByKeyStmt        *sql.Stmt
	deleteUserByIDStmt            *sql.Stmt
	deleteUserByUsernameStmt      *sql.Stmt
	getAgentByHostnameStmt        *sql.Stmt
	getAgentByIDStmt              *sql.Stmt
	getConfigByProfileIDStmt      *sql.Stmt
	getConfigByProfileNameStmt    *sql.Stmt
	getProfileByIDStmt            *sql.Stmt
	getProfileByNameStmt          *sql.Stmt
	getProviderByIDStmt           *sql.Stmt
	getProviderByNameStmt         *sql.Stmt
	getRouterByIDStmt             *sql.Stmt
	getServiceByIDStmt            *sql.Stmt
	getSettingByKeyStmt           *sql.Stmt
	getUserByIDStmt               *sql.Stmt
	getUserByUsernameStmt         *sql.Stmt
	listAgentsStmt                *sql.Stmt
	listConfigsStmt               *sql.Stmt
	listProfilesStmt              *sql.Stmt
	listProvidersStmt             *sql.Stmt
	listRoutersStmt               *sql.Stmt
	listRoutersByProfileIDStmt    *sql.Stmt
	listServicesStmt              *sql.Stmt
	listServicesByProfileIDStmt   *sql.Stmt
	listSettingsStmt              *sql.Stmt
	listUsersStmt                 *sql.Stmt
	updateAgentStmt               *sql.Stmt
	updateConfigStmt              *sql.Stmt
	updateProfileStmt             *sql.Stmt
	updateProviderStmt            *sql.Stmt
	updateRouterStmt              *sql.Stmt
	updateServiceStmt             *sql.Stmt
	updateSettingStmt             *sql.Stmt
	updateUserStmt                *sql.Stmt
	upsertAgentStmt               *sql.Stmt
	upsertProfileStmt             *sql.Stmt
	upsertProviderStmt            *sql.Stmt
	upsertRouterStmt              *sql.Stmt
	upsertServiceStmt             *sql.Stmt
	upsertSettingStmt             *sql.Stmt
	upsertUserStmt                *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                            tx,
		tx:                            tx,
		createAgentStmt:               q.createAgentStmt,
		createConfigStmt:              q.createConfigStmt,
		createProfileStmt:             q.createProfileStmt,
		createProviderStmt:            q.createProviderStmt,
		createRouterStmt:              q.createRouterStmt,
		createServiceStmt:             q.createServiceStmt,
		createSettingStmt:             q.createSettingStmt,
		createUserStmt:                q.createUserStmt,
		deleteAgentByHostnameStmt:     q.deleteAgentByHostnameStmt,
		deleteAgentByIDStmt:           q.deleteAgentByIDStmt,
		deleteConfigByProfileIDStmt:   q.deleteConfigByProfileIDStmt,
		deleteConfigByProfileNameStmt: q.deleteConfigByProfileNameStmt,
		deleteProfileByIDStmt:         q.deleteProfileByIDStmt,
		deleteProfileByNameStmt:       q.deleteProfileByNameStmt,
		deleteProviderByIDStmt:        q.deleteProviderByIDStmt,
		deleteProviderByNameStmt:      q.deleteProviderByNameStmt,
		deleteRouterByIDStmt:          q.deleteRouterByIDStmt,
		deleteRouterByNameStmt:        q.deleteRouterByNameStmt,
		deleteServiceByIDStmt:         q.deleteServiceByIDStmt,
		deleteServiceByNameStmt:       q.deleteServiceByNameStmt,
		deleteSettingByIDStmt:         q.deleteSettingByIDStmt,
		deleteSettingByKeyStmt:        q.deleteSettingByKeyStmt,
		deleteUserByIDStmt:            q.deleteUserByIDStmt,
		deleteUserByUsernameStmt:      q.deleteUserByUsernameStmt,
		getAgentByHostnameStmt:        q.getAgentByHostnameStmt,
		getAgentByIDStmt:              q.getAgentByIDStmt,
		getConfigByProfileIDStmt:      q.getConfigByProfileIDStmt,
		getConfigByProfileNameStmt:    q.getConfigByProfileNameStmt,
		getProfileByIDStmt:            q.getProfileByIDStmt,
		getProfileByNameStmt:          q.getProfileByNameStmt,
		getProviderByIDStmt:           q.getProviderByIDStmt,
		getProviderByNameStmt:         q.getProviderByNameStmt,
		getRouterByIDStmt:             q.getRouterByIDStmt,
		getServiceByIDStmt:            q.getServiceByIDStmt,
		getSettingByKeyStmt:           q.getSettingByKeyStmt,
		getUserByIDStmt:               q.getUserByIDStmt,
		getUserByUsernameStmt:         q.getUserByUsernameStmt,
		listAgentsStmt:                q.listAgentsStmt,
		listConfigsStmt:               q.listConfigsStmt,
		listProfilesStmt:              q.listProfilesStmt,
		listProvidersStmt:             q.listProvidersStmt,
		listRoutersStmt:               q.listRoutersStmt,
		listRoutersByProfileIDStmt:    q.listRoutersByProfileIDStmt,
		listServicesStmt:              q.listServicesStmt,
		listServicesByProfileIDStmt:   q.listServicesByProfileIDStmt,
		listSettingsStmt:              q.listSettingsStmt,
		listUsersStmt:                 q.listUsersStmt,
		updateAgentStmt:               q.updateAgentStmt,
		updateConfigStmt:              q.updateConfigStmt,
		updateProfileStmt:             q.updateProfileStmt,
		updateProviderStmt:            q.updateProviderStmt,
		updateRouterStmt:              q.updateRouterStmt,
		updateServiceStmt:             q.updateServiceStmt,
		updateSettingStmt:             q.updateSettingStmt,
		updateUserStmt:                q.updateUserStmt,
		upsertAgentStmt:               q.upsertAgentStmt,
		upsertProfileStmt:             q.upsertProfileStmt,
		upsertProviderStmt:            q.upsertProviderStmt,
		upsertRouterStmt:              q.upsertRouterStmt,
		upsertServiceStmt:             q.upsertServiceStmt,
		upsertSettingStmt:             q.upsertSettingStmt,
		upsertUserStmt:                q.upsertUserStmt,
	}
}
