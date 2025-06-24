package backup

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/storage"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/util"
)

const BackupPath = "backups"

type BackupManager struct {
	Conn      *store.Connection
	SM        *settings.SettingsManager
	Storage   storage.Backend
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	mu        sync.Mutex
}

func NewManager(conn *store.Connection, sm *settings.SettingsManager) *BackupManager {
	if conn == nil {
		log.Fatal("No database connection provided")
	}
	if sm == nil {
		log.Fatal("No settings manager provided")
	}

	return &BackupManager{
		Conn:     conn,
		SM:       sm,
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
	var err error
	m.Storage, err = storage.GetBackend(ctx, m.SM, BackupPath)
	if err != nil {
		return fmt.Errorf("failed to get storage backend: %w", err)
	}
	return nil
}

func (m *BackupManager) backupLoop(ctx context.Context) {
	defer m.waitGroup.Done()

	// Get backup interval
	interval, ok := m.SM.Get(settings.KeyBackupInterval)
	if !ok {
		slog.Error("failed to get backup interval")
		return
	}

	ticker := time.NewTicker(settings.AsDuration(interval))
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

	dbPath := util.ResolvePath("mantrae.db")
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
	if err := m.SetStorage(ctx); err != nil {
		return fmt.Errorf("failed to set storage: %w", err)
	}
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
	// Check if filename matches pattern *.db
	matched, err := filepath.Match("*.db", filename)
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

	retention, ok := m.SM.Get(settings.KeyBackupKeep)
	if !ok {
		return fmt.Errorf("failed to get retention setting")
	}
	if len(files) <= settings.AsInt(retention) {
		return nil
	}

	// Delete older backups
	for _, file := range files[settings.AsInt(retention):] {
		if err := m.Storage.Delete(ctx, file.Name); err != nil {
			return fmt.Errorf("failed to delete old backup %s: %w", file.Name, err)
		}
	}

	return nil
}
