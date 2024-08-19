package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   []byte `json:"secret"`
}

func randomPassword(length int) []byte {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil
	}
	return bytes
}

func GenerateCreds() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	credsPath := filepath.Join(cwd, "creds.json")
	if _, err = os.Stat(credsPath); os.IsNotExist(err) {
		username := "admin"
		password := base64.StdEncoding.EncodeToString(randomPassword(32))

		jsonCreds, err := json.MarshalIndent(Credentials{
			Username: username,
			Password: password,
			Secret:   randomPassword(64),
		}, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(credsPath, jsonCreds, 0644); err != nil {
			return err
		}

		slog.Info("Generated new credentials", "username", username, "password", password)
	}

	return nil
}

// GetCreds retrieves the credentials from the creds.json file or generates a new one
func (c *Credentials) GetCreds() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	credsPath := filepath.Join(cwd, "creds.json")

	file, err := os.ReadFile(credsPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &c); err != nil {
		return err
	}

	if c.Username == "" || c.Password == "" || c.Secret == nil {
		slog.Info("Invalid credentials, regenerating...")
		if err := os.Remove(credsPath); err != nil {
			return err
		}
		if err := GenerateCreds(); err != nil {
			return err
		}
	}
	return nil
}
