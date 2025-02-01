package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
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
	newDB, err := InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize new database: %w", err)
	}

	// Update to new connection
	c.db = newDB

	return nil
}
