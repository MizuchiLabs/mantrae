// Package config provides application configuration and setup.
package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mizuchilabs/mantrae/internal/backup"
	"github.com/mizuchilabs/mantrae/internal/dns"
	"github.com/mizuchilabs/mantrae/internal/event"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type EnvConfig struct {
	Secret        string `env:"SECRET"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	AdminEmail    string `env:"ADMIN_EMAIL"`
	BaseURL       string `env:"BASE_URL"       envDefault:"http://localhost:3000"`
	FrontendURL   string `env:"FRONTEND_URL"   envDefault:"http://localhost:5173"`
	Debug         bool   `env:"DEBUG"          envDefault:"false"`
}

type App struct {
	// Environment variables
	EnvConfig

	// App state
	Conn  *store.Connection
	BM    *backup.BackupManager
	SM    *settings.SettingsManager
	Event *event.Broadcaster
	DNS   *dns.DNSManager
}

func Setup(ctx context.Context) (*App, error) {
	// Read environment variables
	app, err := env.ParseAs[App]()
	if err != nil {
		return nil, err
	}
	// Setup logger
	app.Logger()

	if !slices.Contains([]int{16, 24, 32}, len(app.Secret)) {
		return nil, fmt.Errorf("secret must be either 16, 24 or 32 bytes, is %d", len(app.Secret))
	}

	app.Conn = store.NewConnection(ctx, "")
	app.SM = settings.NewManager(app.Conn)
	app.SM.Start(ctx)

	app.BM = backup.NewManager(app.Conn, app.SM)
	app.BM.Start(ctx)

	app.Event = event.NewBroadcaster(ctx)
	app.DNS = dns.NewManager(app.Conn, app.Secret)

	return &app, app.setupDefaultData(ctx)
}

func (a *App) Logger() {
	level := slog.LevelInfo
	if a.Debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(colorable.NewColorable(os.Stderr), &tint.Options{
			Level:      level,
			TimeFormat: time.RFC3339,
			NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		}),
	))
}

func (a *App) setupDefaultData(ctx context.Context) error {
	q := a.Conn.Query

	// Ensure at least one user exists
	users, err := q.ListUsers(ctx, &db.ListUsersParams{})
	if err != nil {
		return fmt.Errorf("failed to list admin users: %w", err)
	}

	if len(users) == 0 {
		// Generate password if not provided
		password := a.AdminPassword
		if password == "" {
			password = util.GenPassword(32)
			slog.Info("Generated new admin", "password", password)
		}

		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		email := a.AdminEmail
		if email == "" {
			email = "admin@localhost"
		}

		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}

		if _, err = q.CreateUser(ctx, &db.CreateUserParams{
			ID:       id.String(),
			Username: "admin",
			Password: hash,
			Email:    &email,
		}); err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	// Ensure default profile exists
	profiles, err := q.ListProfiles(ctx, &db.ListProfilesParams{})
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		description := "Default profile"
		if _, err = q.CreateProfile(ctx, &db.CreateProfileParams{
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

// --- BaseURL helpers ---

func (c App) BaseHost() string {
	host, _ := util.ParseHostPort(c.BaseURL)
	return host
}

func (c App) BasePort() string {
	_, port := util.ParseHostPort(c.BaseURL)
	return port
}

func (c App) BaseDomain() string {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return ""
	}
	return u.Host // includes port if present
}

// --- FrontendURL helpers ---

func (c App) FrontendHost() string {
	host, _ := util.ParseHostPort(c.FrontendURL)
	return host
}

func (c App) FrontendPort() string {
	_, port := util.ParseHostPort(c.FrontendURL)
	return port
}

func (c App) FrontendDomain() string {
	u, err := url.Parse(c.FrontendURL)
	if err != nil {
		return ""
	}
	return u.Host
}
