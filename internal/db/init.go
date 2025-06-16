package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/mizuchilabs/mantrae/internal/app"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Default settings for sqlite db
const InitScript = `
	PRAGMA busy_timeout = 5000;
	PRAGMA journal_mode = WAL;
	PRAGMA journal_size_limit = 200000000;
	PRAGMA synchronous = NORMAL;
	PRAGMA foreign_keys = ON;
	PRAGMA temp_store = MEMORY;
	PRAGMA mmap_size = 300000000;
	PRAGMA page_size = 32768;
	PRAGMA cache_size = -16000;
`

func InitDB() (*sql.DB, error) {
	var db *sql.DB
	var err error

	if util.IsTest() {
		db, err = sql.Open("sqlite", "file::memory:?cache=shared")
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
	} else {
		sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, _ string) error {
			_, err = conn.ExecContext(context.Background(), InitScript, nil)
			return err
		})

		dbPath := app.ResolvePath("mantrae.db")
		db, err = sql.Open("sqlite", fmt.Sprintf("file:%s", dbPath))
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		db.SetMaxOpenConns(1)    // Only one writer
		db.SetMaxIdleConns(1)    // Prevent idle connections
		db.SetConnMaxLifetime(0) // Connections are reused forever
	}

	if err := Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
