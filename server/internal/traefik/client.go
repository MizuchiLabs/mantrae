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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
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

func UpdateTraefikInstance(
	r *http.Request,
	q *db.Queries,
	profileID int64,
) (*db.TraefikInstance, error) {
	instanceName := r.Header.Get(meta.HeaderTraefikName)
	instanceURL := r.Header.Get(meta.HeaderTraefikURL)
	if instanceName == "" || instanceURL == "" {
		slog.Debug("Skipping traefik update, missing headers")
		return nil, nil
	}

	params := db.UpsertTraefikInstanceParams{
		ID:        uuid.New().String(),
		ProfileID: profileID,
		Name:      instanceName,
		Url:       instanceURL,
		Tls:       strings.HasPrefix(instanceURL, "https"),
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
			data, err := fetch(params, path)
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
		return nil, nil
	}

	// Decode responses
	if results["raw"] != nil {
		var config schema.Configuration
		if err := json.NewDecoder(results["raw"]).Decode(&config); err != nil {
			slog.Warn("failed to decode raw config", "error", err)
		} else {
			params.Config = &config
		}
	}

	if results["entrypoints"] != nil {
		var entrypoints schema.EntryPoints
		if err := json.NewDecoder(results["entrypoints"]).Decode(&entrypoints); err != nil {
			slog.Warn("failed to decode entrypoints", "error", err)
		} else {
			params.Entrypoints = &entrypoints
		}
	}

	if results["overview"] != nil {
		var overview schema.Overview
		params.Overview = &overview
		if err := json.NewDecoder(results["overview"]).Decode(&overview); err != nil {
			slog.Warn("failed to decode overview", "error", err)
		} else {
			params.Overview = &overview
		}
	}

	if results["version"] != nil {
		var version schema.Version
		if err := json.NewDecoder(results["version"]).Decode(&version); err != nil {
			slog.Warn("failed to decode version", "error", err)
		} else {
			params.Version = &version
		}
	}

	result, err := q.UpsertTraefikInstance(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func fetch(instance db.UpsertTraefikInstanceParams, endpoint string) (io.ReadCloser, error) {
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
