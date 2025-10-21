// Package config provides application configuration and setup.
package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/pkg/logger"
	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/mizuchilabs/mantrae/server/internal/backup"
	"github.com/mizuchilabs/mantrae/server/internal/event"
	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/store"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
)

type App struct {
	Secret string `env:"SECRET"`
	Conn   *store.Connection
	BM     *backup.BackupManager
	SM     *settings.SettingsManager
	Event  *event.Broadcaster
}

func Setup(ctx context.Context) (*App, error) {
	// Setup fancy logger
	logger.Setup()

	// Read flags
	flags := ParseFlags()

	// Read environment variables
	app, err := env.ParseAs[App]()
	if err != nil {
		return nil, err
	}

	if len(app.Secret) != 16 && len(app.Secret) != 24 && len(app.Secret) != 32 {
		return nil, fmt.Errorf("secret must be either 16, 24 or 32 bytes")
	}

	app.Conn = store.NewConnection("")
	app.SM = settings.NewManager(app.Conn)
	app.SM.Start(ctx)

	app.BM = backup.NewManager(app.Conn, app.SM)
	app.BM.Start(ctx)

	app.Event = event.NewBroadcaster(ctx)

	app.resetPassword(ctx, flags)

	return &app, app.setupDefaultData(ctx)
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

func (a *App) resetPassword(ctx context.Context, flags *Flags) {
	if flags.ResetPassword == "" {
		return
	}

	user, err := a.Conn.GetQuery().GetUserByUsername(ctx, flags.ResetUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("failed to get user", "user", flags.ResetUser)
		} else {
			slog.Error("failed to get user", "error", err)
		}
		os.Exit(1)
	}
	hash, err := util.HashPassword(flags.ResetPassword)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		os.Exit(1)
	}
	if err = a.Conn.GetQuery().UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: hash,
	}); err != nil {
		slog.Error("failed to update password for user", "user", flags.ResetUser, "error", err)
		os.Exit(1)
	}

	slog.Info("Reset successful!", "user", flags.ResetUser, "password", flags.ResetPassword)
	os.Exit(1)
}
