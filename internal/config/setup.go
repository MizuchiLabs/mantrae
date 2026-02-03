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
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/urfave/cli/v3"
)

type EnvConfig struct {
	Secret        string `env:"SECRET"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	AdminEmail    string `env:"ADMIN_EMAIL"`
	BaseURL       string `env:"BASE_URL"       envDefault:"http://localhost:3000"`
	FrontendURL   string `env:"FRONTEND_URL"   envDefault:"http://localhost:5173"`
	Debug         bool   `env:"DEBUG"          envDefault:"false"`
	Version       string
}

type App struct {
	// Environment variables
	EnvConfig

	// App state
	Conn *store.Connection
	BM   *backup.BackupManager
	SM   *settings.SettingsManager
	DNS  *dns.DNSManager
}

func New(ctx context.Context, cmd *cli.Command) (*App, error) {
	app, err := env.ParseAs[App]()
	if err != nil {
		return nil, err
	}

	// Merge command line flags (can override environment variables)
	if cmd != nil {
		app.Debug = cmd.Bool("debug")
		app.Version = cmd.Root().Version
	}

	app.initLogger()

	if !slices.Contains([]int{16, 24, 32}, len(app.Secret)) {
		return nil, fmt.Errorf("secret must be either 16, 24 or 32 bytes, is %d", len(app.Secret))
	}

	app.Conn = store.NewConnection(ctx, "")
	app.SM = settings.NewManager(app.Conn)
	app.SM.Start(ctx)

	app.BM = backup.NewManager(app.Conn, app.SM)
	app.BM.Start(ctx)

	app.DNS = dns.NewManager(app.Conn, app.Secret)

	return &app, app.setupDefaultData(ctx)
}

func (a App) initLogger() {
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
	q := a.Conn.Q

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

func (a App) BaseHost() string {
	host, _ := util.ParseHostPort(a.BaseURL)
	return host
}

func (a App) BasePort() string {
	_, port := util.ParseHostPort(a.BaseURL)
	return port
}

func (a App) BaseDomain() string {
	u, err := url.Parse(a.BaseURL)
	if err != nil {
		return ""
	}
	return u.Host // includes port if present
}

// --- FrontendURL helpers ---

func (a App) FrontendHost() string {
	host, _ := util.ParseHostPort(a.FrontendURL)
	return host
}

func (a App) FrontendPort() string {
	_, port := util.ParseHostPort(a.FrontendURL)
	return port
}

func (a App) FrontendDomain() string {
	u, err := url.Parse(a.FrontendURL)
	if err != nil {
		return ""
	}
	return u.Host
}
