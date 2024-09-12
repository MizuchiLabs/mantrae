// Package config provides functions for parsing command-line f and
// setting up the application's default settings.
package config

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

type Flags struct {
	Version  bool
	Port     int
	URL      string
	Username string
	Password string
	Update   bool
	Reset    bool
}

func (f *Flags) Parse() error {
	flag.BoolVar(&f.Version, "version", false, "Print version and exit")
	flag.IntVar(&f.Port, "port", 3000, "Port to listen on")
	flag.StringVar(
		&f.URL,
		"url",
		"",
		"Specify the URL of the Traefik instance (e.g. http://localhost:8080)",
	)
	flag.StringVar(&f.Username, "username", "", "Specify the username for the Traefik instance")
	flag.StringVar(&f.Password, "password", "", "Specify the password for the Traefik instance")
	flag.BoolVar(&f.Update, "update", false, "Update the application")
	flag.BoolVar(&f.Reset, "reset", false, "Reset the default admin password")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(util.Version)
		os.Exit(0)
	}
	if err := SetDefaultAdminUser(); err != nil {
		return err
	}
	if err := SetDefaultSettings(); err != nil {
		return err
	}

	if f.URL != "" {
		if err := SetDefaultProfile(f.URL, f.Username, f.Password); err != nil {
			return err
		}
	}
	if f.Reset {
		if err := ResetAdminUser(); err != nil {
			return err
		}
	}

	util.UpdateSelf(f.Update)

	return nil
}

func SetDefaultAdminUser() error {
	// check if default admin user exists
	creds, err := db.Query.GetUserByUsername(context.Background(), "admin")
	if err != nil {
		password := util.GenPassword(32)
		hash, err := util.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		if _, err := db.Query.CreateUser(context.Background(), db.CreateUserParams{
			Username: "admin",
			Password: hash,
			Type:     "user",
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
		if _, err := db.Query.UpdateUser(context.Background(), db.UpdateUserParams{
			Username: "admin",
			Password: hash,
			Type:     "user",
		}); err != nil {
			return fmt.Errorf("failed to update default admin user: %w", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
	}
	return nil
}

func SetDefaultProfile(url, username, password string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

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
			return fmt.Errorf("failed to create default profile: %w", err)
		}
		slog.Info("Created default profile", "url", url, "username", username, "password", password)
		return nil
	}

	if profile.Url != url || *profile.Username != username || *profile.Password != password {
		if _, err := db.Query.UpdateProfile(context.Background(), db.UpdateProfileParams{
			ID:       profile.ID,
			Name:     "default",
			Url:      url,
			Username: &username,
			Password: &password,
			Tls:      false,
		}); err != nil {
			return fmt.Errorf("failed to update default profile: %w", err)
		}
		slog.Info("Updated default profile", "url", url, "username", username, "password", password)
	}

	return nil
}

func SetDefaultSettings() error {
	baseSettings := []db.Setting{
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
	}

	for _, setting := range baseSettings {
		if _, err := db.Query.GetSettingByKey(context.Background(), setting.Key); err != nil {
			if _, err := db.Query.CreateSetting(context.Background(), db.CreateSettingParams{
				Key:   setting.Key,
				Value: setting.Value,
			}); err != nil {
				return fmt.Errorf("failed to create setting: %w", err)
			}
		}
	}
	return nil
}

// ResetAdminUser resets the default admin user with a new password.
func ResetAdminUser() error {
	creds, err := db.Query.GetUserByUsername(context.Background(), "admin")
	if err != nil {
		return fmt.Errorf("failed to get default admin user: %w", err)
	}

	password := util.GenPassword(32)
	hash, err := util.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if _, err := db.Query.UpdateUser(context.Background(), db.UpdateUserParams{
		ID:       creds.ID,
		Username: creds.Username,
		Password: hash,
		Type:     "user",
	}); err != nil {
		return fmt.Errorf("failed to update default admin user: %w", err)
	}
	slog.Info("Generated new admin password", "password", password)
	return nil
}
