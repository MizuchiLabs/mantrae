package traefik

import (
	"context"
	"database/sql"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
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

		// Merge configurations
		for k, v := range runtimeConfig.Routers {
			params.Config.Routers[k] = v

			service := dynamic.Service{
				LoadBalancer: &dynamic.ServersLoadBalancer{
					Servers: []dynamic.Server{{URL: *agent.ActiveIp}},
				},
			}
			params.Config.Services[k] = &db.ServiceInfo{
				ServiceInfo:  &runtime.ServiceInfo{Service: &service},
				ServerStatus: map[string]string{},
			}
		}
		for k, v := range runtimeConfig.Middlewares {
			params.Config.Middlewares[k] = v
		}
		for k, v := range runtimeConfig.TCPRouters {
			params.Config.TCPRouters[k] = v

			service := dynamic.TCPService{
				LoadBalancer: &dynamic.TCPServersLoadBalancer{
					Servers: []dynamic.TCPServer{{Address: *agent.ActiveIp}},
				},
			}
			params.Config.TCPServices[k] = &runtime.TCPServiceInfo{TCPService: &service}
		}
		for k, v := range runtimeConfig.TCPMiddlewares {
			params.Config.TCPMiddlewares[k] = v
		}
		for k, v := range runtimeConfig.UDPRouters {
			params.Config.UDPRouters[k] = v

			service := dynamic.UDPService{
				LoadBalancer: &dynamic.UDPServersLoadBalancer{
					Servers: []dynamic.UDPServer{{Address: *agent.ActiveIp}},
				},
			}
			params.Config.UDPServices[k] = &runtime.UDPServiceInfo{UDPService: &service}
		}
	}

	q := db.New(DB)
	if err := q.UpsertTraefikAgentConfig(context.Background(), params); err != nil {
		return err
	}

	return nil
}
