// Package backup provides functionality for creating and restoring backups.
package backup

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/storage"
	"github.com/mizuchilabs/mantrae/server/internal/store"
	"github.com/mizuchilabs/mantrae/server/internal/traefik"
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
	interval, ok := m.SM.Get(ctx, settings.KeyBackupInterval)
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

	backupNameDB := fmt.Sprintf("backup_%s.db", time.Now().UTC().Format("20060102_150405"))

	// Create a temporary file for the db backup
	tmpFile, err := os.CreateTemp("", "sqlite_backup_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			slog.Error("failed to remove temp file", "error", err)
		}
		if err = tmpFile.Close(); err != nil {
			slog.Error("failed to close temp file", "error", err)
		}
	}()

	// Perform SQLite backup
	if _, err = m.Conn.Get().Exec("VACUUM INTO ?", tmpFile.Name()); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Rewind the file for reading
	if _, err = tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to rewind temp file: %w", err)
	}

	// Store the backup using the backend
	if err = m.Storage.Store(ctx, backupNameDB, tmpFile); err != nil {
		return fmt.Errorf("failed to store backup: %w", err)
	}

	// Perform YAML backup
	if err = traefik.BackupDynamicConfigs(ctx, m.Conn.GetQuery(), m.Storage); err != nil {
		return fmt.Errorf("failed to backup dynamic configs: %w", err)
	}

	// Clean up older backups
	if err := m.cleanup(ctx); err != nil {
		return fmt.Errorf("failed to cleanup backups: %w", err)
	}

	slog.Info("Backup created successfully", "name", backupNameDB)
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

	// Get the backup from storage
	reader, err := m.Storage.Retrieve(ctx, backupName)
	if err != nil {
		return fmt.Errorf("failed to retrieve backup: %w", err)
	}
	defer func() {
		if err = reader.Close(); err != nil {
			slog.Error("failed to close backup reader", "error", err)
		}
	}()

	// Create a temporary file for the backup
	tmpFile, err := os.CreateTemp("", "restore_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			slog.Error("failed to remove temp file", "error", err)
		}
		if err = tmpFile.Close(); err != nil {
			slog.Error("failed to close temp file", "error", err)
		}
	}()

	// Copy backup to temp file
	if _, err = io.Copy(tmpFile, reader); err != nil {
		return fmt.Errorf("failed to copy backup to temp file: %w", err)
	}

	// Reinitialize the database
	if err = m.Conn.Replace(tmpFile.Name()); err != nil {
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
	// Check if filename matches pattern *.db or *.yaml
	matchedDB, err := filepath.Match("*.db", filename)
	if err != nil {
		return false
	}
	if matchedDB {
		return true
	}

	matchedYaml, err := filepath.Match("*.yaml", filename)
	if err != nil {
		return false
	}
	if matchedYaml {
		return true
	}
	return false
}

func (m *BackupManager) cleanup(ctx context.Context) error {
	// Regex to parse filenames like backup_20230710_101010.db
	dbBackupRegex := regexp.MustCompile(`backup_(\d{8}_\d{6})\.db`)

	// Regex to parse filenames like backup_profileName_20230710_101010.yaml
	yamlBackupRegex := regexp.MustCompile(`backup_(.+)_(\d{8}_\d{6})\.yaml`)

	files, err := m.Storage.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	retention, ok := m.SM.Get(ctx, settings.KeyBackupKeep)
	if !ok {
		return fmt.Errorf("failed to get retention setting")
	}

	// Sort descending by time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(*files[j].Timestamp)
	})

	var dbBackups []storage.StoredFile
	yamlBackupsByProfile := make(map[string][]storage.StoredFile)

	for _, f := range files {
		if m := dbBackupRegex.FindStringSubmatch(f.Name); m != nil {
			dbBackups = append(dbBackups, f)
		} else if m := yamlBackupRegex.FindStringSubmatch(f.Name); m != nil {
			profile := m[1]
			yamlBackupsByProfile[profile] = append(yamlBackupsByProfile[profile], f)
		} else {
			continue
		}
	}

	// Delete older backups
	if len(dbBackups) > settings.AsInt(retention) {
		for _, f := range dbBackups[settings.AsInt(retention):] {
			if err := m.Storage.Delete(ctx, f.Name); err != nil {
				return fmt.Errorf("failed to delete db backup %s: %w", f.Name, err)
			}
		}
	}
	// Cleanup YAML backups per profile
	for profile, backups := range yamlBackupsByProfile {
		if len(backups) <= settings.AsInt(retention) {
			continue
		}
		for _, f := range backups[settings.AsInt(retention):] {
			if err := m.Storage.Delete(ctx, f.Name); err != nil {
				return fmt.Errorf(
					"failed to delete yaml backup %s (profile %s): %w",
					f.Name,
					profile,
					err,
				)
			}
		}
	}

	return nil
}
