package util

import (
	"log"
	"os"
	"path/filepath"
)

var (
	MainDir   = ""
	DBName    = "mantrae.db"
	BackupDir = "backups"
)

func Path(rel string) string {
	if MainDir != "" {
		return filepath.Join(MainDir, rel)
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
