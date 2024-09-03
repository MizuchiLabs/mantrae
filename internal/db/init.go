package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

var Query *Queries

func InitDB() (*sql.DB, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:mantrae.db?mode=rwc&_journal=WAL&_fk=1&_sync=NORMAL")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check if the database is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check database: %w", err)
	}

	if count == 0 {
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	Query = New(db)
	return db, nil
}
