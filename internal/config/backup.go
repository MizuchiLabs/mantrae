package config

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/robfig/cron/v3"
)

var backupCron *cron.Cron

// BackupData is the structure for the full manual backup
type BackupData struct {
	Profiles    []db.Profile    `json:"profiles"`
	Providers   []db.Provider   `json:"providers"`
	Settings    []db.Setting    `json:"settings"`
	Users       []db.User       `json:"users"`
	Routers     []db.Router     `json:"routers"`
	Services    []db.Service    `json:"services"`
	Middlewares []db.Middleware `json:"middlewares"`
}

func DumpBackup(ctx context.Context) (*BackupData, error) {
	var data BackupData
	var err error

	data.Profiles, err = db.Query.ListProfiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %w", err)
	}

	data.Routers, err = db.Query.ListRoutersByProvider(ctx, "http")
	if err != nil {
		return nil, fmt.Errorf("failed to get routers: %w", err)
	}
	data.Services, err = db.Query.ListServicesByProvider(ctx, "http")
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	data.Middlewares, err = db.Query.ListMiddlewaresByProvider(ctx, "http")
	if err != nil {
		return nil, fmt.Errorf("failed to get middlewares: %w", err)
	}

	for i, router := range data.Routers {
		if err := router.DecodeFields(); err != nil {
			slog.Error("Failed to decode router", "name", router.Name, "error", err)
		}
		data.Routers[i] = router
	}
	for i, service := range data.Services {
		if err := service.DecodeFields(); err != nil {
			slog.Error("Failed to decode service", "name", service.Name, "error", err)
		}
		data.Services[i] = service
	}
	for i, middleware := range data.Middlewares {
		if err := middleware.DecodeFields(); err != nil {
			slog.Error("Failed to decode middleware", "name", middleware.Name, "error", err)
		}
		data.Middlewares[i] = middleware
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
	var err error
	db.DB.Close()

	// Delete Database
	if err = os.Remove(util.DBPath()); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete existing database: %w", err)
	}

	// Reopen the database connection
	if err = db.InitDB(); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create defaults
	if err := SetDefaultAdminUser(); err != nil {
		return err
	}
	if err := SetDefaultSettings(); err != nil {
		return err
	}

	for _, provider := range data.Providers {
		if _, err := db.Query.UpsertProvider(ctx, db.UpsertProviderParams(provider)); err != nil {
			return fmt.Errorf("failed to upsert provider: %w", err)
		}
	}

	for _, setting := range data.Settings {
		if _, err := db.Query.UpsertSetting(ctx, db.UpsertSettingParams(setting)); err != nil {
			return fmt.Errorf("failed to upsert setting: %w", err)
		}
	}

	for _, user := range data.Users {
		if _, err := db.Query.UpsertUser(ctx, db.UpsertUserParams(user)); err != nil {
			return fmt.Errorf("failed to upsert user: %w", err)
		}
	}

	// Insert data
	if len(data.Profiles) > 0 {
		for _, profile := range data.Profiles {
			if _, err := db.Query.UpsertProfile(ctx, db.UpsertProfileParams(profile)); err != nil {
				return fmt.Errorf("failed to upsert profile: %w", err)
			}
		}

		for _, router := range data.Routers {
			if err := router.Verify(); err != nil {
				continue
			}
			if _, err := db.Query.UpsertRouter(ctx, db.UpsertRouterParams(router)); err != nil {
				return fmt.Errorf("failed to upsert router: %w", err)
			}
		}

		for _, service := range data.Services {
			params := db.UpsertServiceParams(service)
			if err := params.Verify(); err != nil {
				continue
			}
			if _, err := db.Query.UpsertService(ctx, db.UpsertServiceParams(params)); err != nil {
				return fmt.Errorf("failed to upsert service: %w", err)
			}
		}

		for _, middleware := range data.Middlewares {
			params := db.UpsertMiddlewareParams(middleware)
			if err := params.Verify(); err != nil {
				continue
			}
			if _, err := db.Query.UpsertMiddleware(ctx, db.UpsertMiddlewareParams(params)); err != nil {
				return fmt.Errorf("failed to upsert middleware: %w", err)
			}
		}
	}

	// Trigger fetch
	if len(data.Profiles) > 0 {
		go traefik.GetTraefikConfig()
	}

	return nil
}

func BackupDatabase() error {
	timestamp := time.Now().Format("2006-01-02")
	backupName := fmt.Sprintf("backup-%s.tar.gz", timestamp)
	backupPath := fmt.Sprintf("%s/%s", util.BackupPath(), backupName)

	// Create the backup directory if it doesn't exist
	backupDir := filepath.Dir(backupPath)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(backupDir, 0750); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	// Check if the backup file already exists
	if _, err := os.Stat(backupPath); err == nil {
		return nil // Backup already exists
	}

	// Open the original database file
	dbFile, err := os.Open(util.DBPath())
	if err != nil {
		return fmt.Errorf("failed to open original database file: %w", err)
	}
	defer dbFile.Close()

	// Create the gzip file
	gzipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create gzip file: %w", err)
	}
	defer gzipFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	info, err := dbFile.Stat()
	if err != nil {
		return err
	}

	// Add the original database file to the tar archive
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return fmt.Errorf("failed to get file info header: %w", err)
	}

	header.Name = util.DBPath()
	if err = tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err = io.Copy(tarWriter, dbFile); err != nil {
		return fmt.Errorf("failed to copy database file: %w", err)
	}

	// Generate the JSON backup using DumpBackup
	backupData, err := DumpBackup(context.Background())
	if err != nil {
		return fmt.Errorf("failed to dump backup data: %w", err)
	}

	// Marshal the JSON data
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	// Add the JSON file to tar archive
	jsonFileName := "backup.json"
	jsonHeader := &tar.Header{
		Name: jsonFileName,
		Size: int64(len(jsonData)),
		Mode: 0644,
	}
	if err := tarWriter.WriteHeader(jsonHeader); err != nil {
		return fmt.Errorf("failed to write JSON header: %w", err)
	}
	if _, err := tarWriter.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
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
	files, err := filepath.Glob(fmt.Sprintf("%s/backup-*.tar.gz", util.BackupPath()))
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
