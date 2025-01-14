package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/robfig/cron/v3"
)

type BackupManager struct {
	Config  BackupConfig
	db      *db.Queries
	DBType  string
	DBPath  string
	cronJob *cron.Cron
}

type BackupMetadata struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	Created  time.Time `json:"created"`
	DBType   string    `json:"db_type"`
	Version  string    `json:"version"`
}

func NewBackupManager(config Config, db *db.Queries) (*BackupManager, error) {
	if err := os.MkdirAll(config.Backup.Dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// enabled, err := db.GetSetting(context.Background(), "backup-enabled")
	// if err != nil {
	// 	return nil, err
	// }
	// schedule, err := db.GetSetting(context.Background(), "backup-schedule")
	// if err != nil {
	// 	return nil, err
	// }
	// keep, err := db.GetSetting(context.Background(), "backup-keep")
	// if err != nil {
	// 	return nil, err
	// }

	return &BackupManager{
		Config:  config.Backup,
		db:      db,
		DBType:  config.Database.Type,
		DBPath:  config.DBPath(),
		cronJob: cron.New(),
	}, nil
}

func (bm *BackupManager) Start(ctx context.Context) error {
	if !bm.Config.Enabled {
		return nil
	}

	// Schedule backup job
	_, err := bm.cronJob.AddFunc(bm.Config.Schedule, func() {
		if err := bm.CreateBackup(ctx); err != nil {
			slog.Error("Scheduled backup failed", "error", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule backup: %w", err)
	}

	bm.cronJob.Start()
	return nil
}

func (bm *BackupManager) Stop() {
	if bm.cronJob != nil {
		bm.cronJob.Stop()
	}
}

func (bm *BackupManager) CreateBackup(ctx context.Context) error {
	timestamp := time.Now().UTC()
	backupName := fmt.Sprintf("backup_%s.db", timestamp.Format("20060102_150405"))
	backupPath := filepath.Join(bm.Config.Dir, backupName)

	// Create backup file
	if err := bm.performBackup(ctx, backupPath); err != nil {
		return err
	}

	// Create metadata file
	metadata := BackupMetadata{
		Version: "1.0",
		Created: timestamp,
		DBType:  bm.DBType,
	}

	metadataPath := backupPath + ".json"
	if err := bm.saveMetadata(metadataPath, metadata); err != nil {
		return err
	}

	// Cleanup old backups
	if err := bm.cleanup(); err != nil {
		slog.Error("Backup cleanup failed", "error", err)
	}

	return nil
}

func (bm *BackupManager) performBackup(ctx context.Context, destPath string) error {
	switch bm.DBType {
	case "sqlite":
		return bm.backupSQLite(ctx, destPath)
	case "postgres":
		return bm.backupPostgres(ctx, destPath)
	default:
		return fmt.Errorf("unsupported database type: %s", bm.DBType)
	}
}

func (bm *BackupManager) backupSQLite(ctx context.Context, destPath string) error {
	// For SQLite, we can simply copy the database file
	src, err := os.Open(bm.DBPath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy database: %w", err)
	}

	return nil
}

func (bm *BackupManager) backupPostgres(ctx context.Context, destPath string) error {
	// Implement pg_dump here
	// Example implementation:
	/*
		cmd := exec.CommandContext(ctx, "pg_dump",
			"-h", host,
			"-U", username,
			"-d", dbname,
			"-f", destPath,
		)
		return cmd.Run()
	*/
	return fmt.Errorf("postgres backup not implemented")
}

func (bm *BackupManager) saveMetadata(path string, metadata BackupMetadata) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func (bm *BackupManager) cleanup() error {
	files, err := filepath.Glob(filepath.Join(bm.Config.Dir, "backup_*.db"))
	if err != nil {
		return fmt.Errorf("failed to list backup files: %w", err)
	}

	if len(files) <= bm.Config.Keep {
		return nil
	}

	// Sort files by modification time
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	fileInfos := make([]fileInfo, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{file, info.ModTime()})
	}

	// Sort by modification time (oldest first)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.Before(fileInfos[j].modTime)
	})

	// Remove oldest files
	for i := 0; i < len(fileInfos)-bm.Config.Keep; i++ {
		if err := os.Remove(fileInfos[i].path); err != nil {
			slog.Error("Failed to remove old backup", "file", fileInfos[i].path, "error", err)
			continue
		}
		// Also remove metadata file if it exists
		metadataPath := fileInfos[i].path + ".json"
		if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
			slog.Error("Failed to remove metadata file", "file", metadataPath, "error", err)
		}
	}

	return nil
}

