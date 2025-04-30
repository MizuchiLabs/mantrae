package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration
type Config struct {
	Server     ServerConfig
	Traefik    TraefikConfig
	Background BackgroundJobs
	Secret     string `env:"SECRET" envDefault:""`
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST"         envDefault:"0.0.0.0"`
	Port        string `env:"SERVER_PORT"         envDefault:"3000"`
	ServerURL   string `env:"SERVER_URL"          envDefault:"http://127.0.0.1"`
	BasicAuth   bool   `env:"SERVER_BASIC_AUTH"   envDefault:"false"`
	EnableAgent bool   `env:"SERVER_ENABLE_AGENT" envDefault:"true"`
	LogLevel    string `env:"SERVER_LOG_LEVEL"    envDefault:"info"`
}

type TraefikConfig struct {
	Profile  string `env:"TRAEFIK_PROFILE"  envDefault:"default"`
	URL      string `env:"TRAEFIK_URL"      envDefault:""`
	Username string `env:"TRAEFIK_USERNAME" envDefault:""`
	Password string `env:"TRAEFIK_PASSWORD" envDefault:""`
	TLS      bool   `env:"TRAEFIK_TLS"      envDefault:""`
}

type BackgroundJobs struct {
	Traefik int64 `env:"BACKGROUND_JOBS_TRAEFIK" envDefault:"20"`
	DNS     int64 `env:"BACKGROUND_JOBS_DNS"     envDefault:"180"`
	Agent   int64 `env:"BACKGROUND_JOBS_AGENT"   envDefault:"180"`
}

func ReadConfig() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	if config.Secret == "" {
		return nil, fmt.Errorf("SECRET environment variable not set")
	}

	return &config, nil
}

func ResolvePath(path string) string {
	var basePath string
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		basePath = dbPath
	} else {
		basePath = "data"
	}

	// If the provided path is absolute, return it as-is
	if filepath.IsAbs(path) {
		return path
	}

	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Printf("Warning: failed to create base directory: %v", err)
	}

	return filepath.Join(basePath, path)
}
