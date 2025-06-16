package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/source"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

func GetTraefikConfig(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		profile_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		src := source.Source(r.PathValue("source"))
		if !src.Valid() {
			http.Error(w, "invalid source", http.StatusBadRequest)
			return
		}
		config, err := q.GetTraefikConfigBySource(r.Context(), db.GetTraefikConfigBySourceParams{
			ProfileID: profile_id,
			Source:    src,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PublishTraefikConfig(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		profile, err := q.GetProfileByName(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Initialize merged config
		merged := &db.TraefikConfiguration{
			Routers:        make(map[string]*runtime.RouterInfo),
			Middlewares:    make(map[string]*runtime.MiddlewareInfo),
			Services:       make(map[string]*db.ServiceInfo),
			TCPRouters:     make(map[string]*runtime.TCPRouterInfo),
			TCPMiddlewares: make(map[string]*runtime.TCPMiddlewareInfo),
			TCPServices:    make(map[string]*runtime.TCPServiceInfo),
			UDPRouters:     make(map[string]*runtime.UDPRouterInfo),
			UDPServices:    make(map[string]*runtime.UDPServiceInfo),
		}

		// Get local config
		local, err := q.GetLocalTraefikConfig(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		agents, err := q.GetAgentTraefikConfigs(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Merge agent configurations (agent config)
		for _, agent := range agents {
			merged = traefik.MergeConfigs(merged, agent.Config)
		}
		if local.Config == nil && merged == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Merge with local config
		merged = traefik.MergeConfigs(merged, local.Config)

		// Convert to dynamic
		dynamic := traefik.ConvertToDynamicConfig(merged)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dynamic); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
