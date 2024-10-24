package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

var (
	DB    *sql.DB
	Query *Queries
)

// Default settings for sqlite db
const initScript = `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA foreign_keys = ON;
	PRAGMA mmap_size = 300000000;
	PRAGMA page_size = 32768;
	PRAGMA temp_store = MEMORY;
`

// isTest returns true if the current program is running in a test environment
func isTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func InitDB() error {
	var db *sql.DB
	var err error

	if isTest() {
		db, err = sql.Open("sqlite", "file:test.db?mode=memory")
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
	} else {
		sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, _ string) error {
			_, err = conn.ExecContext(context.Background(), initScript, nil)
			return err
		})

		db, err = sql.Open("sqlite", "file:"+util.Path(util.DBName))
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
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
