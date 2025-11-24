// Package store provides functionality for interacting with the database.
package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

type Connection struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewConnection opens a SQLite connection.
// If `path` is empty, it defaults to "mantrae.db" in the data dir.
// If `path` is ":memory:" or "file::memory:?cache=shared", opens in-memory.
func NewConnection(path string) *Connection {
	dataSource := fmt.Sprintf("file:%s", filepath.ToSlash(util.ResolvePath("mantrae.db")))
	if path != "" {
		dataSource = path
	}

	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := configureSQLite(db); err != nil {
		log.Fatalf("failed to configure db: %v", err)
	}

	conn := &Connection{db: db}
	conn.Migrate()
	return conn
}

// configureSQLite applies performance and safety pragmas.
func configureSQLite(db *sql.DB) error {
	pragmas := `
	PRAGMA busy_timeout = 5000;
	PRAGMA journal_mode = WAL;
	PRAGMA journal_size_limit = 200000000;
	PRAGMA synchronous = NORMAL;
	PRAGMA foreign_keys = ON;
	PRAGMA temp_store = MEMORY;
	PRAGMA mmap_size = 300000000;
	PRAGMA page_size = 32768;
	PRAGMA cache_size = -16000;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, pragmas); err != nil {
		return fmt.Errorf("executing pragmas: %w", err)
	}

	db.SetMaxOpenConns(10)                 // Allow multiple concurrent connections
	db.SetMaxIdleConns(5)                  // Keep some connections alive
	db.SetConnMaxLifetime(0)               // Reuse connections indefinitely
	db.SetConnMaxIdleTime(5 * time.Minute) // Close idle connections after 5min

	return nil
}

func (c *Connection) Get() *sql.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.db.Close()
}

func (c *Connection) Ping() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db.Ping()
}

func (c *Connection) GetQuery() *db.Queries {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return db.New(c.db)
}

// Replace replaces the on‐disk DB with srcPath, then reopens it.
func (c *Connection) Replace(srcPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Checkpoint & close old DB
	if c.db != nil {
		if _, err := c.db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"); err != nil {
			slog.Warn("checkpoint old database failed", "error", err)
		}
		if err := c.db.Close(); err != nil {
			slog.Warn("close old database failed", "error", err)
		}
	}

	// Allow locks to clear
	time.Sleep(500 * time.Millisecond)

	// Copy new file into place
	dst := util.ResolvePath("mantrae.db")
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove old db file: %w", err)
	}
	if err := util.CopyFile(srcPath, dst); err != nil {
		return fmt.Errorf("copy new db file: %w", err)
	}

	// Reinitialize the database connection
	conn := NewConnection("")
	c.db = conn.db
	return nil
}

func (c *Connection) Migrate() {
	goose.SetBaseFS(migrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(c.db, "migrations"); err != nil {
		log.Fatal(err)
	}
}
