package traefik

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func profilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(cwd, "profiles.json")
}

func (p *Profiles) Load() error {
	p.mu.RLock()

	if _, err := os.Stat(profilePath()); os.IsNotExist(err) {
		p.Profiles = make(map[string]Profile)
		p.Profiles["default"] = Profile{Name: "default"}
		p.mu.RUnlock()
		if err := p.Save(); err != nil {
			slog.Error("Failed to save profiles", "error", err)
		}
		return nil
	}

	file, err := os.ReadFile(profilePath())
	if err != nil {
		return fmt.Errorf("failed to read profiles file: %w", err)
	}

	if err := json.Unmarshal(file, &p.Profiles); err != nil {
		return fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	p.mu.RUnlock()
	return nil
}

func (p *Profiles) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	tmpFile, err := os.CreateTemp(os.TempDir(), "profiles-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	profileBytes, err := json.Marshal(p.Profiles)
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	_, err = tmpFile.Write(profileBytes)
	if err != nil {
		return fmt.Errorf("failed to write profiles: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	tmpFile.Close()

	if err := Move(tmpFile.Name(), profilePath()); err != nil {
		return fmt.Errorf("failed to move temp file: %w", err)
	}

	return nil
}

func (p *Profiles) SetDefaultProfile(url, username, password string) error {
	if err := p.Load(); err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}
	if len(p.Profiles) == 0 {
		p.Profiles = make(map[string]Profile)
	}

	p.Profiles["default"] = Profile{
		Name:     "default",
		URL:      url,
		Username: username,
		Password: password,
	}

	return p.Save()
}

func Move(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(source, destination)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return fmt.Errorf("failed to stat source file: %w", err)
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return fmt.Errorf("failed to chmod destination file: %w", err)
	}
	os.Remove(source)
	return nil
}
