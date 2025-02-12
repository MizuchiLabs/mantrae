package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
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
		mergedConfig := &db.TraefikConfiguration{
			Routers:        make(map[string]*runtime.RouterInfo),
			Middlewares:    make(map[string]*runtime.MiddlewareInfo),
			Services:       make(map[string]*db.ServiceInfo),
			TCPRouters:     make(map[string]*runtime.TCPRouterInfo),
			TCPMiddlewares: make(map[string]*runtime.TCPMiddlewareInfo),
			TCPServices:    make(map[string]*runtime.TCPServiceInfo),
			UDPRouters:     make(map[string]*runtime.UDPRouterInfo),
			UDPServices:    make(map[string]*runtime.UDPServiceInfo),
		}

		local, err := q.GetLocalTraefikConfig(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if local.Config == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		agents, err := q.GetAgentTraefikConfigs(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Merge configurations (prefer local)
		for _, agent := range agents {
			mergedConfig = traefik.MergeConfigs(mergedConfig, agent.Config)
		}
		mergedConfig = traefik.MergeConfigs(mergedConfig, local.Config)

		// Convert to dynamic
		dynamic := traefik.ConvertToDynamicConfig(mergedConfig)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dynamic); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
