package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type UpsertMiddlewareParams struct {
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Middleware    *dynamic.Middleware    `json:"middleware"`
	TCPMiddleware *dynamic.TCPMiddleware `json:"tcpMiddleware"`
}

type Plugin struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	Author        string        `json:"author"`
	Type          string        `json:"type"`
	Import        string        `json:"import"`
	Summary       string        `json:"summary"`
	IconUrl       string        `json:"iconUrl"`
	BannerUrl     string        `json:"bannerUrl"`
	Readme        string        `json:"readme"`
	LatestVersion string        `json:"latestVersion"`
	Versions      []string      `json:"versions"`
	Stars         int           `json:"stars"`
	Snippet       PluginSnippet `json:"snippet"`
	CreatedAt     string        `json:"createdAt"`
}

type PluginSnippet struct {
	Yaml string `json:"yaml"`
}

// UpsertMiddleware
func UpsertMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params UpsertMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// if err := validateRouterParams(&params); err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// Get existing config
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		existingConfig, err := q.GetTraefikConfigBySource(r.Context(), db.GetTraefikConfigBySourceParams{
			ProfileID: profileID,
			Source:    source.Local,
		})
		if err != nil {
			http.Error(w, "Failed to get existing config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Initialize maps if nil
		if existingConfig.Config.HTTP.Middlewares == nil {
			existingConfig.Config.HTTP.Middlewares = make(map[string]*dynamic.Middleware)
		}
		if existingConfig.Config.TCP.Middlewares == nil {
			existingConfig.Config.TCP.Middlewares = make(map[string]*dynamic.TCPMiddleware)
		}

		// Ensure name has @http suffix
		if !strings.HasSuffix(params.Name, "@http") {
			params.Name = fmt.Sprintf("%s@http", strings.Split(params.Name, "@")[0])
		}

		// Update configuration based on type
		switch params.Type {
		case "http":
			existingConfig.Config.HTTP.Middlewares[params.Name] = params.Middleware
		case "tcp":
			existingConfig.Config.TCP.Middlewares[params.Name] = params.TCPMiddleware
		default:
			http.Error(w, "invalid middleware type: must be http or tcp", http.StatusBadRequest)
			return
		}

		err = q.UpdateTraefikConfig(r.Context(), db.UpdateTraefikConfigParams{
			ID:     existingConfig.ID,
			Source: source.Local,
			Config: existingConfig.Config,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the updated configuration
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingConfig.Config)
	}
}

// DeleteMiddleware
func DeleteMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		mwName := r.PathValue("name")
		mwType := r.PathValue("type")

		if mwName == "" || mwType == "" {
			http.Error(w, "Missing middleware name or type", http.StatusBadRequest)
			return
		}

		// Ensure name has @http suffix for consistency
		if !strings.HasSuffix(mwName, "@http") {
			mwName = fmt.Sprintf("%s@http", strings.Split(mwName, "@")[0])
		}

		existingConfig, err := q.GetTraefikConfigBySource(r.Context(), db.GetTraefikConfigBySourceParams{
			ProfileID: profileID,
			Source:    source.Local,
		})
		if err != nil {
			http.Error(w, "Failed to get existing config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove router and service based on type
		switch mwType {
		case "http":
			delete(existingConfig.Config.HTTP.Middlewares, mwName)
		case "tcp":
			delete(existingConfig.Config.TCP.Middlewares, mwName)
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpdateTraefikConfig(r.Context(), db.UpdateTraefikConfigParams{
			ID:     existingConfig.ID,
			Source: source.Local,
			Config: existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetMiddlewarePlugins retrieves middleware plugins available for Traefik.
func GetMiddlewarePlugins(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://plugins.traefik.io/api/services/plugins")
	if err != nil {
		http.Error(w, "Failed to fetch plugins", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the response body into the Go Plugin struct
	var plugins []Plugin
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to parse plugins: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Filter out non-middleware plugins
	var middlewarePlugins []Plugin
	for _, plugin := range plugins {
		if plugin.Type == "middleware" {
			middlewarePlugins = append(middlewarePlugins, plugin)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(middlewarePlugins)
}
