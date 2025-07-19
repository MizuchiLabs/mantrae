package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/app"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

type Connection struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewDBConnection() (*Connection, error) {
	database, err := InitDB()
	if err != nil {
		return nil, err
	}
	return &Connection{db: database}, nil
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

func (c *Connection) GetQuery() *Queries {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return New(c.db)
}

func (c *Connection) Replace(srcPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close old connection if it exists
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
	dst := app.ResolvePath("mantrae.db")
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove old db file: %w", err)
	}
	if err := util.CopyFile(srcPath, dst); err != nil {
		return fmt.Errorf("copy new db file: %w", err)
	}

	// Create new connection before closing the old one
	newDB, err := InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize new database: %w", err)
	}

	// Update to new connection
	c.db = newDB

	return nil
}
