package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/traefik/genconf/dynamic"
	ttls "github.com/traefik/genconf/dynamic/tls"
	"sigs.k8s.io/yaml"
)

var rwMutex sync.RWMutex

func profilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	newFilePath := filepath.Join(cwd, "profiles.json")
	return newFilePath
}

func defaultProfile() Profile {
	return Profile{
		Name: "default",
		Instance: Instance{
			URL:      "http://127.0.0.1:8080",
			Username: "",
			Password: "",
			Dynamic:  Dynamic{},
		},
	}
}

func LoadProfiles() ([]Profile, error) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	profiles := []Profile{}

	if _, err := os.Stat(profilePath()); os.IsNotExist(err) {
		profiles = append(profiles, defaultProfile())
		if err := SaveProfiles(profiles); err != nil {
			slog.Error("Failed to save profiles", "error", err)
		}
		return profiles, nil
	}

	file, err := os.ReadFile(profilePath())
	if err != nil {
		return []Profile{}, fmt.Errorf("failed to read profiles file: %w", err)
	}

	if err := json.Unmarshal(file, &profiles); err != nil {
		return []Profile{}, fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	return profiles, nil
}

func SaveProfiles(profiles []Profile) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	tmpFile, err := os.CreateTemp(os.TempDir(), "profiles-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	profileBytes, err := json.Marshal(profiles)
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	_, err = tmpFile.Write(profileBytes)
	if err != nil {
		return fmt.Errorf("failed to write profiles: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	tmpFile.Close()

	if err := Move(tmpFile.Name(), profilePath()); err != nil {
		return fmt.Errorf("failed to move temp file: %w", err)
	}

	return nil
}

func Move(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(source, destination)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return fmt.Errorf("failed to stat source file: %w", err)
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return fmt.Errorf("failed to chmod destination file: %w", err)
	}
	os.Remove(source)
	return nil
}

func ParseConfig(dyn Dynamic) ([]byte, error) {
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

	for _, router := range dyn.Routers {
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
	for _, service := range dyn.Services {
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
	for _, middleware := range dyn.Middlewares {
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

func FetchTraefikConfig() {
	profiles, err := LoadProfiles()
	if err != nil {
		slog.Error("Failed to load profiles", "error", err)
		return
	}

	for idx, profile := range profiles {
		i := profile.Instance
		d := Dynamic{
			Routers:     make(map[string]Router),
			Services:    make(map[string]Service),
			Middlewares: make(map[string]Middleware),
		}

		// Retrieve routers
		var tRouter []Router
		tRouter = append(tRouter, fetchRouters[HTTPRouter](i, HTTPRouterAPI)...)
		tRouter = append(tRouter, fetchRouters[TCPRouter](i, TCPRouterAPI)...)
		tRouter = append(tRouter, fetchRouters[UDPRouter](i, UDPRouterAPI)...)
		for _, r := range tRouter {
			d.Routers[r.Name] = r
		}
		for _, r := range profile.Instance.Dynamic.Routers {
			d.Routers[r.Name] = r
		}

		// Retrieve services
		var tServices []Service
		tServices = append(tServices, fetchServices[HTTPService](i, HTTPServiceAPI)...)
		tServices = append(tServices, fetchServices[TCPService](i, TCPServiceAPI)...)
		tServices = append(tServices, fetchServices[UDPService](i, UDPServiceAPI)...)
		for _, s := range tServices {
			d.Services[s.Name] = s
		}
		for _, s := range profile.Instance.Dynamic.Services {
			d.Services[s.Name] = s
		}

		// Fetch middlewares
		var tMiddlewares []Middleware
		tMiddlewares = append(
			tMiddlewares,
			fetchMiddlewares[HTTPMiddleware](i, HTTPMiddlewaresAPI)...)
		tMiddlewares = append(
			tMiddlewares,
			fetchMiddlewares[TCPMiddleware](i, TCPMiddlewaresAPI)...)
		for _, m := range tMiddlewares {
			d.Middlewares[m.Name] = m
		}
		for _, m := range profile.Instance.Dynamic.Middlewares {
			d.Middlewares[m.Name] = m
		}

		// Retrieve entrypoints
		entrypoints, err := get(i, EntrypointsAPI)
		if err != nil {
			slog.Error("Failed to get entrypoints", "error", err)
			return
		}
		defer entrypoints.Close()

		if err = json.NewDecoder(entrypoints).Decode(&d.Entrypoints); err != nil {
			slog.Error("Failed to decode entrypoints", "error", err)
			return
		}

		// Fetch version
		version, err := get(i, VersionAPI)
		if err != nil {
			slog.Error("Failed to get version", "error", err)
			return
		}
		defer version.Close()

		var v struct {
			Version string `json:"version"`
		}

		if err = json.NewDecoder(version).Decode(&v); err != nil {
			slog.Error("Failed to decode version", "error", err)
			return
		}
		d.Version = v.Version

		profiles[idx].Instance.Dynamic = d
	}

	if err := SaveProfiles(profiles); err != nil {
		slog.Error("Failed to save profiles", "error", err)
	}
}

// Sync periodically syncs the Traefik configuration
func Sync(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for range ticker.C {
		FetchTraefikConfig()
	}
}
