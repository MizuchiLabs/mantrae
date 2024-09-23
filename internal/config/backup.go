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

// BackupData is the structure for the full manual backup
type BackupData struct {
	Profiles  []db.Profile  `json:"profiles"`
	Configs   []db.Config   `json:"configs"`
	Providers []db.Provider `json:"providers"`
	Settings  []db.Setting  `json:"settings"`
	Users     []db.User     `json:"users"`
}

func DumpBackup(ctx context.Context) (*BackupData, error) {
	var data BackupData
	var err error

	data.Profiles, err = db.Query.ListProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %w", err)
	}

	data.Configs, err = db.Query.ListConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get configs: %w", err)
	}

	data.Providers, err = db.Query.ListProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get providers: %w", err)
	}

	data.Settings, err = db.Query.ListSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	data.Users, err = db.Query.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &data, nil
}

func RestoreBackup(ctx context.Context, data *BackupData) error {
	for _, profile := range data.Profiles {
		if _, err := db.Query.UpsertProfile(ctx, db.UpsertProfileParams(profile)); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
	}

	for _, config := range data.Configs {
		if _, err := db.Query.UpsertConfig(ctx, db.UpsertConfigParams(config)); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
	}

	for _, provider := range data.Providers {
		if _, err := db.Query.UpsertProvider(ctx, db.UpsertProviderParams(provider)); err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
		}
	}

	for _, setting := range data.Settings {
		if _, err := db.Query.UpsertSetting(ctx, db.UpsertSettingParams(setting)); err != nil {
			return fmt.Errorf("failed to create setting: %w", err)
		}
	}

	for _, user := range data.Users {
		if _, err := db.Query.UpsertUser(ctx, db.UpsertUserParams(user)); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}
	return nil
}

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
