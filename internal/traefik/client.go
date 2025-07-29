// Package traefik provides functionality for interacting with the Traefik API.
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

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/pkg/meta"
)

const (
	RawAPI         = "/api/rawdata"
	EntrypointsAPI = "/api/entrypoints"
	OverviewAPI    = "/api/overview"
	VersionAPI     = "/api/version"
)

type APIResponse struct {
	Name string
	Data io.ReadCloser
	Err  error
}

func UpdateTraefikInstance(r *http.Request, q *db.Queries, profileID int64) {
	instanceName := r.Header.Get(meta.HeaderTraefikName)
	instanceURL := r.Header.Get(meta.HeaderTraefikURL)
	if instanceName == "" || instanceURL == "" {
		slog.Debug("Skipping traefik update, missing headers")
		return
	}

	instance, err := q.GetTraefikInstanceByName(
		context.Background(),
		db.GetTraefikInstanceByNameParams{
			ProfileID: profileID,
			Name:      instanceName,
		},
	)
	if err != nil {
		// Create new temporary instance
		instance = db.TraefikInstance{
			ProfileID: profileID,
			Name:      instanceName,
			Url:       instanceURL,
		}
	}

	endpoints := []struct {
		name string
		path string
	}{
		{"raw", RawAPI},
		{"entrypoints", EntrypointsAPI},
		{"overview", OverviewAPI},
		{"version", VersionAPI},
	}

	responses := make(chan APIResponse, len(endpoints))

	// Fetch all endpoints concurrently
	for _, endpoint := range endpoints {
		go func(name, path string) {
			data, err := fetch(instance, path)
			responses <- APIResponse{Name: name, Data: data, Err: err}
		}(endpoint.name, endpoint.path)
	}

	// Collect responses
	results := make(map[string]io.ReadCloser)
	var fetchErrors []error

	for range len(endpoints) {
		resp := <-responses
		if resp.Err != nil {
			slog.Error("failed to fetch traefik data", "endpoint", resp.Name, "error", resp.Err)
			fetchErrors = append(fetchErrors, resp.Err)
			continue
		}
		results[resp.Name] = resp.Data
	}

	// Clean up all response bodies
	defer func() {
		for _, body := range results {
			if body != nil {
				if err := body.Close(); err != nil {
					slog.Error("failed to close response body", "error", err)
				}
			}
		}
	}()

	// If any critical fetch failed, abort
	if len(fetchErrors) > 0 {
		return
	}

	// Decode responses
	var config schema.Configuration
	if err := json.NewDecoder(results["raw"]).Decode(&config); err != nil {
		slog.Error("failed to decode raw data", "error", err)
		return
	}

	var entrypoints schema.EntryPoints
	if err := json.NewDecoder(results["entrypoints"]).Decode(&entrypoints); err != nil {
		slog.Error("failed to decode entrypoints", "error", err)
		return
	}

	var overview schema.Overview
	if err := json.NewDecoder(results["overview"]).Decode(&overview); err != nil {
		slog.Error("failed to decode overview", "error", err)
		return
	}

	var version schema.Version
	if err := json.NewDecoder(results["version"]).Decode(&version); err != nil {
		slog.Error("failed to decode version", "error", err)
		return
	}

	// Upsert with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := q.UpsertTraefikInstance(ctx, db.UpsertTraefikInstanceParams{
		Name:        instance.Name,
		Url:         instance.Url,
		Username:    instance.Username,
		Password:    instance.Password,
		Tls:         instance.Tls,
		Entrypoints: &entrypoints,
		Overview:    &overview,
		Config:      &config,
		Version:     &version,
		ProfileID:   profileID,
	}); err != nil {
		slog.Error("failed to update traefik instance", "error", err)
		return
	}
}

func fetch(instance db.TraefikInstance, endpoint string) (io.ReadCloser, error) {
	if instance.Url == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	client := &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !instance.Tls,
			},
		},
	}

	req, err := http.NewRequest("GET", instance.Url+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if instance.Username != nil && instance.Password != nil {
		req.SetBasicAuth(*instance.Username, *instance.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
		return nil, fmt.Errorf("failed to fetch %s: %w", instance.Url+endpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
