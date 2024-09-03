// Package util contains various utility functions
package util

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

func randomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}

func SetDefaultAdminUser() {
	// check if default admin user exists
	creds, err := db.Query.GetCredentialByUsername(context.Background(), "admin")
	if err != nil {
		password := randomPassword(32)
		if err := db.Query.CreateCredential(context.Background(), db.CreateCredentialParams{
			Username: "admin",
			Password: password,
		}); err != nil {
			slog.Error("Failed to create default admin user", "error", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
		return
	}

	// Validate credentials
	if creds.Username != "admin" || creds.Password == "" {
		password := randomPassword(32)
		slog.Info("Invalid credentials, regenerating...")
		if err := db.Query.UpdateCredential(context.Background(), db.UpdateCredentialParams{
			Username: "admin",
			Password: password,
		}); err != nil {
			slog.Error("Failed to update default admin user", "error", err)
		}
		slog.Info("Generated default admin user", "username", "admin", "password", password)
	}
}
