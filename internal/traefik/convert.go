package traefik

import (
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

func ConvertToDynamicConfig(rc *db.TraefikConfiguration) *dynamic.Configuration {
	dc := &dynamic.Configuration{}

	// Only create HTTP config if there are HTTP components
	if len(rc.Routers) > 0 || len(rc.Middlewares) > 0 || len(rc.Services) > 0 {
		dc.HTTP = &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		}

		for name, router := range rc.Routers {
			dc.HTTP.Routers[name] = router.Router
		}
		for name, service := range rc.Services {
			dc.HTTP.Services[name] = service.Service
		}
		for name, middleware := range rc.Middlewares {
			dc.HTTP.Middlewares[name] = middleware.Middleware
		}
	}

	// Only create TCP config if there are TCP components
	if len(rc.TCPRouters) > 0 || len(rc.TCPMiddlewares) > 0 || len(rc.TCPServices) > 0 {
		dc.TCP = &dynamic.TCPConfiguration{
			Routers:     make(map[string]*dynamic.TCPRouter),
			Middlewares: make(map[string]*dynamic.TCPMiddleware),
			Services:    make(map[string]*dynamic.TCPService),
		}

		for name, router := range rc.TCPRouters {
			dc.TCP.Routers[name] = router.TCPRouter
		}
		for name, service := range rc.TCPServices {
			dc.TCP.Services[name] = service.TCPService
		}
		for name, middleware := range rc.TCPMiddlewares {
			dc.TCP.Middlewares[name] = middleware.TCPMiddleware
		}
	}

	// Only create UDP config if there are UDP components
	if len(rc.UDPRouters) > 0 || len(rc.UDPServices) > 0 {
		dc.UDP = &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		}

		for name, router := range rc.UDPRouters {
			dc.UDP.Routers[name] = router.UDPRouter
		}
		for name, service := range rc.UDPServices {
			dc.UDP.Services[name] = service.UDPService
		}
	}

	return dc
}

func MergeConfigs(base, overlay *db.TraefikConfiguration) *db.TraefikConfiguration {
	if overlay == nil {
		return base
	}

	// Merge HTTP components
	mergeRouters(base.Routers, overlay.Routers)
	mergeMiddlewares(base.Middlewares, overlay.Middlewares)
	mergeServices(base.Services, overlay.Services)

	// Merge TCP components
	mergeTCPRouters(base.TCPRouters, overlay.TCPRouters)
	mergeTCPMiddlewares(base.TCPMiddlewares, overlay.TCPMiddlewares)
	mergeTCPServices(base.TCPServices, overlay.TCPServices)

	// Merge UDP components
	mergeUDPRouters(base.UDPRouters, overlay.UDPRouters)
	mergeUDPServices(base.UDPServices, overlay.UDPServices)

	return base
}

// Merge helper functions for each type
func mergeRouters(base, overlay map[string]*runtime.RouterInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeMiddlewares(base, overlay map[string]*runtime.MiddlewareInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeServices(base, overlay map[string]*db.ServiceInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeTCPRouters(base, overlay map[string]*runtime.TCPRouterInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeTCPMiddlewares(base, overlay map[string]*runtime.TCPMiddlewareInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeTCPServices(base, overlay map[string]*runtime.TCPServiceInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeUDPRouters(base, overlay map[string]*runtime.UDPRouterInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}

func mergeUDPServices(base, overlay map[string]*runtime.UDPServiceInfo) {
	if overlay == nil {
		return
	}
	for k, v := range overlay {
		base[k] = v
	}
}
