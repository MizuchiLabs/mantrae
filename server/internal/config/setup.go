// Package config provides application configuration and setup.
package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/server/internal/backup"
	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/store"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/logger"
	"github.com/mizuchilabs/mantrae/pkg/util"
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

	// app.Event = events.NewEventBroadcaster()
	app.Conn = store.NewConnection("")
	app.SM = settings.NewManager(app.Conn)
	app.SM.Start(ctx)

	app.BM = backup.NewManager(app.Conn, app.SM)
	app.BM.Start(ctx)

	if err := app.setupDefaultData(ctx); err != nil {
		return nil, err
	}

	// Start background jobs
	app.setupBackgroundJobs(ctx)
	return &app, nil
}

func (a *App) setupDefaultData(ctx context.Context) error {
	q := a.Conn.GetQuery()

	// Ensure at least one user exists
	users, err := q.ListUsers(ctx, db.ListUsersParams{})
	if err != nil {
		return fmt.Errorf("failed to list admin users: %w", err)
	}

	if len(users) == 0 {
		// Generate password if not provided
		password, ok := os.LookupEnv("ADMIN_PASSWORD")
		if !ok {
			password = util.GenPassword(32)
			slog.Info("Generated new admin", "password", password)
		}

		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		email, ok := os.LookupEnv("ADMIN_EMAIL")
		if !ok {
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
		}); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	// Ensure default profile exists
	profiles, err := q.ListProfiles(ctx, db.ListProfilesParams{})
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		description := "Default profile"
		if _, err = q.CreateProfile(ctx, db.CreateProfileParams{
			Name:        "default",
			Description: &description,
			Token:       util.GenerateToken(6),
		}); err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}
	}

	// Check default server url
	ip, err := util.GetHostIPv4()
	if err != nil {
		return fmt.Errorf("failed to get private IPs: %w", err)
	}
	serverURL, ok := a.SM.Get(ctx, "server_url")
	if !ok || serverURL == "" {
		url := fmt.Sprintf("http://%s:3000", ip)
		if err = a.SM.Set(ctx, "server_url", url); err != nil {
			return fmt.Errorf("failed to set server url: %w", err)
		}
	}
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("failed to parse server url: %w", err)
	}
	if u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" {
		url := fmt.Sprintf("http://%s:3000", ip)
		if err := a.SM.Set(ctx, "server_url", url); err != nil {
			return fmt.Errorf("failed to set server url: %w", err)
		}
	}

	return nil
}
