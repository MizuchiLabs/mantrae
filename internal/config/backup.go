package config

import (
	"context"
	"database/sql"
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
	DBType  string
	DBPath  string
	cronJob *cron.Cron
	db      *sql.DB
	q       *db.Queries
}

type BackupMetadata struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	Created  time.Time `json:"created"`
	DBType   string    `json:"db_type"`
	Version  string    `json:"version"`
}

func NewBackupManager(ctx context.Context, config Config) (*BackupManager, error) {
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

	bm := &BackupManager{
		Config:  config.Backup,
		DBType:  config.Database.Type,
		DBPath:  config.DBPath(),
		cronJob: cron.New(),
	}
	if err := bm.Start(ctx); err != nil {
		return nil, err
	}

	return bm, nil
}

func (bm *BackupManager) Start(ctx context.Context) error {
	// Ensure previous connections are closed
	if bm.db != nil {
		if err := bm.Stop(); err != nil {
			return fmt.Errorf("failed to stop existing connections: %w", err)
		}
	}

	// Small delay before opening new connections
	time.Sleep(100 * time.Millisecond)

	// Initialize new database connection with retries
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		newDB, err := db.InitDB()
		if err == nil {
			bm.q = db.New(newDB)
			bm.db = newDB
			break
		}
		if i == maxRetries-1 {
			return fmt.Errorf(
				"failed to initialize database after %d attempts: %w",
				maxRetries,
				err,
			)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if bm.Config.Enabled {
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
	}

	return nil
}

func (bm *BackupManager) Stop() error {
	if bm.cronJob != nil {
		bm.cronJob.Stop()
	}

	if bm.db != nil {
		if err := bm.db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}
	// Small delay to ensure SQLite cleanup
	time.Sleep(200 * time.Millisecond)

	return nil
}

func (bm *BackupManager) CreateBackup(ctx context.Context) error {
	timestamp := time.Now().UTC()
	backupName := fmt.Sprintf("backup_%s.db", timestamp.Format("20060102_150405"))
	backupPath := filepath.Join(bm.Config.Dir, backupName)

	if err := bm.backupSQLite(backupPath); err != nil {
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

func (bm *BackupManager) backupSQLite(destPath string) error {
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

func (bm *BackupManager) RestoreBackup(ctx context.Context, backupPath string) error {
	ext := filepath.Ext(backupPath)

	if ext == ".db" {
		// For .db files, use the existing file replacement method
		tempDB := bm.DBPath + ".tmp"
		defer os.Remove(tempDB)

		// Copy backup to temporary location
		if err := copyFile(backupPath, tempDB); err != nil {
			return fmt.Errorf("failed to create temporary database: %w", err)
		}

		// Test the temporary database
		if err := bm.testDatabase(tempDB); err != nil {
			return fmt.Errorf("backup verification failed: %w", err)
		}

		// Create a backup of the current database
		currentBackup := bm.DBPath + ".old"
		if err := copyFile(bm.DBPath, currentBackup); err != nil {
			return fmt.Errorf("failed to backup current database: %w", err)
		}

		// Close the current database connection
		if err := bm.Stop(); err != nil {
			return fmt.Errorf("failed to close database connections: %w", err)
		}

		// Replace the current database
		if err := os.Rename(tempDB, bm.DBPath); err != nil {
			// If replacement fails, try to restore the old database
			if err = os.Rename(currentBackup, bm.DBPath); err != nil {
				return fmt.Errorf("failed to restore old database: %w", err)
			}
			slog.Error("Failed to replace database", "error", err)
		}

		// Clean up the old backup
		os.Remove(currentBackup)

		// Start the new database
		if err := bm.Start(ctx); err != nil {
			return fmt.Errorf("failed to start database connections: %w", err)
		}
		return nil
	}

	if ext == ".sql" {
		// For .sql files, execute the SQL statements
		sqlContent, err := os.ReadFile(backupPath)
		if err != nil {
			return fmt.Errorf("failed to read SQL backup file: %w", err)
		}

		// Create a new database connection for the restore
		db, err := sql.Open("sqlite", bm.DBPath)
		if err != nil {
			return fmt.Errorf("failed to open database for restore: %w", err)
		}
		defer db.Close()

		// Start a transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()

		// Execute SQL statements
		if _, err := tx.ExecContext(ctx, string(sqlContent)); err != nil {
			return fmt.Errorf("failed to execute SQL backup: %w", err)
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit restore transaction: %w", err)
		}

		return nil
	}

	return fmt.Errorf("unsupported backup file format: %s", ext)
}

func (bm *BackupManager) testDatabase(dbPath string) error {
	// Try to open the database
	testDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open backup database: %w", err)
	}
	defer testDB.Close()

	// Test if database is accessible with a simple query
	if err = testDB.Ping(); err != nil {
		return fmt.Errorf("backup database is not accessible: %w", err)
	}

	// Start a read-only transaction
	tx, err := testDB.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get list of tables
	rows, err := tx.Query(`
        SELECT name FROM sqlite_master 
        WHERE type='table' 
        AND name NOT LIKE 'sqlite_%' 
        AND name NOT LIKE 'goose_%'
    `)
	if err != nil {
		return fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	// Check each table's schema
	var tables []string
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if len(tables) == 0 {
		return fmt.Errorf("backup database contains no tables")
	}

	// Verify schema for each table
	for _, table := range tables {
		// Get table info
		rows, err := tx.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			return fmt.Errorf("failed to get schema for table %s: %w", table, err)
		}

		var hasColumns bool
		for rows.Next() {
			hasColumns = true
			break
		}
		rows.Close()

		if !hasColumns {
			return fmt.Errorf("table %s has no columns", table)
		}
	}

	// Test foreign key constraints
	_, err = tx.Exec("PRAGMA foreign_key_check")
	if err != nil {
		return fmt.Errorf("foreign key check failed: %w", err)
	}

	// Try to read a row from each table to ensure data is accessible
	for _, table := range tables {
		query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table)
		_, err := tx.Query(query)
		if err != nil {
			return fmt.Errorf("failed to read from table %s: %w", table, err)
		}
	}

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
	matched, err := filepath.Match("backup_*.*", filename)
	if err != nil || !matched {
		return false
	}

	// Check file extension
	ext := filepath.Ext(filename)
	if ext != ".db" && ext != ".sql" {
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
