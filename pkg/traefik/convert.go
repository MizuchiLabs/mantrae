package traefik

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// DecodeFromLabels uses the traefik parses to decode the config from the labels into our Dynamic struct
func DecodeFromLabels(id string, container []byte) error {
	agent, err := db.Query.GetAgentByID(context.Background(), id)
	if err != nil {
		return err
	}

	var containers []*agentv1.Container
	if err := json.Unmarshal(container, &containers); err != nil {
		return err
	}

	for i := range containers {
		// Convert labels to official traefik types
		config := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}
		if err := parser.Decode(
			containers[i].Labels,
			config,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		); err != nil {
			return err
		}

		// Container portmap
		portmap := containers[i].Portmap
		var useIP string
		if agent.PublicIp != nil && agent.ActiveIp == nil {
			useIP = *agent.PublicIp
		}
		if agent.ActiveIp != nil {
			useIP = *agent.ActiveIp
		}

		// Add to our database
		if config.HTTP != nil {
			for i, router := range config.HTTP.Routers {
				priority := int64(router.Priority)
				dbRouter := db.Router{
					ProfileID:   agent.ProfileID,
					Name:        i,
					Provider:    "http",
					Protocol:    "http",
					AgentID:     &agent.ID,
					EntryPoints: router.EntryPoints,
					Middlewares: router.Middlewares,
					Rule:        router.Rule,
					RuleSyntax:  &router.RuleSyntax,
					Service:     router.Service,
					Priority:    &priority,
					Tls:         router.TLS,
				}
				if err := dbRouter.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams(dbRouter)); err != nil {
					return err
				}
			}

			for i, service := range config.HTTP.Services {
				servers := []dynamic.Server{}

				// Get created router
				router, err := db.Query.GetRouterByName(
					context.Background(),
					db.GetRouterByNameParams{
						ProfileID: agent.ProfileID,
						Name:      i,
					},
				)
				if err != nil {
					slog.Error("Failed to get router", "error", err)
					continue
				}
				if err := router.DecodeFields(); err != nil {
					return err
				}

				// Build servers
				for _, server := range service.LoadBalancer.Servers {
					newServer := dynamic.Server{}

					internalPort, err := strconv.ParseInt(server.Port, 10, 32)
					if err != nil {
						slog.Error("Failed to parse internal port", "port", server.Port)
						continue
					}
					if externalPort, ok := portmap[int32(internalPort)]; ok {
						if externalPort == 443 {
							newServer.URL = fmt.Sprintf("https://%s:%d", useIP, externalPort)
						} else {
							newServer.URL = fmt.Sprintf("http://%s:%d", useIP, externalPort)
						}
						router.UpdateError("agent", "")
					} else {
						router.UpdateError("agent", "No external port mapping found for internal port "+server.Port)
						continue
					}
					servers = append(servers, newServer)
				}
				if len(servers) == 0 {
					// Fallback
					newServer := dynamic.Server{
						URL: fmt.Sprintf("http://%s:%d", useIP, 80),
					}
					servers = append(servers, newServer)
				}
				loadBalancer := dynamic.ServersLoadBalancer{
					Sticky:             service.LoadBalancer.Sticky,
					Servers:            servers,
					HealthCheck:        service.LoadBalancer.HealthCheck,
					PassHostHeader:     service.LoadBalancer.PassHostHeader,
					ResponseForwarding: service.LoadBalancer.ResponseForwarding,
					ServersTransport:   service.LoadBalancer.ServersTransport,
				}
				dbService := db.UpsertServiceParams{
					ProfileID:    agent.ProfileID,
					Name:         i,
					Provider:     "http",
					Protocol:     "http",
					AgentID:      &agent.ID,
					LoadBalancer: &loadBalancer,
					Weighted:     service.Weighted,
					Mirroring:    service.Mirroring,
					Failover:     service.Failover,
				}
				if err := dbService.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertService(context.Background(), dbService); err != nil {
					return err
				}
			}
		}

		if config.TCP != nil {
			for i, router := range config.TCP.Routers {
				priority := int64(router.Priority)
				dbRouter := db.Router{
					ProfileID:   agent.ProfileID,
					Name:        i,
					Provider:    "http",
					Protocol:    "tcp",
					AgentID:     &agent.ID,
					EntryPoints: router.EntryPoints,
					Middlewares: router.Middlewares,
					Rule:        router.Rule,
					RuleSyntax:  &router.RuleSyntax,
					Service:     router.Service,
					Priority:    &priority,
					Tls:         router.TLS,
				}
				if err := dbRouter.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams(dbRouter)); err != nil {
					return err
				}
			}

			for i, service := range config.TCP.Services {
				servers := []dynamic.TCPServer{}

				// Build servers
				for _, server := range service.LoadBalancer.Servers {
					newServer := dynamic.TCPServer{}

					internalPort, err := strconv.ParseInt(server.Port, 10, 32)
					if err != nil {
						slog.Error("Failed to parse internal port", "port", server.Port)
						continue
					}
					if externalPort, ok := portmap[int32(internalPort)]; ok {
						newServer.Address = fmt.Sprintf("%s:%d", useIP, externalPort)
					} else {
						slog.Error("No external port mapping found for internal port", "port", server.Port)
						continue
					}
					servers = append(servers, newServer)
				}
				if len(servers) == 0 {
					slog.Error("No servers found for service", "service", service)
					// Fallback
					newServer := dynamic.TCPServer{
						Address: fmt.Sprintf("%s:%d", useIP, 80),
					}
					servers = append(servers, newServer)
				}

				loadBalancer := dynamic.TCPServersLoadBalancer{
					ProxyProtocol:    service.LoadBalancer.ProxyProtocol,
					Servers:          servers,
					ServersTransport: service.LoadBalancer.ServersTransport,
					TerminationDelay: service.LoadBalancer.TerminationDelay,
				}
				dbService := db.UpsertServiceParams{
					ProfileID:    agent.ProfileID,
					Name:         i,
					Provider:     "http",
					Protocol:     "tcp",
					AgentID:      &agent.ID,
					LoadBalancer: &loadBalancer,
					Weighted:     service.Weighted,
				}
				if err := dbService.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertService(context.Background(), dbService); err != nil {
					return err
				}
			}
		}

		if config.UDP != nil {
			for i, router := range config.UDP.Routers {
				dbRouter := db.Router{
					ProfileID:   agent.ProfileID,
					Name:        i,
					Provider:    "http",
					Protocol:    "udp",
					AgentID:     &agent.ID,
					EntryPoints: router.EntryPoints,
					Service:     router.Service,
				}
				if err := dbRouter.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams(dbRouter)); err != nil {
					return err
				}
			}

			for i, service := range config.UDP.Services {
				servers := []dynamic.UDPServer{}

				// Build servers
				for _, server := range service.LoadBalancer.Servers {
					newServer := dynamic.UDPServer{}

					internalPort, err := strconv.ParseInt(server.Port, 10, 32)
					if err != nil {
						slog.Error("Failed to parse internal port", "port", server.Port)
						continue
					}
					if externalPort, ok := portmap[int32(internalPort)]; ok {
						newServer.Address = fmt.Sprintf("%s:%d", useIP, externalPort)
					} else {
						slog.Error("No external port mapping found for internal port", "port", server.Port)
						continue
					}
					servers = append(servers, newServer)
				}
				if len(servers) == 0 {
					slog.Error("No servers found for service", "service", service)
					// Fallback
					newServer := dynamic.UDPServer{
						Address: fmt.Sprintf("%s:%d", useIP, 80),
					}
					servers = append(servers, newServer)
				}

				loadBalancer := dynamic.UDPServersLoadBalancer{Servers: servers}
				dbService := db.UpsertServiceParams{
					ProfileID:    agent.ProfileID,
					Name:         i,
					Provider:     "http",
					Protocol:     "udp",
					AgentID:      &agent.ID,
					LoadBalancer: &loadBalancer,
					Weighted:     service.Weighted,
				}

				if err := dbService.Verify(); err != nil {
					return err
				}

				if _, err := db.Query.UpsertService(context.Background(), dbService); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
