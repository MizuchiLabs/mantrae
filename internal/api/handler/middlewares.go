package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

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

func GetHTTPMiddlewaresBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.GetHTTPMiddlewaresBySourceParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		middlewares, err := q.GetHTTPMiddlewaresBySource(r.Context(), middleware)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middlewares)
	}
}

func GetTCPMiddlewaresBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.GetTCPMiddlewaresBySourceParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		middlewares, err := q.GetTCPMiddlewaresBySource(r.Context(), middleware)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middlewares)
	}
}

func GetHTTPMiddlewareByName(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.GetHTTPMiddlewareByNameParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		middlewares, err := q.GetHTTPMiddlewareByName(r.Context(), middleware)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middlewares)
	}
}

func GetTCPMiddlewareByName(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.GetTCPMiddlewareByNameParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		middlewares, err := q.GetTCPMiddlewareByName(r.Context(), middleware)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middlewares)
	}
}

func UpsertHTTPMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.UpsertHTTPMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.UpsertHTTPMiddleware(r.Context(), middleware); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpsertTCPMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.UpsertTCPMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.UpsertTCPMiddleware(r.Context(), middleware); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteHTTPMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.DeleteHTTPMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.DeleteHTTPMiddleware(r.Context(), middleware); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteTCPMiddleware(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var middleware db.DeleteTCPMiddlewareParams
		if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.DeleteTCPMiddleware(r.Context(), middleware); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
