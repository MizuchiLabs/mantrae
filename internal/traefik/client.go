package traefik

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

const (
	RawAPI         = "/api/rawdata"
	EntrypointsAPI = "/api/entrypoints"
	OverviewAPI    = "/api/overview"
	VersionAPI     = "/api/version"
)

func GetTraefikConfig() {
	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
		return
	}

	// TODO: Ignore fetching for tests for now
	if util.IsTest() {
		return
	}

	for _, profile := range profiles {
		if profile.Url == "" {
			continue
		}

		rawResponse, err := fetch(profile, RawAPI)
		if err != nil {
			slog.Error("Failed to fetch raw data", "error", err)
			continue
		}
		defer rawResponse.Close()

		var config runtime.Configuration
		if err := json.NewDecoder(rawResponse).Decode(&config); err != nil {
			slog.Error("Failed to decode raw data", "error", err)
			continue
		}

		epResponse, err := fetch(profile, EntrypointsAPI)
		if err != nil {
			slog.Error("Failed to fetch raw data", "error", err)
			continue
		}
		defer epResponse.Close()

		var entrypoints db.TraefikEntryPoints
		if err := json.NewDecoder(epResponse).Decode(&entrypoints); err != nil {
			slog.Error("Failed to decode raw data", "error", err)
			continue
		}
		oResponse, err := fetch(profile, OverviewAPI)
		if err != nil {
			slog.Error("Failed to fetch raw data", "error", err)
			continue
		}
		defer epResponse.Close()

		var overview db.TraefikOverview
		if err := json.NewDecoder(oResponse).Decode(&overview); err != nil {
			slog.Error("Failed to decode raw data", "error", err)
			continue
		}

		// Check if traefik config exists
		traefikConfig, err := db.Query.GetTraefikConfig(context.Background(), profile.ID)
		if err != nil {
			if err := db.Query.CreateTraefikConfig(context.Background(), db.CreateTraefikConfigParams{
				ProfileID:   profile.ID,
				Entrypoints: &entrypoints,
				Overview:    &overview,
				Config:      &config,
				Source:      "external",
			}); err != nil {
				slog.Error("Failed to upsert traefik api", "error", err)
				continue
			}
		} else {
			if err := db.Query.UpdateTraefikConfig(context.Background(), db.UpdateTraefikConfigParams{
				ID:          traefikConfig.ID,
				Entrypoints: &entrypoints,
				Overview:    &overview,
				Config:      &config,
				Source:      "external",
			}); err != nil {
				slog.Error("Failed to insert traefik config", "error", err)
				continue
			}
		}
	}

	// Broadcast the update to all clients
	util.Broadcast <- util.EventMessage{
		Type:    "profile_updated",
		Message: "Profile updated",
	}
}

func fetch(profile db.Profile, endpoint string) (io.ReadCloser, error) {
	if profile.Url == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	client := http.Client{Timeout: time.Second * 10}
	if !profile.Tls {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("GET", profile.Url+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if profile.Username != nil && profile.Password != nil {
		req.SetBasicAuth(*profile.Username, *profile.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", profile.Url+endpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
