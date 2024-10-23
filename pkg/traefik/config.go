package traefik

import (
	"strings"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	ttls "github.com/traefik/traefik/v3/pkg/tls"
)

func GenerateConfig(d Dynamic) (*dynamic.Configuration, error) {
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
			name := strings.Split(router.Name, "@")[0]
			switch router.RouterType {
			case "http":
				var tlsConfig *dynamic.RouterTLSConfig
				if router.TLS != nil {
					tlsConfig = &dynamic.RouterTLSConfig{
						Options:      router.TLS.Options,
						CertResolver: router.TLS.CertResolver,
						Domains:      router.TLS.Domains,
					}
				}
				config.HTTP.Routers[name] = &dynamic.Router{
					EntryPoints: router.Entrypoints,
					Middlewares: router.Middlewares,
					Rule:        router.Rule,
					RuleSyntax:  router.RuleSyntax,
					Service:     router.Service,
					Priority:    router.Priority,
					TLS:         tlsConfig,
				}
			case "tcp":
				config.TCP.Routers[name] = &dynamic.TCPRouter{
					EntryPoints: router.Entrypoints,
					Middlewares: router.Middlewares,
					Rule:        router.Rule,
					RuleSyntax:  router.RuleSyntax,
					Service:     router.Service,
					Priority:    router.Priority,
					TLS:         router.TLS,
				}
			case "udp":
				config.UDP.Routers[name] = &dynamic.UDPRouter{
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
					LoadBalancer: service.LoadBalancer.ToHTTPServersLoadBalancer(),
					Weighted:     service.Weighted.ToHTTPWeightedRoundRobin(),
					Mirroring:    service.Mirroring,
					Failover:     service.Failover,
				}
			case "tcp":
				config.TCP.Services[name] = &dynamic.TCPService{
					LoadBalancer: service.LoadBalancer.ToTCPServersLoadBalancer(),
					Weighted:     service.Weighted.ToTCPWeightedRoundRobin(),
				}
			case "udp":
				config.UDP.Services[name] = &dynamic.UDPService{
					LoadBalancer: service.LoadBalancer.ToUDPServersLoadBalancer(),
					Weighted:     service.Weighted.ToUDPWeightedRoundRobin(),
				}
			}
		}
	}
	for _, middleware := range d.Middlewares {
		if middleware.Provider == "http" {
			name := strings.Split(middleware.Name, "@")[0]
			switch middleware.MiddlewareType {
			case "http":
				config.HTTP.Middlewares[name] = &dynamic.Middleware{
					AddPrefix:         middleware.AddPrefix,
					StripPrefix:       middleware.StripPrefix,
					StripPrefixRegex:  middleware.StripPrefixRegex,
					ReplacePath:       middleware.ReplacePath,
					ReplacePathRegex:  middleware.ReplacePathRegex,
					Chain:             middleware.Chain,
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
					Plugin:            middleware.Plugin,
				}
			case "tcp":
				if middleware.IPAllowList != nil {
					config.TCP.Middlewares[name] = &dynamic.TCPMiddleware{
						IPAllowList: &dynamic.TCPIPAllowList{
							SourceRange: middleware.IPAllowList.SourceRange,
						},
					}
				} else {
					config.TCP.Middlewares[name] = &dynamic.TCPMiddleware{
						InFlightConn: middleware.InFlightConn,
					}
				}
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

	return config, nil
}

// Convert LoadBalancer to HTTP ServersLoadBalancer
func (lb *LoadBalancer) ToHTTPServersLoadBalancer() *dynamic.ServersLoadBalancer {
	if lb == nil {
		return nil
	}

	httpServers := make([]dynamic.Server, len(lb.Servers))
	for _, server := range lb.Servers {
		httpServers = append(httpServers, dynamic.Server{
			URL: server.URL,
		})
	}

	return &dynamic.ServersLoadBalancer{
		Sticky:             lb.Sticky,
		Servers:            httpServers,
		HealthCheck:        lb.HealthCheck,
		PassHostHeader:     lb.PassHostHeader,
		ResponseForwarding: lb.ResponseForwarding,
		ServersTransport:   lb.ServersTransport,
	}
}

// Convert LoadBalancer to TCPServersLoadBalancer
func (lb *LoadBalancer) ToTCPServersLoadBalancer() *dynamic.TCPServersLoadBalancer {
	if lb == nil {
		return nil
	}

	tcpServers := make([]dynamic.TCPServer, len(lb.Servers))
	for _, server := range lb.Servers {
		tcpServers = append(tcpServers, dynamic.TCPServer{
			Address: server.Address,
		})
	}

	return &dynamic.TCPServersLoadBalancer{
		TerminationDelay: lb.TerminationDelay,
		ProxyProtocol:    lb.ProxyProtocol,
		Servers:          tcpServers,
	}
}

// Convert LoadBalancer to UDPServersLoadBalancer (same structure as TCP)
func (lb *LoadBalancer) ToUDPServersLoadBalancer() *dynamic.UDPServersLoadBalancer {
	if lb == nil {
		return nil
	}

	udpServers := make([]dynamic.UDPServer, len(lb.Servers))
	for _, server := range lb.Servers {
		udpServers = append(udpServers, dynamic.UDPServer{
			Address: server.Address,
		})
	}

	return &dynamic.UDPServersLoadBalancer{
		Servers: udpServers,
	}
}

// Convert WeightedRoundRobin to HTTP WeightedRoundRobin
func (wrr *WeightedRoundRobin) ToHTTPWeightedRoundRobin() *dynamic.WeightedRoundRobin {
	if wrr == nil {
		return nil
	}

	httpServices := make([]dynamic.WRRService, len(wrr.Services))
	for _, service := range wrr.Services {
		httpServices = append(httpServices, dynamic.WRRService{
			Name:   service.Name,
			Weight: service.Weight, // HTTP uses Weight
		})
	}

	return &dynamic.WeightedRoundRobin{
		Services:    httpServices,
		Sticky:      wrr.Sticky,
		HealthCheck: wrr.HealthCheck,
	}
}

// Convert WeightedRoundRobin to TCP WeightedRoundRobin
func (wrr *WeightedRoundRobin) ToTCPWeightedRoundRobin() *dynamic.TCPWeightedRoundRobin {
	if wrr == nil {
		return nil
	}

	tcpServices := make([]dynamic.TCPWRRService, len(wrr.Services))
	for _, service := range wrr.Services {
		tcpServices = append(tcpServices, dynamic.TCPWRRService{
			Name:   service.Name,
			Weight: service.Weight,
		})
	}

	return &dynamic.TCPWeightedRoundRobin{
		Services: tcpServices,
	}
}

func (wrr *WeightedRoundRobin) ToUDPWeightedRoundRobin() *dynamic.UDPWeightedRoundRobin {
	if wrr == nil {
		return nil
	}

	udpServices := make([]dynamic.UDPWRRService, len(wrr.Services))
	for _, service := range wrr.Services {
		udpServices = append(udpServices, dynamic.UDPWRRService{
			Name:   service.Name,
			Weight: service.Weight,
		})
	}

	return &dynamic.UDPWeightedRoundRobin{
		Services: udpServices,
	}
}
