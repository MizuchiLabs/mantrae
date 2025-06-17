package traefik

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
)

const (
	RawAPI         = "/api/rawdata"
	EntrypointsAPI = "/api/entrypoints"
	OverviewAPI    = "/api/overview"
	VersionAPI     = "/api/version"
)

func UpdateTraefikAPI(DB *sql.DB, instanceID int64) error {
	q := db.New(DB)
	instance, err := q.GetTraefikInstance(context.Background(), instanceID)
	if err != nil {
		return fmt.Errorf("failed to get traefik instance: %w", err)
	}

	rawResp, err := fetch(instance, RawAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", instance.Url+RawAPI, err)
	}
	defer rawResp.Close()

	var config schema.Configuration
	if err = json.NewDecoder(rawResp).Decode(&config); err != nil {
		return fmt.Errorf("failed to decode raw data: %w", err)
	}

	entrypointsResp, err := fetch(instance, EntrypointsAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", instance.Url+EntrypointsAPI, err)
	}
	defer entrypointsResp.Close()

	var entrypoints schema.EntryPoints
	if err = json.NewDecoder(entrypointsResp).Decode(&entrypoints); err != nil {
		return fmt.Errorf("failed to decode entrypoints: %w", err)
	}

	overviewResp, err := fetch(instance, OverviewAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", instance.Url+OverviewAPI, err)
	}
	defer overviewResp.Close()

	var overview schema.Overview
	if err = json.NewDecoder(overviewResp).Decode(&overview); err != nil {
		return fmt.Errorf("failed to decode overview: %w", err)
	}

	versionResp, err := fetch(instance, VersionAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", instance.Url+VersionAPI, err)
	}
	defer versionResp.Close()

	var version schema.Version
	if err := json.NewDecoder(versionResp).Decode(&version); err != nil {
		return fmt.Errorf("failed to decode version: %w", err)
	}

	if _, err := q.UpdateTraefikInstance(context.Background(), db.UpdateTraefikInstanceParams{
		ID:          instance.ID,
		Url:         instance.Url,
		Username:    instance.Username,
		Password:    instance.Password,
		Tls:         instance.Tls,
		Entrypoints: &entrypoints,
		Overview:    &overview,
		Version:     &version.Version,
		Config:      &config,
	}); err != nil {
		return fmt.Errorf("failed to update api data: %w", err)
	}

	return nil
}

func fetch(instance db.TraefikInstance, endpoint string) (io.ReadCloser, error) {
	if instance.Url == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	client := http.Client{Timeout: time.Second * 10}
	if !instance.Tls {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
