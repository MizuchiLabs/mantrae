package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"gopkg.in/yaml.v3"
)

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
