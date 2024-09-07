// Package config provides functions for parsing command-line flags and
// setting up the application's default settings.
package config

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

type Flags struct {
	Version  bool
	Port     int
	URL      string
	Username string
	Password string
}

func ParseFlags() *Flags {
	var flags Flags

	flag.BoolVar(&flags.Version, "version", false, "Print version and exit")
	flag.IntVar(&flags.Port, "port", 3000, "Port to listen on")
	flag.StringVar(
		&flags.URL,
		"url",
		"",
		"Specify the URL of the Traefik instance (e.g. http://localhost:8080)",
	)
	flag.StringVar(&flags.Username, "username", "", "Specify the username for the Traefik instance")
	flag.StringVar(&flags.Password, "password", "", "Specify the password for the Traefik instance")

	flag.Parse()

	if flags.Version {
		fmt.Println(util.Version)
		os.Exit(0)
	}

	if flags.URL != "" {
		SetDefaultProfile(flags.URL, flags.Username, flags.Password)
	}

	return &flags
}

func SetDefaultAdminUser() {
	// check if default admin user exists
	creds, err := db.Query.GetUserByUsername(context.Background(), "admin")
	if err != nil {
		password := util.GenPassword(32)
		hash, err := util.HashPassword(password)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			return
		}

		if _, err := db.Query.CreateUser(context.Background(), db.CreateUserParams{
			Username: "admin",
			Password: hash,
			Type:     "user",
		}); err != nil {
			slog.Error("Failed to create default admin user", "error", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
		return
	}

	// Validate credentials
	if creds.Username != "admin" || creds.Password == "" {
		password := util.GenPassword(32)
		hash, err := util.HashPassword(password)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			return
		}
		slog.Info("Invalid credentials, regenerating...")
		if _, err := db.Query.UpdateUser(context.Background(), db.UpdateUserParams{
			Username: "admin",
			Password: hash,
			Type:     "user",
		}); err != nil {
			slog.Error("Failed to update default admin user", "error", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
	}
}

func SetDefaultProfile(url, username, password string) {
	profile, err := db.Query.GetProfileByName(context.Background(), "default")
	if err != nil {
		_, err := db.Query.CreateProfile(context.Background(), db.CreateProfileParams{
			Name:     "default",
			Url:      url,
			Username: &username,
			Password: &password,
			Tls:      false,
		})
		if err != nil {
			slog.Error("Failed to create default profile", "error", err)
		}
		slog.Info("Generated default profile", "url", url)
		return
	}
	if profile.Url != url || profile.Username != &username || profile.Password != &password {
		if _, err := db.Query.UpdateProfile(context.Background(), db.UpdateProfileParams{
			ID:       profile.ID,
			Name:     "default",
			Url:      url,
			Username: &username,
			Password: &password,
			Tls:      false,
		}); err != nil {
			slog.Error("Failed to update default profile", "error", err)
		}
	}
}

func SetDefaultSettings() {
	baseSettings := []db.Setting{
		{
			Key:   "backup-schedule",
			Value: "0 2 * * 1", // Weekly at 02:00 AM on Monday
		},
		{
			Key:   "backup-keep",
			Value: "3", // Keep 3 backups
		},
	}

	for _, setting := range baseSettings {
		if _, err := db.Query.GetSettingByKey(context.Background(), setting.Key); err != nil {
			if _, err := db.Query.CreateSetting(context.Background(), db.CreateSettingParams{
				Key:   setting.Key,
				Value: setting.Value,
			}); err != nil {
				slog.Error("Failed to create setting", "error", err)
			}
		}
	}
}