// RestoreBackup restores a database from a backup file
func (bm *BackupManager) RestoreBackup(ctx context.Context, backupPath string) error {
	// Verify backup metadata
	metadataPath := backupPath + ".json"
	metadata, err := bm.loadMetadata(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to load backup metadata: %w", err)
	}

	if metadata.DBType != bm.DBType {
		return fmt.Errorf("backup database type (%s) doesn't match current type (%s)",
			metadata.DBType, bm.DBType)
	}

	switch bm.DBType {
	case "sqlite":
		return bm.restoreSQLite(ctx, backupPath)
	case "postgres":
		return bm.restorePostgres(ctx, backupPath)
	default:
		return fmt.Errorf("unsupported database type: %s", bm.DBType)
	}
}

func (bm *BackupManager) loadMetadata(path string) (BackupMetadata, error) {
	var metadata BackupMetadata
	file, err := os.Open(path)
	if err != nil {
		return metadata, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return metadata, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return metadata, nil
}

func (bm *BackupManager) restoreSQLite(ctx context.Context, backupPath string) error {
	// Create a temporary database path
	tempDB := bm.DBPath + ".tmp"

	// Copy backup to temporary location
	if err := copyFile(backupPath, tempDB); err != nil {
		return fmt.Errorf("failed to create temporary database: %w", err)
	}

	// Test the temporary database
	if err := bm.testDatabase(tempDB); err != nil {
		os.Remove(tempDB)
		return fmt.Errorf("backup verification failed: %w", err)
	}

	// Replace the current database
	if err := os.Rename(tempDB, bm.DBPath); err != nil {
		os.Remove(tempDB)
		return fmt.Errorf("failed to replace database: %w", err)
	}

	return nil
}

func (bm *BackupManager) restorePostgres(ctx context.Context, backupPath string) error {
	// Implement pg_restore here
	return fmt.Errorf("postgres restore not implemented")
}

func (bm *BackupManager) testDatabase(dbPath string) error {
	// Implement database verification here
	// For example, try to open the database and run a simple query
	return nil
}

// ListBackups returns a list of available backups
func (bm *BackupManager) ListBackups() ([]BackupMetadata, error) {
	files, err := filepath.Glob(filepath.Join(bm.Config.Dir, "backup_*.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to list backup files: %w", err)
	}

	backups := make([]BackupMetadata, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		metadata, err := bm.loadMetadata(file + ".json")
		if err != nil {
			// If metadata is missing, use file info
			backups = append(backups, BackupMetadata{
				Filename: filepath.Base(file),
				Size:     info.Size(),
				Created:  info.ModTime(),
				DBType:   bm.DBType,
			})
			continue
		}

		backups = append(backups, BackupMetadata{
			Filename: filepath.Base(file),
			Size:     info.Size(),
			Created:  metadata.Created,
			DBType:   metadata.DBType,
			Version:  metadata.Version,
		})
	}

	// Sort by creation time, newest first
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Created.After(backups[j].Created)
	})

	return backups, nil
}

// GetLatestBackup returns the filename of the most recent backup
func (bm *BackupManager) GetLatestBackup() (string, error) {
	backups, err := bm.ListBackups()
	if err != nil {
		return "", fmt.Errorf("failed to list backups: %w", err)
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found")
	}

	// Since ListBackups returns backups sorted by creation time (newest first),
	// we can just return the first one
	return backups[0].Filename, nil
}

// IsValidBackupFile checks if the given filename is a valid backup file
func (bm *BackupManager) IsValidBackupFile(filename string) bool {
	// Prevent directory traversal
	if strings.Contains(filename, "..") {
		return false
	}

	// Check if file matches expected pattern
	matched, err := filepath.Match("backup_*.db", filename)
	if err != nil || !matched {
		return false
	}

	// Check if file exists in backup directory
	fullPath := filepath.Join(bm.Config.Dir, filename)
	if _, err := os.Stat(fullPath); err != nil {
		return false
	}

	return true
}

// GetBackupPath returns the full path for a given backup filename
func (bm *BackupManager) GetBackupPath(filename string) (string, error) {
	if !bm.IsValidBackupFile(filename) {
		return "", fmt.Errorf("invalid backup filename: %s", filename)
	}
	return filepath.Join(bm.Config.Dir, filename), nil
}

// DeleteBackup removes a specific backup and its metadata
func (bm *BackupManager) DeleteBackup(filename string) error {
	if !bm.IsValidBackupFile(filename) {
		return fmt.Errorf("invalid backup filename: %s", filename)
	}

	backupPath := filepath.Join(bm.Config.Dir, filename)
	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	// Try to remove metadata file if it exists
	metadataPath := backupPath + ".json"
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata file: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
