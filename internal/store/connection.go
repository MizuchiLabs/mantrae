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

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/sqlite-schema-diff/pkg/diff"
	"github.com/mizuchilabs/sqlite-schema-diff/pkg/parser"
	_ "modernc.org/sqlite"
)

//go:embed schemas/*.sql
var schemaFS embed.FS

type Connection struct {
	mu  sync.RWMutex
	ctx context.Context
	db  *sql.DB
	Q   *db.Queries
}

// NewConnection opens a SQLite connection.
// If `path` is empty, it defaults to "mantrae.db" in the data dir.
// If `path` is ":memory:" or "file::memory:?cache=shared", opens in-memory.
func NewConnection(ctx context.Context, path string) *Connection {
	dataSource := fmt.Sprintf(
		"file:%s?_txlock=immediate",
		filepath.ToSlash(util.ResolvePath("mantrae.db")),
	)
	if path != "" {
		dataSource = path
	}

	sqliteDB, err := sql.Open("sqlite", dataSource)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := setupSQLite(sqliteDB); err != nil {
		log.Fatalf("failed to configure db: %v", err)
	}
	migrate(sqliteDB)

	conn := &Connection{
		ctx: ctx,
		db:  sqliteDB,
		Q:   db.New(sqliteDB),
	}

	// Wait for shutdown signal
	go func() {
		<-ctx.Done()
		if err := sqliteDB.Close(); err != nil {
			slog.Error("Failed to close database", "error", err)
		}
	}()

	return conn
}

// setupSQLite applies performance and safety pragmas.
func setupSQLite(db *sql.DB) error {
	pragmas := `
	PRAGMA busy_timeout = 5000;
	PRAGMA journal_mode = WAL;
	PRAGMA journal_size_limit = 200000000;
	PRAGMA synchronous = NORMAL;
	PRAGMA foreign_keys = ON;
	PRAGMA temp_store = MEMORY;
	PRAGMA mmap_size = 300000000;
	PRAGMA cache_size = -16000;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, pragmas); err != nil {
		return fmt.Errorf("executing pragmas: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	return nil
}

func (c *Connection) Get() *sql.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
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

	// Also remove WAL and SHM files
	_ = os.Remove(dst + "-wal")
	_ = os.Remove(dst + "-shm")

	if err := util.CopyFile(srcPath, dst); err != nil {
		return fmt.Errorf("copy new db file: %w", err)
	}

	// Reopen the database with the same data source
	dataSource := fmt.Sprintf(
		"file:%s?_txlock=immediate",
		filepath.ToSlash(dst),
	)

	sqliteDB, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return fmt.Errorf("failed to reopen db: %w", err)
	}

	if err := setupSQLite(sqliteDB); err != nil {
		return fmt.Errorf("failed to configure db: %w", err)
	}

	// Run migrations on the restored database
	migrate(sqliteDB)

	// Update connection fields
	c.db = sqliteDB
	c.Q = db.New(sqliteDB)
	return nil
}

func migrate(db *sql.DB) {
	parser.SetBaseFS(schemaFS)
	if err := diff.Apply(db, "schemas", diff.ApplyOptions{}); err != nil {
		slog.Error("failed to apply schema changes", "error", err)
		return
	}
}
