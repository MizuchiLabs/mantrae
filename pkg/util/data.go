package util

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	Secret          string `env:"SECRET"            envDefault:""`
	Port            string `env:"PORT"              envDefault:"3000"`
	AgentPort       string `env:"AGENT_PORT"        envDefault:"8090"`
	ServerURL       string `env:"SERVER_URL"        envDefault:"http://localhost:8090"`
	EnableBasicAuth bool   `env:"ENABLE_BASIC_AUTH" envDefault:"false"`
	EnableAgent     bool   `env:"ENABLE_AGENT"      envDefault:"true"`
	ConfigDir       string `env:"CONFIG_DIR"        envDefault:""`
	BackupDir       string `env:"BACKUP_DIR"        envDefault:"backups"`
	LogLevel        string `env:"LOG_LEVEL"         envDefault:"info"`

	// Database
	DBType string `env:"DB_TYPE" envDefault:"sqlite"`
	DBName string `env:"DB_NAME" envDefault:"mantrae"`
}

var App AppConfig

func init() {
	if err := env.Parse(&App); err != nil {
		log.Fatal(err)
	}

	if App.Secret == "" {
		log.Fatal("SECRET environment variable not set")
	}
}

func Path(rel string) string {
	if App.ConfigDir != "" {
		return filepath.Join(App.ConfigDir, rel)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(cwd, rel)
}

func DBPath() string {
	if App.DBType == "sqlite" && App.DBName != "" {
		slog.Debug("Using SQLite database", "path", Path(App.DBName+".db"))
		return Path(App.DBName + ".db")
	}

	slog.Error("Invalid database type", "type", App.DBType, "name", App.DBName)
	return ""
}

func BackupPath() string {
	if App.BackupDir == "" {
		slog.Error("BACKUP_DIR environment variable empty")
		return ""
	}
	return Path(App.BackupDir)
}
