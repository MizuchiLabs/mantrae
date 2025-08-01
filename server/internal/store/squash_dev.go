//go:build dev
// +build dev

package store

import (
	"fmt"
	"log/slog"
	"os"
)

func Squash() {
	conn := NewConnection(":memory:")
	conn.Migrate()

	db := conn.db
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	var currentVersion int64
	err := db.QueryRow("SELECT version_id FROM goose_db_version ORDER BY id DESC LIMIT 1").
		Scan(&currentVersion)
	if err != nil {
		panic(err)
	}

	// Dump the complete schema
	rows, err := db.Query(`
		SELECT sql FROM sqlite_master 
		WHERE sql IS NOT NULL 
		AND type = 'table' 
		AND name NOT LIKE 'sqlite_%'
		AND name != 'goose_db_version'
	`)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	// Create new base migration
	baseFile := "-- +goose Up\n"
	for rows.Next() {
		var sql string
		if err = rows.Scan(&sql); err != nil {
			panic(err)
		}
		baseFile += sql + ";\n\n"
	}

	// Add indexes and triggers
	rows, err = db.Query(`
		SELECT type, sql FROM sqlite_master 
		WHERE sql IS NOT NULL 
		AND (type = 'index' OR type = 'trigger')
		AND name NOT LIKE 'sqlite_%'
	`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var sqlType, sql string
		if err = rows.Scan(&sqlType, &sql); err != nil {
			panic(err)
		}

		// Wrap triggers in StatementBegin/End blocks
		if sqlType == "trigger" {
			baseFile += "-- +goose StatementBegin\n"
			baseFile += sql + ";\n"
			baseFile += "-- +goose StatementEnd\n\n"
		} else {
			// Regular indexes don't need the statement blocks
			baseFile += sql + ";\n\n"
		}
	}
	// Add goose Down section
	baseFile += "-- +goose Down\n"
	rows, err = db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type = 'table' 
		AND name NOT LIKE 'sqlite_%'
		AND name != 'goose_db_version'
	`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			panic(err)
		}
		baseFile += fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", name)
	}

	// Write new reference
	err = os.WriteFile("00000_base.sql", []byte(baseFile), 0644)
	if err != nil {
		panic(err)
	}

	slog.Info("Migrations squashed successfully!")
}
