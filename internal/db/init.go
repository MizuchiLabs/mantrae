package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

var (
	DB    *sql.DB
	Query *Queries
)

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
		dsn := fmt.Sprintf("file:%s?mode=rwc&_journal=WAL&_fk=1&_sync=NORMAL", util.Path(util.DBName))
		db, err = sql.Open("sqlite", dsn)
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
