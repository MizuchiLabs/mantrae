package config

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/robfig/cron/v3"
)

var backupCron *cron.Cron

func BackupDatabase() error {
	timestamp := time.Now().Format("2006-01-02")
	backupPath := fmt.Sprintf("%s/backup-%s.sql.gz", util.Path(util.BackupDir), timestamp)

	// Create the backup directory if it doesn't exist
	backupDir := filepath.Dir(backupPath)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(backupDir, 0750); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	// Check if the backup file already exists
	if _, err := os.Stat(backupPath); err == nil {
		return nil
	}

	sqlFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create SQL file: %w", err)
	}
	defer sqlFile.Close()

	cmd := exec.Command("sqlite3", util.Path(util.DBName), ".dump")
	sqlPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	// Create the gzip file
	gzipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create gzip file: %w", err)
	}
	defer gzipFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	// Start the sqlite3 command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sqlite3 command: %w", err)
	}

	// Pipe the SQL dump into the gzip writer
	if _, err := io.Copy(gzipWriter, sqlPipe); err != nil {
		return fmt.Errorf("failed to compress SQL dump: %w", err)
	}

	// Wait for the sqlite3 command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("sqlite3 command failed: %w", err)
	}

	slog.Info("Database backup created", "file", backupPath)
	return nil
}

// CleanupBackups deletes the oldest backup files to keep only the specified number of backups
func CleanupBackups() error {
	keepSetting, err := db.Query.GetSettingByKey(context.Background(), "backup-keep")
	if err != nil {
		return fmt.Errorf("failed to get backup-keep setting: %w", err)
	}
	if keepSetting.Value == "" {
		return nil
	}

	keep, err := strconv.Atoi(keepSetting.Value)
	if err != nil {
		return fmt.Errorf("failed to parse backup-keep setting: %w", err)
	}
	// If keep is 0, do not delete any backups
	if keep == 0 {
		return nil
	}

	// Get the list of backup files

	files, err := filepath.Glob(fmt.Sprintf("%s/backup-*.sql.gz", util.Path(util.BackupDir)))
	if err != nil {
		return fmt.Errorf("failed to list backup files: %w", err)
	}

	if len(files) <= keep {
		return nil
	}

	// Sort the backup files by modification time (oldest to newest)
	sort.Slice(files, func(i, j int) bool {
		iInfo, _ := os.Stat(files[i])
		jInfo, _ := os.Stat(files[j])
		return iInfo.ModTime().Before(jInfo.ModTime())
	})

	// Delete the oldest backup files (the first N - keep files)
	for i := 0; i < len(files)-keep; i++ {
		if err := os.Remove(files[i]); err != nil {
			return fmt.Errorf("failed to delete backup file: %w", err)
		}
		slog.Info("Deleted old backup file", "file", files[i])
	}

	return nil
}

func ScheduleBackups() error {
	enabled, err := db.Query.GetSettingByKey(context.Background(), "backup-enabled")
	if err != nil {
		return fmt.Errorf("failed to get backup enabled setting: %w", err)
	}
	if enabled.Value != "true" {
		stopBackupCron()
		return nil
	}

	schedule, err := db.Query.GetSettingByKey(context.Background(), "backup-schedule")
	if err != nil {
		return fmt.Errorf("failed to get backup schedule: %w", err)
	}
	if schedule.Value == "" {
		stopBackupCron()
		return nil
	}

	backupCron = cron.New()
	if _, err := backupCron.AddFunc(schedule.Value, func() {
		if err := BackupDatabase(); err != nil {
			slog.Error("Failed to backup database", "error", err)
		}
		if err := CleanupBackups(); err != nil {
			slog.Error("Failed to cleanup backups", "error", err)
		}
	}); err != nil {
		return fmt.Errorf("failed to schedule backup: %w", err)
	}
	backupCron.Start()

	return nil
}

func stopBackupCron() {
	if backupCron != nil {
		backupCron.Stop()
		backupCron = nil
	}
}
