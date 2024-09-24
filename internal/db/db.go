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
	if q.createConfigStmt, err = db.PrepareContext(ctx, createConfig); err != nil {
		return nil, fmt.Errorf("error preparing query CreateConfig: %w", err)
	}
	if q.createProfileStmt, err = db.PrepareContext(ctx, createProfile); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProfile: %w", err)
	}
	if q.createProviderStmt, err = db.PrepareContext(ctx, createProvider); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProvider: %w", err)
	}
	if q.createSettingStmt, err = db.PrepareContext(ctx, createSetting); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSetting: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
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
	if q.deleteSettingByKeyStmt, err = db.PrepareContext(ctx, deleteSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSettingByKey: %w", err)
	}
	if q.deleteUserByIDStmt, err = db.PrepareContext(ctx, deleteUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserByID: %w", err)
	}
	if q.deleteUserByUsernameStmt, err = db.PrepareContext(ctx, deleteUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserByUsername: %w", err)
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
	if q.getSettingByKeyStmt, err = db.PrepareContext(ctx, getSettingByKey); err != nil {
		return nil, fmt.Errorf("error preparing query GetSettingByKey: %w", err)
	}
	if q.getUserByIDStmt, err = db.PrepareContext(ctx, getUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByID: %w", err)
	}
	if q.getUserByUsernameStmt, err = db.PrepareContext(ctx, getUserByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByUsername: %w", err)
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
	if q.listSettingsStmt, err = db.PrepareContext(ctx, listSettings); err != nil {
		return nil, fmt.Errorf("error preparing query ListSettings: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
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
	if q.updateSettingStmt, err = db.PrepareContext(ctx, updateSetting); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSetting: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	if q.upsertProfileStmt, err = db.PrepareContext(ctx, upsertProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertProfile: %w", err)
	}
	if q.upsertProviderStmt, err = db.PrepareContext(ctx, upsertProvider); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertProvider: %w", err)
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
	createConfigStmt              *sql.Stmt
	createProfileStmt             *sql.Stmt
	createProviderStmt            *sql.Stmt
	createSettingStmt             *sql.Stmt
	createUserStmt                *sql.Stmt
	deleteConfigByProfileIDStmt   *sql.Stmt
	deleteConfigByProfileNameStmt *sql.Stmt
	deleteProfileByIDStmt         *sql.Stmt
	deleteProfileByNameStmt       *sql.Stmt
	deleteProviderByIDStmt        *sql.Stmt
	deleteProviderByNameStmt      *sql.Stmt
	deleteSettingByKeyStmt        *sql.Stmt
	deleteUserByIDStmt            *sql.Stmt
	deleteUserByUsernameStmt      *sql.Stmt
	getConfigByProfileIDStmt      *sql.Stmt
	getConfigByProfileNameStmt    *sql.Stmt
	getProfileByIDStmt            *sql.Stmt
	getProfileByNameStmt          *sql.Stmt
	getProviderByIDStmt           *sql.Stmt
	getProviderByNameStmt         *sql.Stmt
	getSettingByKeyStmt           *sql.Stmt
	getUserByIDStmt               *sql.Stmt
	getUserByUsernameStmt         *sql.Stmt
	listConfigsStmt               *sql.Stmt
	listProfilesStmt              *sql.Stmt
	listProvidersStmt             *sql.Stmt
	listSettingsStmt              *sql.Stmt
	listUsersStmt                 *sql.Stmt
	updateConfigStmt              *sql.Stmt
	updateProfileStmt             *sql.Stmt
	updateProviderStmt            *sql.Stmt
	updateSettingStmt             *sql.Stmt
	updateUserStmt                *sql.Stmt
	upsertProfileStmt             *sql.Stmt
	upsertProviderStmt            *sql.Stmt
	upsertSettingStmt             *sql.Stmt
	upsertUserStmt                *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                            tx,
		tx:                            tx,
		createConfigStmt:              q.createConfigStmt,
		createProfileStmt:             q.createProfileStmt,
		createProviderStmt:            q.createProviderStmt,
		createSettingStmt:             q.createSettingStmt,
		createUserStmt:                q.createUserStmt,
		deleteConfigByProfileIDStmt:   q.deleteConfigByProfileIDStmt,
		deleteConfigByProfileNameStmt: q.deleteConfigByProfileNameStmt,
		deleteProfileByIDStmt:         q.deleteProfileByIDStmt,
		deleteProfileByNameStmt:       q.deleteProfileByNameStmt,
		deleteProviderByIDStmt:        q.deleteProviderByIDStmt,
		deleteProviderByNameStmt:      q.deleteProviderByNameStmt,
		deleteSettingByKeyStmt:        q.deleteSettingByKeyStmt,
		deleteUserByIDStmt:            q.deleteUserByIDStmt,
		deleteUserByUsernameStmt:      q.deleteUserByUsernameStmt,
		getConfigByProfileIDStmt:      q.getConfigByProfileIDStmt,
		getConfigByProfileNameStmt:    q.getConfigByProfileNameStmt,
		getProfileByIDStmt:            q.getProfileByIDStmt,
		getProfileByNameStmt:          q.getProfileByNameStmt,
		getProviderByIDStmt:           q.getProviderByIDStmt,
		getProviderByNameStmt:         q.getProviderByNameStmt,
		getSettingByKeyStmt:           q.getSettingByKeyStmt,
		getUserByIDStmt:               q.getUserByIDStmt,
		getUserByUsernameStmt:         q.getUserByUsernameStmt,
		listConfigsStmt:               q.listConfigsStmt,
		listProfilesStmt:              q.listProfilesStmt,
		listProvidersStmt:             q.listProvidersStmt,
		listSettingsStmt:              q.listSettingsStmt,
		listUsersStmt:                 q.listUsersStmt,
		updateConfigStmt:              q.updateConfigStmt,
		updateProfileStmt:             q.updateProfileStmt,
		updateProviderStmt:            q.updateProviderStmt,
		updateSettingStmt:             q.updateSettingStmt,
		updateUserStmt:                q.updateUserStmt,
		upsertProfileStmt:             q.upsertProfileStmt,
		upsertProviderStmt:            q.upsertProviderStmt,
		upsertSettingStmt:             q.upsertSettingStmt,
		upsertUserStmt:                q.upsertUserStmt,
	}
}
