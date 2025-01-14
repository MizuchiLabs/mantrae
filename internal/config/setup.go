package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/lmittmann/tint"
)

type App struct {
	Config *Config
	DB     *db.Queries
	BM     *BackupManager
}

func Setup() (*App, error) {
	// Initialize database
	DB, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	// Read environment variables
	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	// Read flags
	flags, err := ParseFlags()
	if err != nil {
		return nil, err
	}

	// Setup backup manager
	bm, err := NewBackupManager(*config, db.New(DB))
	if err != nil {
		return nil, fmt.Errorf("failed to create backup manager: %w", err)
	}

	app := App{
		Config: config,
		DB:     db.New(DB),
		BM:     bm,
	}

	// Start background jobs
	app.setupBackgroundJobs(context.Background())

	app.setupLogger()
	app.setDefaultAdminUser()
	app.setDefaultSettings()
	app.setDefaultProfile()
	if flags.Reset {
		app.resetAdminUser()
	}

	return &app, nil
}

func (a *App) setupLogger() {
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logLevel, exists := levelMap[strings.ToLower(a.Config.Server.LogLevel)]
	if !exists {
		logLevel = slog.LevelInfo
	}

	opts := &tint.Options{Level: logLevel}
	logger := slog.New(tint.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func (a *App) setDefaultAdminUser() error {
	// check if default admin user exists
	creds, err := a.DB.GetUserByUsername(context.Background(), "admin")
	if err != nil {
		password := util.GenPassword(32)
		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		if err := a.DB.CreateUser(context.Background(), db.CreateUserParams{
			Username: "admin",
			Password: hash,
			IsAdmin:  true,
		}); err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
		return nil
	}

	// Validate credentials
	if creds.Username != "admin" || creds.Password == "" {
		password := util.GenPassword(32)
		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		slog.Info("Invalid credentials, regenerating...")
		if err := a.DB.UpdateUser(context.Background(), db.UpdateUserParams{
			Username: "admin",
			Password: hash,
			IsAdmin:  true,
		}); err != nil {
			return fmt.Errorf("failed to update default admin user: %w", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
	}
	return nil
}

func (a *App) setDefaultProfile() error {
	if a.Config.Traefik.URL == "" {
		return nil
	}

	if !strings.HasPrefix(a.Config.Traefik.URL, "http://") && !strings.HasPrefix(a.Config.Traefik.URL, "https://") {
		a.Config.Traefik.URL = "http://" + a.Config.Traefik.URL
	}

	_, err := a.DB.GetProfileByName(context.Background(), a.Config.Traefik.Profile)
	if err != nil {
		err := a.DB.CreateProfile(context.Background(), db.CreateProfileParams{
			Name:     a.Config.Traefik.Profile,
			Url:      a.Config.Traefik.URL,
			Username: &a.Config.Traefik.Username,
			Password: &a.Config.Traefik.Password,
			Tls:      a.Config.Traefik.TLS,
		})
		if err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}
		slog.Info("Created default profile", "url", a.Config.Traefik.URL, "username", a.Config.Traefik.Username, "password", a.Config.Traefik.Password)
		return nil
	}

	if err := a.DB.UpdateProfile(context.Background(), db.UpdateProfileParams{
		Name:     a.Config.Traefik.Profile,
		Url:      a.Config.Traefik.URL,
		Username: &a.Config.Traefik.Username,
		Password: &a.Config.Traefik.Password,
		Tls:      a.Config.Traefik.TLS,
	}); err != nil {
		return fmt.Errorf("failed to update default profile: %w", err)
	}
	slog.Info("Updated default profile", "url", a.Config.Traefik.URL, "username", a.Config.Traefik.Username, "password", a.Config.Traefik.Password)
	return nil
}

func (a *App) setDefaultSettings() error {
	baseSettings := []db.Setting{
		{
			Key:   "server-url",
			Value: a.Config.Server.ServerURL,
		},
		{
			Key:   "backup-enabled",
			Value: "true",
		},
		{
			Key:   "backup-schedule",
			Value: "0 2 * * 1", // Weekly at 02:00 AM on Monday
		},
		{
			Key:   "backup-keep",
			Value: "3", // Keep 3 backups
		},
		{
			Key:   "agent-cleanup-enabled",
			Value: "true",
		},
		{
			Key:   "agent-cleanup-timeout",
			Value: "168h",
		},
		{
			Key:   "email-host",
			Value: a.Config.Email.Host,
		},
		{
			Key:   "email-port",
			Value: a.Config.Email.Port,
		},
		{
			Key:   "email-username",
			Value: a.Config.Email.Username,
		},
		{
			Key:   "email-password",
			Value: a.Config.Email.Password,
		},
		{
			Key:   "email-from",
			Value: a.Config.Email.From,
		},
	}

	for _, setting := range baseSettings {
		if err := a.DB.UpsertSetting(context.Background(), db.UpsertSettingParams{
			Key:   setting.Key,
			Value: setting.Value,
		}); err != nil {
			return fmt.Errorf("failed to create setting: %w", err)
		}
	}
	return nil
}

// ResetAdminUser resets the default admin user with a new password.
func (a *App) resetAdminUser() error {
	creds, err := a.DB.GetUserByUsername(context.Background(), "admin")
	if err != nil {
		return fmt.Errorf("failed to get default admin user: %w", err)
	}

	password := util.GenPassword(32)
	hash, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := a.DB.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: creds.Username,
		Password: hash,
		IsAdmin:  true,
	}); err != nil {
		return fmt.Errorf("failed to update default admin user: %w", err)
	}
	slog.Info("Reset default admin user", "username", "admin", "password", password)
	return nil
}
