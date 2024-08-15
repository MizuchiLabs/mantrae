package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/traefik/genconf/dynamic"
	ttls "github.com/traefik/genconf/dynamic/tls"
	"sigs.k8s.io/yaml"
)

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
	var profiles []Profile

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
	file, err := os.OpenFile(profilePath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open profiles file: %w", err)
	}
	defer file.Close()

	profileBytes, err := json.Marshal(profiles)
	if err != nil {
		return fmt.Errorf("failed to marshal profiles: %w", err)
	}

	_, err = file.Write(profileBytes)
	if err != nil {
		return fmt.Errorf("failed to write profiles: %w", err)
	}

	return nil
}

func VerifyProfile(profile Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if profile.Instance.URL != "" {
		if !isValidURL(profile.Instance.URL) {
			return fmt.Errorf("invalid url")
		}
	} else {
		return fmt.Errorf("url cannot be empty")
	}
	return nil
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port != "" {
		if _, err = net.LookupPort("tcp", port); err != nil {
			return false
		}
	}

	ip := net.ParseIP(host)
	if ip != nil {
		return true
	}

	domainRegex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		return false
	}

	return matched
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
					Priority:    int(router.Priority.Int64()),
					TLS:         router.TLS,
				}
			case "tcp":
				config.TCP.Routers[router.Service] = &dynamic.TCPRouter{
					EntryPoints: router.Entrypoints,
					Middlewares: router.Middlewares,
					Service:     router.Service,
					Priority:    int(router.Priority.Int64()),
					Rule:        router.Rule,
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
	for _, middleware := range dyn.HTTPMiddlewares {
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
		}
	}
	for _, middleware := range dyn.TCPMiddlewares {
		if middleware.Provider == "http" {
			config.TCP.Middlewares[middleware.Name] = &dynamic.TCPMiddleware{
				InFlightConn: middleware.InFlightConn,
				IPWhiteList:  middleware.IPWhiteList,
				IPAllowList:  middleware.IPAllowList,
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

	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	endpoints := map[string]string{
		"entrypoints":     "/api/entrypoints",
		"http-routers":    "/api/http/routers",
		"tcp-routers":     "/api/tcp/routers",
		"udp-routers":     "/api/udp/routers",
		"http-services":   "/api/http/services",
		"tcp-services":    "/api/tcp/services",
		"udp-services":    "/api/udp/services",
		"httpmiddlewares": "/api/http/middlewares",
		"tcpmiddlewares":  "/api/tcp/middlewares",
		"version":         "/api/version",
	}

	for idx, profile := range profiles {
		dynamic := profile.Instance.Dynamic

		for endpoint, url := range endpoints {
			req, err := http.NewRequest("GET", profile.Instance.URL+url, nil)
			if err != nil {
				slog.Error("Failed to create request", "error", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			if profile.Instance.Username != "" && profile.Instance.Password != "" {
				req.SetBasicAuth(profile.Instance.Username, profile.Instance.Password)
			}
			resp, err := client.Do(req)
			if err != nil {
				slog.Error("Failed to make request", "error", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				slog.Error("Unexpected status code", "status", resp.StatusCode)
				return
			}

			switch endpoint {
			case "entrypoints":
				if err := json.NewDecoder(resp.Body).Decode(&dynamic.Entrypoints); err != nil {
					slog.Error("Failed to decode entrypoints", "error", err)
					return
				}
			case "http-routers", "tcp-routers", "udp-routers":
				var routers []Router
				if err := json.NewDecoder(resp.Body).Decode(&routers); err != nil {
					slog.Error("Failed to decode routers", "error", err)
					return
				}
				routerType := strings.Split(endpoint, "-")[0] // Extracts "http", "tcp", or "udp"
				for i := range routers {
					routers[i].RouterType = routerType
				}
				dynamic.Routers = append(dynamic.Routers, routers...)
			case "htttp-services", "tcp-services", "udp-services":
				var services []Service
				if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
					slog.Error("Failed to decode services", "error", err)
					return
				}
				serviceType := strings.Split(endpoint, "-")[0]
				for i := range services {
					services[i].ServiceType = serviceType
				}
				dynamic.Services = append(dynamic.Services, services...)
			case "httpmiddlewares":
				if err := json.NewDecoder(resp.Body).Decode(&dynamic.HTTPMiddlewares); err != nil {
					slog.Error("Failed to decode http middlewares", "error", err)
					return
				}
			case "tcpmiddlewares":
				if err := json.NewDecoder(resp.Body).Decode(&dynamic.TCPMiddlewares); err != nil {
					slog.Error("Failed to decode tcp middlewares", "error", err)
					return
				}
			case "version":
				var version struct {
					Version  string `json:"version"`
					Codename string `json:"codename"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
					slog.Error("Failed to decode version", "error", err)
					return
				}
				dynamic.Version = version.Version
			}
		}
		profiles[idx].Instance.Dynamic = dynamic
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
