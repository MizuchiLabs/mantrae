// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

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
	if q.createDNSProviderStmt, err = db.PrepareContext(ctx, createDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query CreateDNSProvider: %w", err)
	}
	if q.createProfileStmt, err = db.PrepareContext(ctx, createProfile); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProfile: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteAgentStmt, err = db.PrepareContext(ctx, deleteAgent); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAgent: %w", err)
	}
	if q.deleteDNSProviderStmt, err = db.PrepareContext(ctx, deleteDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteDNSProvider: %w", err)
	}
	if q.deleteProfileStmt, err = db.PrepareContext(ctx, deleteProfile); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfile: %w", err)
	}
	if q.deleteRouterDNSProviderStmt, err = db.PrepareContext(ctx, deleteRouterDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRouterDNSProvider: %w", err)
	}
	if q.deleteSettingStmt, err = db.PrepareContext(ctx, deleteSetting); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSetting: %w", err)
	}
	if q.deleteTraefikConfigStmt, err = db.PrepareContext(ctx, deleteTraefikConfig); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTraefikConfig: %w", err)
	}
	if q.deleteTraefikConfigByAgentStmt, err = db.PrepareContext(ctx, deleteTraefikConfigByAgent); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTraefikConfigByAgent: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.getAPITraefikConfigStmt, err = db.PrepareContext(ctx, getAPITraefikConfig); err != nil {
		return nil, fmt.Errorf("error preparing query GetAPITraefikConfig: %w", err)
	}
	if q.getActiveDNSProviderStmt, err = db.PrepareContext(ctx, getActiveDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query GetActiveDNSProvider: %w", err)
	}
	if q.getAgentStmt, err = db.PrepareContext(ctx, getAgent); err != nil {
		return nil, fmt.Errorf("error preparing query GetAgent: %w", err)
	}
	if q.getAgentTraefikConfigsStmt, err = db.PrepareContext(ctx, getAgentTraefikConfigs); err != nil {
		return nil, fmt.Errorf("error preparing query GetAgentTraefikConfigs: %w", err)
	}
	if q.getDNSProviderStmt, err = db.PrepareContext(ctx, getDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query GetDNSProvider: %w", err)
	}
	if q.getLocalTraefikConfigStmt, err = db.PrepareContext(ctx, getLocalTraefikConfig); err != nil {
		return nil, fmt.Errorf("error preparing query GetLocalTraefikConfig: %w", err)
	}
	if q.getProfileStmt, err = db.PrepareContext(ctx, getProfile); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfile: %w", err)
	}
	if q.getProfileByNameStmt, err = db.PrepareContext(ctx, getProfileByName); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileByName: %w", err)
	}
	if q.getRouterDNSProviderStmt, err = db.PrepareContext(ctx, getRouterDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query GetRouterDNSProvider: %w", err)
	}
	if q.getSettingStmt, err = db.PrepareContext(ctx, getSetting); err != nil {
		return nil, fmt.Errorf("error preparing query GetSetting: %w", err)
	}
	if q.getTraefikConfigByIDStmt, err = db.PrepareContext(ctx, getTraefikConfigByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetTraefikConfigByID: %w", err)
	}
	if q.getTraefikConfigBySourceStmt, err = db.PrepareContext(ctx, getTraefikConfigBySource); err != nil {
		return nil, fmt.Errorf("error preparing query GetTraefikConfigBySource: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
	}
	if q.listAgentsStmt, err = db.PrepareContext(ctx, listAgents); err != nil {
		return nil, fmt.Errorf("error preparing query ListAgents: %w", err)
	}
	if q.listAgentsByProfileStmt, err = db.PrepareContext(ctx, listAgentsByProfile); err != nil {
		return nil, fmt.Errorf("error preparing query ListAgentsByProfile: %w", err)
	}
	if q.listDNSProvidersStmt, err = db.PrepareContext(ctx, listDNSProviders); err != nil {
		return nil, fmt.Errorf("error preparing query ListDNSProviders: %w", err)
	}
	if q.listProfilesStmt, err = db.PrepareContext(ctx, listProfiles); err != nil {
		return nil, fmt.Errorf("error preparing query ListProfiles: %w", err)
	}
	if q.listRouterDNSProvidersByTraefikIDStmt, err = db.PrepareContext(ctx, listRouterDNSProvidersByTraefikID); err != nil {
		return nil, fmt.Errorf("error preparing query ListRouterDNSProvidersByTraefikID: %w", err)
	}
	if q.listSettingsStmt, err = db.PrepareContext(ctx, listSettings); err != nil {
		return nil, fmt.Errorf("error preparing query ListSettings: %w", err)
	}
	if q.listTraefikIDsStmt, err = db.PrepareContext(ctx, listTraefikIDs); err != nil {
		return nil, fmt.Errorf("error preparing query ListTraefikIDs: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
	}
	if q.updateAgentStmt, err = db.PrepareContext(ctx, updateAgent); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAgent: %w", err)
	}
	if q.updateAgentIPStmt, err = db.PrepareContext(ctx, updateAgentIP); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAgentIP: %w", err)
	}
	if q.updateAgentTokenStmt, err = db.PrepareContext(ctx, updateAgentToken); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAgentToken: %w", err)
	}
	if q.updateDNSProviderStmt, err = db.PrepareContext(ctx, updateDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateDNSProvider: %w", err)
	}
	if q.updateProfileStmt, err = db.PrepareContext(ctx, updateProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfile: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	if q.updateUserLastLoginStmt, err = db.PrepareContext(ctx, updateUserLastLogin); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserLastLogin: %w", err)
	}
	if q.updateUserPasswordStmt, err = db.PrepareContext(ctx, updateUserPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserPassword: %w", err)
	}
	if q.updateUserResetTokenStmt, err = db.PrepareContext(ctx, updateUserResetToken); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserResetToken: %w", err)
	}
	if q.upsertRouterDNSProviderStmt, err = db.PrepareContext(ctx, upsertRouterDNSProvider); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertRouterDNSProvider: %w", err)
	}
	if q.upsertSettingStmt, err = db.PrepareContext(ctx, upsertSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertSetting: %w", err)
	}
	if q.upsertTraefikAgentConfigStmt, err = db.PrepareContext(ctx, upsertTraefikAgentConfig); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertTraefikAgentConfig: %w", err)
	}
	if q.upsertTraefikConfigStmt, err = db.PrepareContext(ctx, upsertTraefikConfig); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertTraefikConfig: %w", err)
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
	if q.createDNSProviderStmt != nil {
		if cerr := q.createDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createDNSProviderStmt: %w", cerr)
		}
	}
	if q.createProfileStmt != nil {
		if cerr := q.createProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProfileStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteAgentStmt != nil {
		if cerr := q.deleteAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAgentStmt: %w", cerr)
		}
	}
	if q.deleteDNSProviderStmt != nil {
		if cerr := q.deleteDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteDNSProviderStmt: %w", cerr)
		}
	}
	if q.deleteProfileStmt != nil {
		if cerr := q.deleteProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileStmt: %w", cerr)
		}
	}
	if q.deleteRouterDNSProviderStmt != nil {
		if cerr := q.deleteRouterDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRouterDNSProviderStmt: %w", cerr)
		}
	}
	if q.deleteSettingStmt != nil {
		if cerr := q.deleteSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSettingStmt: %w", cerr)
		}
	}
	if q.deleteTraefikConfigStmt != nil {
		if cerr := q.deleteTraefikConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTraefikConfigStmt: %w", cerr)
		}
	}
	if q.deleteTraefikConfigByAgentStmt != nil {
		if cerr := q.deleteTraefikConfigByAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTraefikConfigByAgentStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.getAPITraefikConfigStmt != nil {
		if cerr := q.getAPITraefikConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAPITraefikConfigStmt: %w", cerr)
		}
	}
	if q.getActiveDNSProviderStmt != nil {
		if cerr := q.getActiveDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getActiveDNSProviderStmt: %w", cerr)
		}
	}
	if q.getAgentStmt != nil {
		if cerr := q.getAgentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAgentStmt: %w", cerr)
		}
	}
	if q.getAgentTraefikConfigsStmt != nil {
		if cerr := q.getAgentTraefikConfigsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAgentTraefikConfigsStmt: %w", cerr)
		}
	}
	if q.getDNSProviderStmt != nil {
		if cerr := q.getDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDNSProviderStmt: %w", cerr)
		}
	}
	if q.getLocalTraefikConfigStmt != nil {
		if cerr := q.getLocalTraefikConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLocalTraefikConfigStmt: %w", cerr)
		}
	}
	if q.getProfileStmt != nil {
		if cerr := q.getProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileStmt: %w", cerr)
		}
	}
	if q.getProfileByNameStmt != nil {
		if cerr := q.getProfileByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileByNameStmt: %w", cerr)
		}
	}
	if q.getRouterDNSProviderStmt != nil {
		if cerr := q.getRouterDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRouterDNSProviderStmt: %w", cerr)
		}
	}
	if q.getSettingStmt != nil {
		if cerr := q.getSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSettingStmt: %w", cerr)
		}
	}
	if q.getTraefikConfigByIDStmt != nil {
		if cerr := q.getTraefikConfigByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTraefikConfigByIDStmt: %w", cerr)
		}
	}
	if q.getTraefikConfigBySourceStmt != nil {
		if cerr := q.getTraefikConfigBySourceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTraefikConfigBySourceStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
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
	if q.listAgentsByProfileStmt != nil {
		if cerr := q.listAgentsByProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAgentsByProfileStmt: %w", cerr)
		}
	}
	if q.listDNSProvidersStmt != nil {
		if cerr := q.listDNSProvidersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listDNSProvidersStmt: %w", cerr)
		}
	}
	if q.listProfilesStmt != nil {
		if cerr := q.listProfilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listProfilesStmt: %w", cerr)
		}
	}
	if q.listRouterDNSProvidersByTraefikIDStmt != nil {
		if cerr := q.listRouterDNSProvidersByTraefikIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRouterDNSProvidersByTraefikIDStmt: %w", cerr)
		}
	}
	if q.listSettingsStmt != nil {
		if cerr := q.listSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listSettingsStmt: %w", cerr)
		}
	}
	if q.listTraefikIDsStmt != nil {
		if cerr := q.listTraefikIDsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listTraefikIDsStmt: %w", cerr)
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
	if q.updateAgentIPStmt != nil {
		if cerr := q.updateAgentIPStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAgentIPStmt: %w", cerr)
		}
	}
	if q.updateAgentTokenStmt != nil {
		if cerr := q.updateAgentTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAgentTokenStmt: %w", cerr)
		}
	}
	if q.updateDNSProviderStmt != nil {
		if cerr := q.updateDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateDNSProviderStmt: %w", cerr)
		}
	}
	if q.updateProfileStmt != nil {
		if cerr := q.updateProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileStmt: %w", cerr)
		}
	}
	if q.updateUserStmt != nil {
		if cerr := q.updateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStmt: %w", cerr)
		}
	}
	if q.updateUserLastLoginStmt != nil {
		if cerr := q.updateUserLastLoginStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserLastLoginStmt: %w", cerr)
		}
	}
	if q.updateUserPasswordStmt != nil {
		if cerr := q.updateUserPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserPasswordStmt: %w", cerr)
		}
	}
	if q.updateUserResetTokenStmt != nil {
		if cerr := q.updateUserResetTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserResetTokenStmt: %w", cerr)
		}
	}
	if q.upsertRouterDNSProviderStmt != nil {
		if cerr := q.upsertRouterDNSProviderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertRouterDNSProviderStmt: %w", cerr)
		}
	}
	if q.upsertSettingStmt != nil {
		if cerr := q.upsertSettingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertSettingStmt: %w", cerr)
		}
	}
	if q.upsertTraefikAgentConfigStmt != nil {
		if cerr := q.upsertTraefikAgentConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertTraefikAgentConfigStmt: %w", cerr)
		}
	}
	if q.upsertTraefikConfigStmt != nil {
		if cerr := q.upsertTraefikConfigStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertTraefikConfigStmt: %w", cerr)
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
	db                                    DBTX
	tx                                    *sql.Tx
	createAgentStmt                       *sql.Stmt
	createDNSProviderStmt                 *sql.Stmt
	createProfileStmt                     *sql.Stmt
	createUserStmt                        *sql.Stmt
	deleteAgentStmt                       *sql.Stmt
	deleteDNSProviderStmt                 *sql.Stmt
	deleteProfileStmt                     *sql.Stmt
	deleteRouterDNSProviderStmt           *sql.Stmt
	deleteSettingStmt                     *sql.Stmt
	deleteTraefikConfigStmt               *sql.Stmt
	deleteTraefikConfigByAgentStmt        *sql.Stmt
	deleteUserStmt                        *sql.Stmt
	getAPITraefikConfigStmt               *sql.Stmt
	getActiveDNSProviderStmt              *sql.Stmt
	getAgentStmt                          *sql.Stmt
	getAgentTraefikConfigsStmt            *sql.Stmt
	getDNSProviderStmt                    *sql.Stmt
	getLocalTraefikConfigStmt             *sql.Stmt
	getProfileStmt                        *sql.Stmt
	getProfileByNameStmt                  *sql.Stmt
	getRouterDNSProviderStmt              *sql.Stmt
	getSettingStmt                        *sql.Stmt
	getTraefikConfigByIDStmt              *sql.Stmt
	getTraefikConfigBySourceStmt          *sql.Stmt
	getUserStmt                           *sql.Stmt
	getUserByUsernameStmt                 *sql.Stmt
	listAgentsStmt                        *sql.Stmt
	listAgentsByProfileStmt               *sql.Stmt
	listDNSProvidersStmt                  *sql.Stmt
	listProfilesStmt                      *sql.Stmt
	listRouterDNSProvidersByTraefikIDStmt *sql.Stmt
	listSettingsStmt                      *sql.Stmt
	listTraefikIDsStmt                    *sql.Stmt
	listUsersStmt                         *sql.Stmt
	updateAgentStmt                       *sql.Stmt
	updateAgentIPStmt                     *sql.Stmt
	updateAgentTokenStmt                  *sql.Stmt
	updateDNSProviderStmt                 *sql.Stmt
	updateProfileStmt                     *sql.Stmt
	updateUserStmt                        *sql.Stmt
	updateUserLastLoginStmt               *sql.Stmt
	updateUserPasswordStmt                *sql.Stmt
	updateUserResetTokenStmt              *sql.Stmt
	upsertRouterDNSProviderStmt           *sql.Stmt
	upsertSettingStmt                     *sql.Stmt
	upsertTraefikAgentConfigStmt          *sql.Stmt
	upsertTraefikConfigStmt               *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                    tx,
		tx:                                    tx,
		createAgentStmt:                       q.createAgentStmt,
		createDNSProviderStmt:                 q.createDNSProviderStmt,
		createProfileStmt:                     q.createProfileStmt,
		createUserStmt:                        q.createUserStmt,
		deleteAgentStmt:                       q.deleteAgentStmt,
		deleteDNSProviderStmt:                 q.deleteDNSProviderStmt,
		deleteProfileStmt:                     q.deleteProfileStmt,
		deleteRouterDNSProviderStmt:           q.deleteRouterDNSProviderStmt,
		deleteSettingStmt:                     q.deleteSettingStmt,
		deleteTraefikConfigStmt:               q.deleteTraefikConfigStmt,
		deleteTraefikConfigByAgentStmt:        q.deleteTraefikConfigByAgentStmt,
		deleteUserStmt:                        q.deleteUserStmt,
		getAPITraefikConfigStmt:               q.getAPITraefikConfigStmt,
		getActiveDNSProviderStmt:              q.getActiveDNSProviderStmt,
		getAgentStmt:                          q.getAgentStmt,
		getAgentTraefikConfigsStmt:            q.getAgentTraefikConfigsStmt,
		getDNSProviderStmt:                    q.getDNSProviderStmt,
		getLocalTraefikConfigStmt:             q.getLocalTraefikConfigStmt,
		getProfileStmt:                        q.getProfileStmt,
		getProfileByNameStmt:                  q.getProfileByNameStmt,
		getRouterDNSProviderStmt:              q.getRouterDNSProviderStmt,
		getSettingStmt:                        q.getSettingStmt,
		getTraefikConfigByIDStmt:              q.getTraefikConfigByIDStmt,
		getTraefikConfigBySourceStmt:          q.getTraefikConfigBySourceStmt,
		getUserStmt:                           q.getUserStmt,
		getUserByUsernameStmt:                 q.getUserByUsernameStmt,
		listAgentsStmt:                        q.listAgentsStmt,
		listAgentsByProfileStmt:               q.listAgentsByProfileStmt,
		listDNSProvidersStmt:                  q.listDNSProvidersStmt,
		listProfilesStmt:                      q.listProfilesStmt,
		listRouterDNSProvidersByTraefikIDStmt: q.listRouterDNSProvidersByTraefikIDStmt,
		listSettingsStmt:                      q.listSettingsStmt,
		listTraefikIDsStmt:                    q.listTraefikIDsStmt,
		listUsersStmt:                         q.listUsersStmt,
		updateAgentStmt:                       q.updateAgentStmt,
		updateAgentIPStmt:                     q.updateAgentIPStmt,
		updateAgentTokenStmt:                  q.updateAgentTokenStmt,
		updateDNSProviderStmt:                 q.updateDNSProviderStmt,
		updateProfileStmt:                     q.updateProfileStmt,
		updateUserStmt:                        q.updateUserStmt,
		updateUserLastLoginStmt:               q.updateUserLastLoginStmt,
		updateUserPasswordStmt:                q.updateUserPasswordStmt,
		updateUserResetTokenStmt:              q.updateUserResetTokenStmt,
		upsertRouterDNSProviderStmt:           q.upsertRouterDNSProviderStmt,
		upsertSettingStmt:                     q.upsertSettingStmt,
		upsertTraefikAgentConfigStmt:          q.upsertTraefikAgentConfigStmt,
		upsertTraefikConfigStmt:               q.upsertTraefikConfigStmt,
	}
}
