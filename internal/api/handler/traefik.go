package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/meta"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"golang.org/x/sync/singleflight"
	"gopkg.in/yaml.v3"
)

var (
	updateGroup    singleflight.Group
	lastUpdateTime sync.Map
)

func scheduleUpdate(r *http.Request, app *config.App, profileID int64) {
	instanceName := r.Header.Get(meta.HeaderTraefikName)
	instanceURL := r.Header.Get(meta.HeaderTraefikURL)
	if instanceName == "" || instanceURL == "" {
		slog.Debug("Skipping traefik update, missing headers")
		return
	}

	// Check if we recently updated this instance
	if lastUpdate, ok := lastUpdateTime.Load(instanceName); ok {
		if time.Since(lastUpdate.(time.Time)) < 5*time.Second {
			return // Skip if updated recently
		}
	}

	// Use singleflight to ensure only one update per instance
	go func() {
		_, _, _ = updateGroup.Do(instanceName, func() (any, error) {
			lastUpdateTime.Store(instanceName, time.Now())
			result, err := traefik.UpdateTraefikInstance(r, app.Conn.Query, profileID)
			if err != nil {
				slog.Error("failed to update traefik instance", "error", err)
				return nil, nil
			}

			if result != nil {
				app.Event.Broadcast(&mantraev1.EventStreamResponse{
					Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
					Data: &mantraev1.EventStreamResponse_TraefikInstance{
						TraefikInstance: result.ToProto(),
					},
				})
			}
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

		profile, err := a.Conn.Query.GetProfileByName(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if urlToken != profile.Token && headerToken != profile.Token {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Create or update traefik instance
		scheduleUpdate(r, a, profile.ID)

		cfg, err := traefik.BuildDynamicConfig(r.Context(), a.Conn.Query, *profile)
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
