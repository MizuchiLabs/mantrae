package traefik

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
	"golang.org/x/exp/maps"
)

func DecodeAgentConfig(DB *sql.DB, agent db.Agent) error {
	// Initialize merged configuration
	params := db.UpsertTraefikAgentConfigParams{
		ProfileID: agent.ProfileID,
		AgentID:   &agent.ID,
		Config: &db.TraefikConfiguration{
			Routers:        make(map[string]*runtime.RouterInfo),
			Middlewares:    make(map[string]*runtime.MiddlewareInfo),
			Services:       make(map[string]*db.ServiceInfo),
			TCPRouters:     make(map[string]*runtime.TCPRouterInfo),
			TCPMiddlewares: make(map[string]*runtime.TCPMiddlewareInfo),
			TCPServices:    make(map[string]*runtime.TCPServiceInfo),
			UDPRouters:     make(map[string]*runtime.UDPRouterInfo),
			UDPServices:    make(map[string]*runtime.UDPServiceInfo),
		},
	}

	for _, container := range *agent.Containers {
		dynConfig := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}

		err := parser.Decode(
			container.Labels,
			dynConfig,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		)
		if err != nil {
			return err
		}

		runtimeConfig := runtime.NewConfig(*dynConfig)

		// Merge routers
		maps.Copy(params.Config.Routers, runtimeConfig.Routers)
		maps.Copy(params.Config.TCPRouters, runtimeConfig.TCPRouters)
		maps.Copy(params.Config.UDPRouters, runtimeConfig.UDPRouters)

		// Merge services
		for k, v := range runtimeConfig.Services {
			service := v.Service
			for i, svc := range service.LoadBalancer.Servers {
				url := fmt.Sprintf("%s://%s:%s", svc.Scheme, *agent.ActiveIp, svc.Port)
				service.LoadBalancer.Servers[i].URL = url
			}
			params.Config.Services[k] = &db.ServiceInfo{
				ServiceInfo: &runtime.ServiceInfo{Service: service},
			}
		}

		for k, v := range runtimeConfig.TCPServices {
			service := v.TCPService
			for i, svc := range service.LoadBalancer.Servers {
				address := fmt.Sprintf("%s:%s", *agent.ActiveIp, svc.Port)
				service.LoadBalancer.Servers[i].Address = address
			}
			params.Config.TCPServices[k] = &runtime.TCPServiceInfo{TCPService: service}
		}

		for k, v := range runtimeConfig.UDPServices {
			service := v.UDPService
			for i, svc := range service.LoadBalancer.Servers {
				address := fmt.Sprintf("%s:%s", *agent.ActiveIp, svc.Port)
				service.LoadBalancer.Servers[i].Address = address
			}
			params.Config.UDPServices[k] = &runtime.UDPServiceInfo{UDPService: service}
		}

		// Merge middlewares
		maps.Copy(params.Config.Middlewares, runtimeConfig.Middlewares)
		maps.Copy(params.Config.TCPMiddlewares, runtimeConfig.TCPMiddlewares)
	}

	q := db.New(DB)
	if err := q.UpsertTraefikAgentConfig(context.Background(), params); err != nil {
		return err
	}

	return nil
}
