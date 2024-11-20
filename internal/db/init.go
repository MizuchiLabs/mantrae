package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
	//_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

var (
	DB    *sql.DB
	Query *Queries
)

// Default settings for sqlite db
const initScript = `
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

func InitDB() error {
	var db *sql.DB
	var err error

	if util.IsTest() {
		db, err = sql.Open("sqlite", "file::memory:?cache=shared")
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
	} else {
		sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, _ string) error {
			_, err = conn.ExecContext(context.Background(), initScript, nil)
			return err
		})

		db, err = sql.Open("sqlite", "file:"+util.DBPath())
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		db.SetMaxOpenConns(1) // Only one writer
		db.SetMaxIdleConns(1) // Prevent idle connections
	}

	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	Query = New(db)
	return nil
}
