package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/traefik"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"golang.org/x/sync/singleflight"
	"gopkg.in/yaml.v3"
)

var (
	updateGroup    singleflight.Group
	lastUpdateTime sync.Map // map[string]time.Time
)

func scheduleUpdate(r *http.Request, q *db.Queries, profileID int64) {
	instanceName := r.Header.Get(meta.HeaderTraefikName)
	if instanceName == "" {
		return
	}

	// Check if we recently updated this instance
	if lastUpdate, ok := lastUpdateTime.Load(instanceName); ok {
		if time.Since(lastUpdate.(time.Time)) < 30*time.Second {
			return // Skip if updated recently
		}
	}

	// Use singleflight to ensure only one update per instance
	go func() {
		_, _, _ = updateGroup.Do(instanceName, func() (any, error) {
			lastUpdateTime.Store(instanceName, time.Now())
			traefik.UpdateTraefikInstance(r, q, profileID)
			return nil, nil
		})
	}()
}

func PublishTraefikConfig(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		urlToken := r.URL.Query().Get("token")
		headerToken := r.Header.Get(meta.HeaderTraefikToken)
		accept := r.Header.Get("Accept")

		profile, err := a.Conn.GetQuery().GetProfileByName(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if urlToken != profile.Token && headerToken != profile.Token {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Create or update traefik instance
		scheduleUpdate(r, a.Conn.GetQuery(), profile.ID)

		cfg, err := traefik.BuildDynamicConfig(r.Context(), a.Conn.GetQuery(), profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Determine response format: prefer query param over header
		if format == "yaml" || (format == "" && strings.Contains(accept, "yaml")) {
			w.Header().Set("Content-Type", "application/x-yaml")
			enc := yaml.NewEncoder(w)
			enc.SetIndent(2)
			if err := enc.Encode(cfg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// Default to JSON
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ") // Indent the output with two spaces
		if err := enc.Encode(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
