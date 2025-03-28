package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/app"
	"github.com/MizuchiLabs/mantrae/internal/backup"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/storage"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/lmittmann/tint"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Config *app.Config
	Conn   *db.Connection
	BM     *backup.BackupManager
	SM     *SettingsManager
}

func Setup(ctx context.Context) (*App, error) {
	// Read flags
	ParseFlags()

	// Read environment variables
	config, err := app.ReadConfig()
	if err != nil {
		return nil, err
	}

	conn, err := db.NewDBConnection()
	if err != nil {
		if conn != nil {
			defer conn.Close()
		}
		return nil, err
	}

	storage, err := storage.NewLocalStorage(config.Backup.BackupPath)
	if err != nil {
		return nil, err
	}
	bm := backup.NewManager(conn, config.Backup, storage)
	bm.Start(ctx)

	// Setup settings manager
	sm := NewSettingsManager(conn.Get())
	if err := sm.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize settings: %w", err)
	}

	app := App{
		Config: config,
		Conn:   conn,
		BM:     bm,
		SM:     sm,
	}

	app.setupLogger()
	if err := app.setDefaultAdminUser(ctx); err != nil {
		return nil, err
	}
	if err := app.setDefaultProfile(ctx); err != nil {
		return nil, err
	}

	// Start background jobs
	app.setupBackgroundJobs(ctx)
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

func (a *App) setDefaultAdminUser(ctx context.Context) error {
	// Generate password if not provided
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = util.GenPassword(32)
		slog.Info(
			"Generated new admin password (Please use ADMIN_PASSWORD env var to set it)",
			"password",
			password,
		)
	}

	hash, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Try to get existing admin user
	q := a.Conn.GetQuery()
	user, err := q.GetUser(ctx, 1)
	// If admin doesn't exist, create new admin
	if err != nil {
		adminMail := "admin@mantrae"
		if err = q.CreateUser(ctx, db.CreateUserParams{
			Username: "admin",
			Email:    &adminMail,
			Password: hash,
			IsAdmin:  true,
		}); err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}
		slog.Info("Generated default 'admin' user")
		return nil
	}

	userPassword, err := q.GetUserPassword(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user password: %w", err)
	}
	// Skip if password is correct
	passwordErr := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if passwordErr == nil {
		return nil
	} else {
		// Update password on change
		if err = q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			ID:       user.ID,
			Password: hash,
		}); err != nil {
			return fmt.Errorf("failed to update default admin user password: %w", err)
		}
	}

	return nil
}

func (a *App) setDefaultProfile(ctx context.Context) error {
	if a.Config.Traefik.URL == "" {
		return nil
	}

	if !strings.HasPrefix(a.Config.Traefik.URL, "http://") &&
		!strings.HasPrefix(a.Config.Traefik.URL, "https://") {
		a.Config.Traefik.URL = "http://" + a.Config.Traefik.URL
	}

	q := a.Conn.GetQuery()
	profile, err := q.GetProfileByName(ctx, a.Config.Traefik.Profile)
	if err != nil {
		profileID, err := q.CreateProfile(ctx, db.CreateProfileParams{
			Name:     a.Config.Traefik.Profile,
			Url:      a.Config.Traefik.URL,
			Username: &a.Config.Traefik.Username,
			Password: &a.Config.Traefik.Password,
			Tls:      a.Config.Traefik.TLS,
		})
		if err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}

		// Create default local config
		if err := q.UpsertTraefikConfig(ctx, db.UpsertTraefikConfigParams{
			ProfileID: profileID,
			Source:    source.Local,
		}); err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}

		slog.Info(
			"Created default profile",
			"url",
			a.Config.Traefik.URL,
			"username",
			a.Config.Traefik.Username,
			"password",
			a.Config.Traefik.Password,
		)
		return nil
	}

	if profile.Url == a.Config.Traefik.URL &&
		*profile.Username == a.Config.Traefik.Username &&
		*profile.Password == a.Config.Traefik.Password &&
		profile.Tls == a.Config.Traefik.TLS {
		return nil
	}

	if err := q.UpdateProfile(ctx, db.UpdateProfileParams{
		ID:       profile.ID,
		Name:     a.Config.Traefik.Profile,
		Url:      a.Config.Traefik.URL,
		Username: &a.Config.Traefik.Username,
		Password: &a.Config.Traefik.Password,
		Tls:      a.Config.Traefik.TLS,
	}); err != nil {
		return fmt.Errorf("failed to update default profile: %w", err)
	}
	slog.Info(
		"Updated default profile",
		"url",
		a.Config.Traefik.URL,
		"username",
		a.Config.Traefik.Username,
		"password",
		a.Config.Traefik.Password,
	)
	return nil
}
