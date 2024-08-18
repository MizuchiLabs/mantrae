package traefik

import (
	"strings"

	"github.com/traefik/genconf/dynamic"
	ttls "github.com/traefik/genconf/dynamic/tls"
	"sigs.k8s.io/yaml"
)

func GenerateConfig(d Dynamic) ([]byte, error) {
	config := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:           make(map[string]*dynamic.Router),
			Middlewares:       make(map[string]*dynamic.Middleware),
			Services:          make(map[string]*dynamic.Service),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:     make(map[string]*dynamic.TCPRouter),
			Services:    make(map[string]*dynamic.TCPService),
			Middlewares: make(map[string]*dynamic.TCPMiddleware),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]ttls.Store),
			Options: make(map[string]ttls.Options),
		},
	}

	for _, router := range d.Routers {
		// Only add routers by our provider
		if router.Provider == "http" {
			switch router.RouterType {
			case "http":
				config.HTTP.Routers[router.Service] = &dynamic.Router{
					EntryPoints: router.Entrypoints,
					Middlewares: router.Middlewares,
					Service:     router.Service,
					Rule:        router.Rule,
					// Priority:    int(router.Priority.Int64()),
					TLS: router.TLS,
				}
			case "tcp":
				config.TCP.Routers[router.Service] = &dynamic.TCPRouter{
					EntryPoints: router.Entrypoints,
					Middlewares: router.Middlewares,
					Service:     router.Service,
					// Priority:    int(router.Priority.Int64()),
					Rule: router.Rule,
				}
			case "udp":
				config.UDP.Routers[router.Service] = &dynamic.UDPRouter{
					EntryPoints: router.Entrypoints,
					Service:     router.Service,
				}
			}
		}
	}
	for _, service := range d.Services {
		// Only add services by our provider
		if service.Provider == "http" {
			name := strings.Split(service.Name, "@")[0]
			switch service.ServiceType {
			case "http":
				config.HTTP.Services[name] = &dynamic.Service{
					LoadBalancer: service.LoadBalancer,
					Weighted:     service.Weighted,
					Mirroring:    service.Mirroring,
					Failover:     service.Failover,
				}
			case "tcp":
				config.TCP.Services[name], _ = convertService(service)
			case "udp":
				_, config.UDP.Services[name] = convertService(service)
			}
		}
	}
	for _, middleware := range d.Middlewares {
		if middleware.Provider == "http" {
			config.HTTP.Middlewares[middleware.Name] = &dynamic.Middleware{
				AddPrefix:         middleware.AddPrefix,
				StripPrefix:       middleware.StripPrefix,
				StripPrefixRegex:  middleware.StripPrefixRegex,
				ReplacePath:       middleware.ReplacePath,
				ReplacePathRegex:  middleware.ReplacePathRegex,
				Chain:             middleware.Chain,
				IPWhiteList:       middleware.IPWhiteList,
				IPAllowList:       middleware.IPAllowList,
				Headers:           middleware.Headers,
				Errors:            middleware.Errors,
				RateLimit:         middleware.RateLimit,
				RedirectRegex:     middleware.RedirectRegex,
				RedirectScheme:    middleware.RedirectScheme,
				BasicAuth:         middleware.BasicAuth,
				DigestAuth:        middleware.DigestAuth,
				ForwardAuth:       middleware.ForwardAuth,
				InFlightReq:       middleware.InFlightReq,
				Buffering:         middleware.Buffering,
				CircuitBreaker:    middleware.CircuitBreaker,
				Compress:          middleware.Compress,
				PassTLSClientCert: middleware.PassTLSClientCert,
				Retry:             middleware.Retry,
				ContentType:       middleware.ContentType,
				Plugin:            middleware.Plugin,
			}
			config.TCP.Middlewares[middleware.Name] = &dynamic.TCPMiddleware{
				InFlightConn: middleware.InFlightConn,
			}
		}
	}

	// Remove empty configurations
	if len(config.HTTP.Routers) == 0 && len(config.HTTP.Services) == 0 &&
		len(config.HTTP.Middlewares) == 0 {
		config.HTTP = nil
	}
	if len(config.TCP.Routers) == 0 && len(config.TCP.Services) == 0 &&
		len(config.TCP.Middlewares) == 0 {
		config.TCP = nil
	}
	if len(config.UDP.Routers) == 0 && len(config.UDP.Services) == 0 {
		config.UDP = nil
	}
	if len(config.TLS.Stores) == 0 && len(config.TLS.Options) == 0 {
		config.TLS = nil
	}

	yamlConfig, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return yamlConfig, nil
}

func convertService(service Service) (*dynamic.TCPService, *dynamic.UDPService) {
	var tcpServer []dynamic.TCPServer
	var udpServer []dynamic.UDPServer

	for _, lb := range service.LoadBalancer.Servers {
		if lb.URL != "" {
			tcpServer = append(tcpServer, dynamic.TCPServer{
				Address: lb.URL,
			})
			udpServer = append(udpServer, dynamic.UDPServer{
				Address: lb.URL,
			})
		}
	}
	tcpService := &dynamic.TCPService{
		LoadBalancer: &dynamic.TCPServersLoadBalancer{
			Servers: tcpServer,
		},
	}
	udpService := &dynamic.UDPService{
		LoadBalancer: &dynamic.UDPServersLoadBalancer{
			Servers: udpServer,
		},
	}
	return tcpService, udpService
}
