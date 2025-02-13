package backup

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/app"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/storage"
)

type BackupManager struct {
	Conn      *db.Connection
	Config    *app.BackupConfig
	Backend   storage.Backend
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	mu        sync.Mutex
}

func NewManager(
	conn *db.Connection,
	config app.BackupConfig,
	backend storage.Backend,
) *BackupManager {
	return &BackupManager{
		Conn:     conn,
		Config:   &config,
		Backend:  backend,
		stopChan: make(chan struct{}),
	}
}

func (m *BackupManager) Start(ctx context.Context) {
	m.waitGroup.Add(1)
	go m.backupLoop(ctx)
}

func (m *BackupManager) Stop() {
	close(m.stopChan)
	m.waitGroup.Wait()
}

func (m *BackupManager) backupLoop(ctx context.Context) {
	defer m.waitGroup.Done()
	ticker := time.NewTicker(m.Config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.Stop()
			return
		case <-ticker.C:
			if err := m.Create(ctx); err != nil {
				slog.Error("backup failed", "error", err)
			}
			if err := m.cleanup(ctx); err != nil {
				slog.Error("backup cleanup failed", "error", err)
			}
		}
	}
}

func (m *BackupManager) Create(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	backupName := fmt.Sprintf("backup_%s.db", time.Now().UTC().Format("20060102_150405"))

	// Create a temporary file for the backup
	tmpFile, err := os.CreateTemp("", "sqlite_backup_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db := m.Conn.Get()
	// Perform SQLite backup
	if _, err := db.Exec("VACUUM INTO ?", tmpFile.Name()); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Rewind the file for reading
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to rewind temp file: %w", err)
	}

	// Store the backup using the backend
	if err := m.Backend.Store(ctx, backupName, tmpFile); err != nil {
		return fmt.Errorf("failed to store backup: %w", err)
	}

	// Clean up older backups
	if err := m.cleanup(ctx); err != nil {
		return fmt.Errorf("failed to cleanup backups: %w", err)
	}

	return nil
}

func (m *BackupManager) Restore(ctx context.Context, backupName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dbPath := m.Config.DatabaseName + ".db"
	walPath := dbPath + "-wal"
	shmPath := dbPath + "-shm"

	// Get the backup from storage
	reader, err := m.Backend.Retrieve(ctx, backupName)
	if err != nil {
		return fmt.Errorf("failed to retrieve backup: %w", err)
	}
	defer reader.Close()

	// Create a temporary file for the backup
	tmpFile, err := os.CreateTemp("", "restore_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Copy backup to temp file
	if _, err = io.Copy(tmpFile, reader); err != nil {
		return fmt.Errorf("failed to copy backup to temp file: %w", err)
	}

	// Close the temp file to ensure all data is written
	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	// Close current database connections
	if err = m.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	// Remove WAL and SHM files if they exist
	os.Remove(walPath)
	os.Remove(shmPath)

	// Copy the temp file to the database location instead of rename (invalid cross-device link)
	srcFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open temp file for copying: %w", err)
	}
	defer srcFile.Close()

	// Create new database file
	dstFile, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create new database file: %w", err)
	}
	defer dstFile.Close()

	// Copy the contents
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy database contents: %w", err)
	}

	// Ensure all data is written to disk
	if err = dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync database file: %w", err)
	}

	// Reinitialize the database
	if err = m.Conn.Replace(); err != nil {
		return err
	}

	return nil
}

func (m *BackupManager) IsValidBackupFile(filename string) bool {
	// Prevent directory traversal
	if strings.Contains(filename, "..") {
		return false
	}
	// Check if filename matches pattern: backup_YYYYMMDD_HHMMSS.db
	matched, err := filepath.Match("backup_[0-9]*_[0-9]*.db", filename)
	if err != nil {
		return false
	}
	return matched
}

func (m *BackupManager) cleanup(ctx context.Context) error {
	files, err := m.Backend.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if len(files) <= m.Config.Keep {
		return nil
	}

	// Delete older backups
	for _, file := range files[m.Config.Keep:] {
		if err := m.Backend.Delete(ctx, file.Name); err != nil {
			return fmt.Errorf("failed to delete old backup %s: %w", file.Name, err)
		}
	}

	return nil
}
