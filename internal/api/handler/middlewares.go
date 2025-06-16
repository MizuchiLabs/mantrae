package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/source"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

type UpsertMiddlewareParams struct {
	Name          string                     `json:"name"`
	Protocol      string                     `json:"protocol"`
	Middleware    *runtime.MiddlewareInfo    `json:"middleware"`
	TCPMiddleware *runtime.TCPMiddlewareInfo `json:"tcpMiddleware"`
}

type DeleteMiddlewareParams struct {
	ProfileID int64  `json:"profileId"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
}

type BulkDeleteMiddlewareParams struct {
	ProfileID int64 `json:"profileId"`
	Items     []struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
	} `json:"items"`
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
func UpsertMiddleware(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var params UpsertMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get existing config
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), profileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Initialize maps if nil
		if existingConfig.Config == nil {
			existingConfig.Config = &db.TraefikConfiguration{}
		}
		if existingConfig.Config.Middlewares == nil {
			existingConfig.Config.Middlewares = make(map[string]*runtime.MiddlewareInfo)
		}
		if existingConfig.Config.TCPMiddlewares == nil {
			existingConfig.Config.TCPMiddlewares = make(map[string]*runtime.TCPMiddlewareInfo)
		}

		// Ensure name has no @
		params.Name = strings.Split(params.Name, "@")[0]

		if params.Middleware == nil {
			http.Error(w, "Missing middleware configuration", http.StatusBadRequest)
			return
		}

		// Update configuration based on type
		switch params.Protocol {
		case "http":
			if err = db.VerifyMiddleware(params.Middleware.Middleware); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			existingConfig.Config.Middlewares[params.Name] = params.Middleware
		case "tcp":
			existingConfig.Config.TCPMiddlewares[params.Name] = params.TCPMiddleware
		default:
			http.Error(w, "invalid middleware type: must be http or tcp", http.StatusBadRequest)
			return
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryTraefik,
		}

		// Return the updated configuration
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(existingConfig.Config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteMiddleware deletes a middleware
func DeleteMiddleware(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params DeleteMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.ProfileID == 0 || params.Name == "" || params.Protocol == "" {
			http.Error(w, "Missing middleware name or protocol", http.StatusBadRequest)
			return
		}

		q := a.Conn.GetQuery()
		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), params.ProfileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		switch params.Protocol {
		case "http":
			delete(existingConfig.Config.Middlewares, params.Name)
		case "tcp":
			delete(existingConfig.Config.TCPMiddlewares, params.Name)
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryTraefik,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func BulkDeleteMiddleware(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params BulkDeleteMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.ProfileID == 0 {
			http.Error(w, "Missing profile ID", http.StatusBadRequest)
			return
		}

		q := a.Conn.GetQuery()
		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), params.ProfileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		for _, item := range params.Items {
			if item.Name == "" || item.Protocol == "" {
				continue // Skip invalid entries
			}

			switch item.Protocol {
			case "http":
				delete(existingConfig.Config.Middlewares, item.Name)
			case "tcp":
				delete(existingConfig.Config.TCPMiddlewares, item.Name)
			default:
				continue
			}
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryTraefik,
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
	if err := json.NewEncoder(w).Encode(middlewarePlugins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
