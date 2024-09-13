package util

import (
	"log"
	"os"
	"path/filepath"
)

var (
	DBName    = "mantrae.db"
	BackupDir = "backups"
)

func Path(rel string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	mainDir := filepath.Join(configDir, "mantrae")
	if err = os.MkdirAll(mainDir, 0700); err != nil {
		log.Fatal(err)
	}

	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(mainDir, rel)
}
