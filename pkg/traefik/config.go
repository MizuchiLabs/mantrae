package traefik

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	ttls "github.com/traefik/traefik/v3/pkg/tls"
)

func GenerateConfig(d *Dynamic) (*dynamic.Configuration, error) {
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

	routers, err := db.Query.ListRoutersByProfileID(context.Background(), d.ProfileID)
	if err != nil {
		return nil, err
	}

	for _, router := range routers {
		if err := router.DecodeFields(); err != nil {
			continue
		}

		// Only add routers by our provider
		if router.Provider == "http" {
			name := strings.Split(router.Name, "@")[0]
			switch router.Protocol {
			case "http":
				var rtr *dynamic.Router
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.HTTP.Routers[name] = rtr
			case "tcp":
				var rtr *dynamic.TCPRouter
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}

				config.TCP.Routers[name] = rtr
			case "udp":
				var rtr *dynamic.UDPRouter
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.UDP.Routers[name] = rtr
			}
		}
	}

	services, err := db.Query.ListServicesByProfileID(context.Background(), d.ProfileID)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if err := service.DecodeFields(); err != nil {
			continue
		}

		// Only add services by our provider
		if service.Provider == "http" {
			name := strings.Split(service.Name, "@")[0]
			switch service.Protocol {
			case "http":
				var svc *dynamic.Service
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.HTTP.Services[name] = svc
			case "tcp":
				var svc *dynamic.TCPService
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.TCP.Services[name] = svc
			case "udp":
				var svc *dynamic.UDPService
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.UDP.Services[name] = svc
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
