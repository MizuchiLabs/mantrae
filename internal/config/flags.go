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
