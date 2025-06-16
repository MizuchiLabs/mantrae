package traefik

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
	"golang.org/x/exp/maps"
	"sigs.k8s.io/yaml"
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

func ConvertFileToDynamicConfig(
	file multipart.File,
	extension string,
) (*dynamic.Configuration, error) {
	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a new configuration
	configuration := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:           make(map[string]*dynamic.Router),
			Middlewares:       make(map[string]*dynamic.Middleware),
			Services:          make(map[string]*dynamic.Service),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:           make(map[string]*dynamic.TCPRouter),
			Services:          make(map[string]*dynamic.TCPService),
			Middlewares:       make(map[string]*dynamic.TCPMiddleware),
			ServersTransports: make(map[string]*dynamic.TCPServersTransport),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
	}

	// Unmarshal yaml or json
	if extension == ".json" {
		if err := json.Unmarshal(content, configuration); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	}
	if extension == ".yaml" || extension == ".yml" {
		if err := yaml.Unmarshal(content, configuration); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}
	return configuration, nil
}

func ConvertDynamicToRuntime(dc *dynamic.Configuration) *db.TraefikConfiguration {
	rc := runtime.NewConfig(*dc)
	config := &db.TraefikConfiguration{
		Routers:        rc.Routers,
		Middlewares:    rc.Middlewares,
		TCPRouters:     rc.TCPRouters,
		TCPMiddlewares: rc.TCPMiddlewares,
		TCPServices:    rc.TCPServices,
		UDPRouters:     rc.UDPRouters,
		UDPServices:    rc.UDPServices,
	}
	config.Services = make(map[string]*db.ServiceInfo)
	for k, v := range rc.Services {
		config.Services[k] = db.FromRuntimeServiceInfo(v)
	}
	return config
}

func MergeConfigs(base, overlay *db.TraefikConfiguration) *db.TraefikConfiguration {
	if base == nil {
		return overlay
	}
	if overlay == nil {
		return base
	}

	initializeMaps(base)

	// Merge components
	mergeMaps(base.Routers, overlay.Routers)
	mergeMaps(base.Middlewares, overlay.Middlewares)
	mergeMaps(base.Services, overlay.Services)
	mergeMaps(base.TCPRouters, overlay.TCPRouters)
	mergeMaps(base.TCPMiddlewares, overlay.TCPMiddlewares)
	mergeMaps(base.TCPServices, overlay.TCPServices)
	mergeMaps(base.UDPRouters, overlay.UDPRouters)
	mergeMaps(base.UDPServices, overlay.UDPServices)

	return base
}

// Initialize nil maps in the configuration
func initializeMaps(config *db.TraefikConfiguration) {
	if config.Routers == nil {
		config.Routers = make(map[string]*runtime.RouterInfo)
	}
	if config.Middlewares == nil {
		config.Middlewares = make(map[string]*runtime.MiddlewareInfo)
	}
	if config.Services == nil {
		config.Services = make(map[string]*db.ServiceInfo)
	}
	if config.TCPRouters == nil {
		config.TCPRouters = make(map[string]*runtime.TCPRouterInfo)
	}
	if config.TCPMiddlewares == nil {
		config.TCPMiddlewares = make(map[string]*runtime.TCPMiddlewareInfo)
	}
	if config.TCPServices == nil {
		config.TCPServices = make(map[string]*runtime.TCPServiceInfo)
	}
	if config.UDPRouters == nil {
		config.UDPRouters = make(map[string]*runtime.UDPRouterInfo)
	}
	if config.UDPServices == nil {
		config.UDPServices = make(map[string]*runtime.UDPServiceInfo)
	}
}

// Generic merge function for maps
func mergeMaps[K comparable, V any](base, overlay map[K]V) {
	if base == nil || overlay == nil {
		return
	}

	maps.Copy(base, overlay)
}
