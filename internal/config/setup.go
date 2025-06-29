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
	"github.com/mizuchilabs/mantrae/internal/backup"
	"github.com/mizuchilabs/mantrae/internal/events"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/logger"
	"github.com/mizuchilabs/mantrae/pkg/util"
)

type App struct {
	Secret string `env:"SECRET"`
	Conn   *store.Connection
	Event  *events.EventBroadcaster
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
	app.Event = events.NewEventBroadcaster()
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

	// Ensure at least one admin user exists
	admins, err := q.ListAdminUsers(ctx, db.ListAdminUsersParams{Limit: 1, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to list admin users: %w", err)
	}

	if len(admins) == 0 {
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
		if _, err = q.CreateProfile(ctx, db.CreateProfileParams{
			Name:        "default",
			Description: &description,
		}); err != nil {
			return fmt.Errorf("failed to create default profile: %w", err)
		}
	}

	// Check default server url
	ips, err := util.GetPrivateIPs()
	if err != nil {
		return fmt.Errorf("failed to get private IPs: %w", err)
	}
	serverURL, ok := a.SM.Get("server_url")
	if !ok || serverURL == "" {
		if err = a.SM.Set(ctx, "server_url", ips.IPv4); err != nil {
			return fmt.Errorf("failed to set server url: %w", err)
		}
	}
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("failed to parse server url: %w", err)
	}
	if u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" {
		if err := a.SM.Set(ctx, "server_url", ips.IPv4); err != nil {
			return fmt.Errorf("failed to set server url: %w", err)
		}
	}

	return nil
}
