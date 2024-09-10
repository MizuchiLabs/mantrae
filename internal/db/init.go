package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

var (
	DB    *sql.DB
	Query *Queries
)

func InitDB() error {
	db, err := sql.Open("sqlite3", "file:mantrae.db?mode=rwc&_journal=WAL&_fk=1&_sync=NORMAL")
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to open database: %w", err)
	}

	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		db.Close()
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	Query = New(db)
	return nil
}
