package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

func GetTraefikConfig(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
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
		json.NewEncoder(w).Encode(config)
	}
}

func PublishTraefikConfig(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
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

		localT, err := q.GetLocalTraefikConfig(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		agentT, err := q.GetAgentTraefikConfigs(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Merge configurations from each agent
		for _, a := range agentT {
			if a.Config == nil {
				continue
			}
			if a.Config.Routers != nil {
				for k, v := range a.Config.Routers {
					mergedConfig.Routers[k] = v
				}
			}
			if a.Config.Services != nil {
				for k, v := range a.Config.Services {
					mergedConfig.Services[k] = v
				}
			}
		}

		if localT.Config == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Overlay local config to ensure it takes precedence
		if localT.Config.Routers != nil {
			for k, v := range localT.Config.Routers {
				mergedConfig.Routers[k] = v
			}
		}
		if localT.Config.Middlewares != nil {
			for k, v := range localT.Config.Middlewares {
				mergedConfig.Middlewares[k] = v
			}
		}
		if localT.Config.Services != nil {
			for k, v := range localT.Config.Services {
				mergedConfig.Services[k] = v
			}
		}
		if localT.Config.TCPRouters != nil {
			for k, v := range localT.Config.TCPRouters {
				mergedConfig.TCPRouters[k] = v
			}
		}
		if localT.Config.TCPMiddlewares != nil {
			for k, v := range localT.Config.TCPMiddlewares {
				mergedConfig.TCPMiddlewares[k] = v
			}
		}
		if localT.Config.TCPServices != nil {
			for k, v := range localT.Config.TCPServices {
				mergedConfig.TCPServices[k] = v
			}
		}
		if localT.Config.UDPRouters != nil {
			for k, v := range localT.Config.UDPRouters {
				mergedConfig.UDPRouters[k] = v
			}
		}
		if localT.Config.UDPServices != nil {
			for k, v := range localT.Config.UDPServices {
				mergedConfig.UDPServices[k] = v
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mergedConfig)
	}
}
