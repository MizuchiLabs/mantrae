package traefik

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

const (
	RawAPI         = "/api/rawdata"
	EntrypointsAPI = "/api/entrypoints"
	OverviewAPI    = "/api/overview"
	VersionAPI     = "/api/version"
)

func UpdateTraefikAPI(DB *sql.DB, profile db.Profile) error {
	rawResponse, err := fetch(profile, RawAPI)
	if err != nil {
		slog.Error("Failed to fetch raw data", "error", err)

		// Clear api data
		if err = ClearTraefikAPI(DB, profile.ID); err != nil {
			slog.Error("Failed to update api data", "error", err)
		}
		return err
	}
	defer rawResponse.Close()

	var config db.TraefikConfiguration
	if err = json.NewDecoder(rawResponse).Decode(&config); err != nil {
		return fmt.Errorf("failed to decode raw data: %w", err)
	}

	epResponse, err := fetch(profile, EntrypointsAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", profile.Url+EntrypointsAPI, err)
	}
	defer epResponse.Close()

	var entrypoints db.TraefikEntryPoints
	if err = json.NewDecoder(epResponse).Decode(&entrypoints); err != nil {
		return fmt.Errorf("failed to decode entrypoints: %w", err)
	}
	oResponse, err := fetch(profile, OverviewAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", profile.Url+OverviewAPI, err)
	}
	defer epResponse.Close()

	var overview db.TraefikOverview
	if err = json.NewDecoder(oResponse).Decode(&overview); err != nil {
		return fmt.Errorf("failed to decode overview: %w", err)
	}

	vResponse, err := fetch(profile, VersionAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", profile.Url+VersionAPI, err)
	}
	defer vResponse.Close()

	var version db.TraefikVersion
	if err := json.NewDecoder(vResponse).Decode(&version); err != nil {
		return fmt.Errorf("failed to decode version: %w", err)
	}

	q := db.New(DB)
	if err := q.UpsertTraefikConfig(context.Background(), db.UpsertTraefikConfigParams{
		ProfileID:   profile.ID,
		Entrypoints: &entrypoints,
		Overview:    &overview,
		Version:     &version.Version,
		Config:      &config,
		Source:      source.API,
	}); err != nil {
		return fmt.Errorf("failed to update api data: %w", err)
	}

	return nil
}

func ClearTraefikAPI(DB *sql.DB, profileID int64) error {
	q := db.New(DB)
	if err := q.UpsertTraefikConfig(context.Background(), db.UpsertTraefikConfigParams{
		ProfileID:   profileID,
		Source:      source.API,
		Entrypoints: nil,
		Overview:    nil,
		Version:     nil,
		Config:      nil,
	}); err != nil {
		return fmt.Errorf("failed to update api data: %w", err)
	}
	return nil
}

func GetTraefikConfig(DB *sql.DB) {
	q := db.New(DB)
	profiles, err := q.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
		return
	}

	for _, profile := range profiles {
		if profile.Url == "" {
			continue
		}

		if err := UpdateTraefikAPI(DB, profile); err != nil {
			slog.Error("Failed to update api data", "error", err)
			continue
		}
	}

	// Broadcast the update to all clients
	util.Broadcast <- util.EventMessage{
		Type:     util.EventTypeUpdate,
		Category: util.EventCategoryTraefik,
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
