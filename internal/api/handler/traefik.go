package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"gopkg.in/yaml.v3"
)

func PublishTraefikConfig(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, err := traefik.BuildDynamicConfig(r.Context(), a.Conn.GetQuery(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		format := r.URL.Query().Get("format")
		accept := r.Header.Get("Accept")

		// Determine response format: prefer query param over header
		if format == "yaml" || (format == "" && strings.Contains(accept, "yaml")) {
			w.Header().Set("Content-Type", "application/x-yaml")
			if err := yaml.NewEncoder(w).Encode(cfg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		// Default to JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
