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
	SM     *SettingsManager
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

	// Setup settings manager
	sm := NewSettingsManager(db.New(DB))
	if err := sm.Initialize(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize settings: %w", err)
	}

	app := App{
		Config: config,
		DB:     db.New(DB),
		BM:     bm,
		SM:     sm,
	}

	app.setupLogger()
	app.setDefaultAdminUser()
	app.setDefaultProfile()
	if flags.Reset {
		app.resetAdminUser()
	}

	// Update self
	util.UpdateSelf(flags.Update)

	// Start background jobs
	app.setupBackgroundJobs(context.Background())
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
