package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/backup"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/logger"
)

type App struct {
	Secret string `env:"SECRET"`
	Conn   *store.Connection
	BM     *backup.BackupManager
	SM     *settings.SettingsManager
}

func Setup(ctx context.Context) (*App, error) {
	// Setup fancy logger
	logger.Setup()

	// Read flags
	ParseFlags()

	// Read environment variables
	app, err := env.ParseAs[App]()
	if err != nil {
		return nil, err
	}

	app.Conn = store.NewConnection("")
	app.BM = backup.NewManager(app.Conn, app.SM)
	app.BM.Start(ctx)

	app.SM = settings.NewManager(app.Conn)
	if err := app.SM.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize settings: %w", err)
	}

	if err := app.setupDefaultData(ctx); err != nil {
		return nil, err
	}

	// Start background jobs
	app.setupBackgroundJobs(ctx)
	return &app, nil
}

func (a *App) setupDefaultData(ctx context.Context) error {
	q := a.Conn.GetQuery()

	// Ensure at least one admin user exists
	admins, err := q.ListAdminUsers(ctx, db.ListAdminUsersParams{Limit: 1, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to list admin users: %w", err)
	}

	if len(admins) == 0 {
		// Generate password if not provided
		password := os.Getenv("ADMIN_PASSWORD")
		if password == "" {
			password = util.GenPassword(32)
			slog.Info("Generated new admin", "password", password)
		}

		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		email := os.Getenv("ADMIN_EMAIL")
		if email == "" {
			email = "admin@localhost"
		}

		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}

		if _, err = q.CreateUser(ctx, db.CreateUserParams{
			ID:       id.String(),
			Username: "admin",
			Password: hash,
			Email:    &email,
			IsAdmin:  true,
		}); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	// Ensure default profile exists
	profiles, err := q.ListProfiles(ctx, db.ListProfilesParams{Limit: 1, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		description := "Default profile"
		if _, err := q.CreateProfile(ctx, db.CreateProfileParams{
			Name:        "default",
			Description: &description,
		}); err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}
	}

	return nil
}
