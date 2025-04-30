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
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/storage"
)

type BackupManager struct {
	Conn      *db.Connection
	Settings  *settings.SettingsManager
	Storage   storage.Backend
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	mu        sync.Mutex
}

func NewManager(
	conn *db.Connection,
	settings *settings.SettingsManager,
) *BackupManager {
	return &BackupManager{
		Conn:     conn,
		Settings: settings,
		stopChan: make(chan struct{}),
	}
}

func (m *BackupManager) Start(ctx context.Context) {
	// Init storage
	if err := m.SetStorage(ctx); err != nil {
		slog.Error("backup failed", "error", err)
	}

	m.waitGroup.Add(1)
	go m.backupLoop(ctx)
}

func (m *BackupManager) Stop() {
	close(m.stopChan)
	m.waitGroup.Wait()
}

func (m *BackupManager) SetStorage(ctx context.Context) error {
	storageSet, _ := m.Settings.Get(ctx, settings.KeyBackupStorage)
	storageType := storage.BackendType(storageSet.String("local"))
	if !storageType.Valid() {
		return fmt.Errorf("storage backend not configured")
	}

	var err error
	var newStorage storage.Backend
	switch storageType {
	case storage.BackendTypeLocal:
		pathSet, err := m.Settings.Get(ctx, settings.KeyBackupPath)
		if err != nil {
			return fmt.Errorf("failed to get backup path: %w", err)
		}

		path := pathSet.String("backups")
		newStorage, err = storage.NewLocalStorage(path)
		if err != nil {
			return fmt.Errorf("failed to create local storage: %w", err)
		}
		slog.Debug("backup storage set to local", "path", path)

	case storage.BackendTypeS3:
		newStorage, err = storage.NewS3Storage(ctx, m.Settings)
		if err != nil {
			return fmt.Errorf("failed to create s3 storage: %w", err)
		}
		slog.Debug("backup storage set to S3")

	default:
		return fmt.Errorf("unsupported backend type: %s", storageType)
	}

	m.Storage = newStorage
	return nil
}

func (m *BackupManager) backupLoop(ctx context.Context) {
	defer m.waitGroup.Done()

	// Get backup interval
	intervalSet, err := m.Settings.Get(ctx, settings.KeyBackupInterval)
	if err != nil {
		slog.Error("failed to get backup interval", "error", err)
		return
	}
	interval := intervalSet.Duration(24)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.Stop()
			return
		case <-ticker.C:
			// Set storage before creating backup (live change)
			if err := m.SetStorage(ctx); err != nil {
				slog.Error("backup failed", "error", err)
				continue
			}
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

	// Set storage
	if err := m.SetStorage(ctx); err != nil {
		return fmt.Errorf("failed to set storage: %w", err)
	}

	backupName := fmt.Sprintf("backup_%s.db", time.Now().UTC().Format("20060102_150405"))

	// Create a temporary file for the backup
	tmpFile, err := os.CreateTemp("", "sqlite_backup_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

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
	if err := m.Storage.Store(ctx, backupName, tmpFile); err != nil {
		return fmt.Errorf("failed to store backup: %w", err)
	}

	// Clean up older backups
	if err := m.cleanup(ctx); err != nil {
		return fmt.Errorf("failed to cleanup backups: %w", err)
	}

	slog.Info("Backup created successfully", "name", backupName)
	return nil
}

func (m *BackupManager) Restore(ctx context.Context, backupName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Set storage
	if err := m.SetStorage(ctx); err != nil {
		return fmt.Errorf("failed to set storage: %w", err)
	}

	// Validate backup name for security
	if !m.IsValidBackupFile(backupName) {
		return fmt.Errorf("invalid backup file name")
	}

	dbPath := app.ResolvePath("mantrae.db")
	walPath := dbPath + "-wal"
	shmPath := dbPath + "-shm"

	// Get the backup from storage
	reader, err := m.Storage.Retrieve(ctx, backupName)
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

func (m *BackupManager) List(ctx context.Context) ([]storage.StoredFile, error) {
	// Set storage
	if err := m.SetStorage(ctx); err != nil {
		return nil, fmt.Errorf("failed to set storage: %w", err)
	}

	files, err := m.Storage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	// Filter out any non-backup files
	var backups []storage.StoredFile
	for _, file := range files {
		if m.IsValidBackupFile(file.Name) {
			backups = append(backups, file)
		}
	}

	return backups, nil
}

func (m *BackupManager) Delete(ctx context.Context, id string) error {
	// Set storage
	if err := m.SetStorage(ctx); err != nil {
		return fmt.Errorf("failed to set storage: %w", err)
	}

	// Delete backup file
	if err := m.Storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete backup %s: %w", id, err)
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
	files, err := m.Storage.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	retentionSet, err := m.Settings.Get(ctx, settings.KeyBackupKeep)
	if err != nil {
		return fmt.Errorf("failed to get retention setting: %w", err)
	}
	retention := retentionSet.Int(7)
	if len(files) <= retention {
		return nil
	}

	// Delete older backups
	for _, file := range files[retention:] {
		if err := m.Storage.Delete(ctx, file.Name); err != nil {
			return fmt.Errorf("failed to delete old backup %s: %w", file.Name, err)
		}
	}

	return nil
}
