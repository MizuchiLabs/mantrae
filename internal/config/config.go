package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Email      EmailConfig
	Traefik    TraefikConfig
	Backup     BackupConfig
	Background BackgroundJobs
	Secret     string `env:"SECRET" envDefault:""`
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port        string `env:"SERVER_PORT" envDefault:"3000"`
	ServerURL   string `env:"SERVER_URL" envDefault:"http://127.0.0.1"`
	BasicAuth   bool   `env:"SERVER_BASIC_AUTH" envDefault:"false"`
	EnableAgent bool   `env:"SERVER_ENABLE_AGENT" envDefault:"true"`
	LogLevel    string `env:"SERVER_LOG_LEVEL" envDefault:"info"`
}

type DatabaseConfig struct {
	Type string `env:"DB_TYPE" envDefault:"sqlite"`
	Name string `env:"DB_NAME" envDefault:"mantrae"`
}

type EmailConfig struct {
	Host     string `env:"EMAIL_HOST" envDefault:"localhost"`
	Port     string `env:"EMAIL_PORT" envDefault:"587"`
	Username string `env:"EMAIL_USERNAME" envDefault:""`
	Password string `env:"EMAIL_PASSWORD" envDefault:""`
	From     string `env:"EMAIL_FROM" envDefault:"mantrae@localhost"`
}

type TraefikConfig struct {
	Profile  string `env:"TRAEFIK_PROFILE" envDefault:"default"`
	URL      string `env:"TRAEFIK_URL" envDefault:""`
	Username string `env:"TRAEFIK_USERNAME" envDefault:""`
	Password string `env:"TRAEFIK_PASSWORD" envDefault:""`
	TLS      bool   `env:"TRAEFIK_TLS" envDefault:""`
}

type BackupConfig struct {
	Enabled  bool   `env:"BACKUP_ENABLED" envDefault:"true"`
	Dir      string `env:"BACKUP_DIR" envDefault:"backups"`
	Schedule string `env:"BACKUP_SCHEDULE" envDefault:"0 2 * * 1"`
	Keep     int    `env:"BACKUP_KEEP" envDefault:"3"`
}

type BackgroundJobs struct {
	Traefik int64 `env:"BACKGROUND_JOBS_TRAEFIK" envDefault:"20"`
	DNS     int64 `env:"BACKGROUND_JOBS_DNS" envDefault:"300"`
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

func Path(rel string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(cwd, rel)
}

func (c *Config) DBPath() string {
	if c.Database.Type == "sqlite" && c.Database.Name != "" {
		slog.Debug("Using SQLite database", "path", Path(c.Database.Name+".db"))
		return Path(c.Database.Name + ".db")
	}

	slog.Error("Invalid database type", "type", c.Database.Type, "name", c.Database.Name)
	return ""
}
