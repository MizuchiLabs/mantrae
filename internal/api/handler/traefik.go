package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func PublishTraefikConfig(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		profile, err := q.GetProfileByName(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cfg := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{
				Routers:     make(map[string]*dynamic.Router),
				Middlewares: make(map[string]*dynamic.Middleware),
				Services:    make(map[string]*dynamic.Service),
			},
			TCP: &dynamic.TCPConfiguration{
				Routers:     make(map[string]*dynamic.TCPRouter),
				Middlewares: make(map[string]*dynamic.TCPMiddleware),
				Services:    make(map[string]*dynamic.TCPService),
			},
			UDP: &dynamic.UDPConfiguration{
				Routers:  make(map[string]*dynamic.UDPRouter),
				Services: make(map[string]*dynamic.UDPService),
			},
		}

		// Routers
		httpRouters, err := q.GetHttpRoutersByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tcpRouters, err := q.GetTcpRoutersByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		udpRouters, err := q.GetUdpRoutersByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Services
		httpServices, err := q.GetHttpServicesByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tcpServices, err := q.GetTcpServicesByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		udpServices, err := q.GetUdpServicesByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Middlewares
		httpMiddlewares, err := q.GetHttpMiddlewaresByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tcpMiddlewares, err := q.GetTcpMiddlewaresByProfile(r.Context(), profile.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, r := range httpRouters {
			cfg.HTTP.Routers[r.Name] = r.Config
		}
		for _, r := range tcpRouters {
			cfg.TCP.Routers[r.Name] = r.Config
		}
		for _, r := range udpRouters {
			cfg.UDP.Routers[r.Name] = r.Config
		}

		for _, s := range httpServices {
			cfg.HTTP.Services[s.Name] = s.Config
		}
		for _, s := range tcpServices {
			cfg.TCP.Services[s.Name] = s.Config
		}
		for _, s := range udpServices {
			cfg.UDP.Services[s.Name] = s.Config
		}

		for _, m := range httpMiddlewares {
			cfg.HTTP.Middlewares[m.Name] = m.Config
		}
		for _, m := range tcpMiddlewares {
			cfg.TCP.Middlewares[m.Name] = m.Config
		}

		// Cleanup empty sections (to avoid Traefik {} block warnings)
		if len(cfg.HTTP.Routers) == 0 && len(cfg.HTTP.Middlewares) == 0 &&
			len(cfg.HTTP.Services) == 0 {
			cfg.HTTP = nil
		}
		if len(cfg.TCP.Routers) == 0 && len(cfg.TCP.Middlewares) == 0 &&
			len(cfg.TCP.Services) == 0 {
			cfg.TCP = nil
		}
		if len(cfg.UDP.Routers) == 0 && len(cfg.UDP.Services) == 0 {
			cfg.UDP = nil
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
