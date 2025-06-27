package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/util"
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
	var dataSource string
	if path == "" {
		dataSource = fmt.Sprintf("file:%s", filepath.ToSlash(util.ResolvePath("mantrae.db")))
	} else {
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

func (c *Connection) Replace() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close old connection if it exists
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			// Log error but don't fail - we already have new connection
			fmt.Printf("warning: failed to close old database connection: %v\n", err)
		}
	}

	// Wait a small amount of time for SQLite to release locks
	time.Sleep(100 * time.Millisecond)

	// Create new connection before closing the old one
	conn := NewConnection("")

	// Update to new connection
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
