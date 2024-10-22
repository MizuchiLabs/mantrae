package traefik

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

const (
	HTTPRouterAPI      = "/api/http/routers"
	TCPRouterAPI       = "/api/tcp/routers"
	UDPRouterAPI       = "/api/udp/routers"
	HTTPServiceAPI     = "/api/http/services"
	TCPServiceAPI      = "/api/tcp/services"
	UDPServiceAPI      = "/api/udp/services"
	HTTPMiddlewaresAPI = "/api/http/middlewares"
	TCPMiddlewaresAPI  = "/api/tcp/middlewares"
	OverviewAPI        = "/api/overview"
	EntrypointsAPI     = "/api/entrypoints"
	VersionAPI         = "/api/version"
)

// Extra fields from the API endpoint
type BaseFields struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Status   string `json:"status,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// Extended routers
type HTTPRouter struct {
	BaseFields
	RouterType  string `json:"routerType"`
	DNSProvider string `json:"dnsProvider"`
	SSLError    string `json:"sslError,omitempty"`
	dynamic.Router
}

type TCPRouter struct {
	BaseFields
	RouterType  string `json:"routerType"`
	DNSProvider string `json:"dnsProvider"`
	dynamic.TCPRouter
}

type UDPRouter struct {
	BaseFields
	RouterType  string `json:"routerType"`
	DNSProvider string `json:"dnsProvider"`
	dynamic.UDPRouter
}

type Routerable interface {
	ToRouter() *Router
}

func (r HTTPRouter) ToRouter() *Router {
	var tlsConfig *dynamic.RouterTCPTLSConfig

	if r.TLS != nil {
		tlsConfig = &dynamic.RouterTCPTLSConfig{
			Options:      r.TLS.Options,
			CertResolver: r.TLS.CertResolver,
			Domains:      r.TLS.Domains,
		}
	}

	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "http",
		Entrypoints: r.EntryPoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		RuleSyntax:  r.RuleSyntax,
		Service:     r.Service,
		Priority:    r.Priority,
		TLS:         tlsConfig,
	}
}

func (r TCPRouter) ToRouter() *Router {
	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "tcp",
		Entrypoints: r.EntryPoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		RuleSyntax:  r.RuleSyntax,
		Service:     r.Service,
		Priority:    r.Priority,
		TLS:         r.TLS,
	}
}

func (r UDPRouter) ToRouter() *Router {
	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "udp",
		Entrypoints: r.EntryPoints,
		Service:     r.Service,
	}
}

func getRouters[T Routerable](profile db.Profile, endpoint string) map[string]Router {
	body, err := fetch(profile, endpoint)
	if err != nil {
		slog.Error("Failed to get routers", "error", err)
		return nil
	}
	defer body.Close()
	if body == nil {
		return nil
	}
	var routerables []T
	if err := json.NewDecoder(body).Decode(&routerables); err != nil {
		slog.Error("Failed to decode routers", "error", err)
		return nil
	}

	routers := make(map[string]Router, len(routerables))
	for _, r := range routerables {
		newRouter := r.ToRouter()
		if newRouter.Name == "" {
			continue
		}
		routers[newRouter.Name] = *newRouter
	}
	return routers
}

type HTTPService struct {
	BaseFields
	ServiceType  string            `json:"serviceType,omitempty"`
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.Service
}

type TCPService struct {
	BaseFields
	ServiceType  string            `json:"serviceType,omitempty"`
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.TCPService
}

type UDPService struct {
	BaseFields
	ServiceType  string            `json:"serviceType,omitempty"`
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.UDPService
}

type Serviceable interface {
	ToService() *Service
}

func (s HTTPService) ToService() *Service {
	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				servers[i] = Server{URL: server.URL}
			}
		}
		lb = &LoadBalancer{
			Servers:            servers,
			Sticky:             s.LoadBalancer.Sticky,
			PassHostHeader:     s.LoadBalancer.PassHostHeader,
			HealthCheck:        s.LoadBalancer.HealthCheck,
			ResponseForwarding: s.LoadBalancer.ResponseForwarding,
			ServersTransport:   s.LoadBalancer.ServersTransport,
		}
	}

	var weighted *WeightedRoundRobin
	if s.Weighted != nil {
		weighted = &WeightedRoundRobin{
			Services:    s.Weighted.Services,
			Sticky:      s.Weighted.Sticky,
			HealthCheck: s.Weighted.HealthCheck,
		}
	}

	return &Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       s.Status,
		ServiceType:  "http",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
		Mirroring:    s.Mirroring,
		Failover:     s.Failover,
	}
}

func (s TCPService) ToService() *Service {
	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				servers[i] = Server{Address: server.Address}
			}
		}

		lb = &LoadBalancer{
			Servers:          servers,
			ServersTransport: s.LoadBalancer.ServersTransport,
			ProxyProtocol:    s.LoadBalancer.ProxyProtocol,
			TerminationDelay: s.LoadBalancer.TerminationDelay,
		}
	}

	var weighted *WeightedRoundRobin
	var services []dynamic.WRRService
	if s.Weighted != nil && s.Weighted.Services != nil {
		for _, service := range s.Weighted.Services {
			services = append(services, dynamic.WRRService{
				Name:   service.Name,
				Weight: service.Weight,
			})
		}

		weighted = &WeightedRoundRobin{Services: services}
	}
	return &Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       s.Status,
		ServiceType:  "tcp",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
	}
}

func (s UDPService) ToService() *Service {
	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil && s.LoadBalancer.Servers != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				servers[i] = Server{
					Address: server.Address,
				}
			}
		}
		lb = &LoadBalancer{
			Servers: servers,
		}
	}

	var weighted *WeightedRoundRobin
	var services []dynamic.WRRService
	if s.Weighted != nil && s.Weighted.Services != nil {
		for _, service := range s.Weighted.Services {
			services = append(services, dynamic.WRRService{
				Name:   service.Name,
				Weight: service.Weight,
			})
		}

		weighted = &WeightedRoundRobin{Services: services}
	}

	return &Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       s.Status,
		ServiceType:  "udp",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
	}
}

func getServices[T Serviceable](profile db.Profile, endpoint string) map[string]Service {
	body, err := fetch(profile, endpoint)
	if err != nil {
		slog.Error("Failed to get services", "error", err)
		return nil
	}
	defer body.Close()

	var serviceables []T
	if err := json.NewDecoder(body).Decode(&serviceables); err != nil {
		slog.Error("Failed to decode services", "error", err)
		return nil
	}

	services := make(map[string]Service, len(serviceables))
	for _, s := range serviceables {
		newService := s.ToService()
		if newService.Name == "" {
			continue
		}
		services[newService.Name] = *newService
	}

	return services
}

type HTTPMiddleware struct {
	BaseFields
	MiddlewareType string `json:"middlewareType,omitempty"`
	dynamic.Middleware
}

type TCPMiddleware struct {
	BaseFields
	MiddlewareType string `json:"middlewareType,omitempty"`
	dynamic.TCPMiddleware
}

type Middlewareable interface {
	ToMiddleware() *Middleware
}

func (m HTTPMiddleware) ToMiddleware() *Middleware {
	return &Middleware{
		Name:              m.Name,
		Provider:          m.Provider,
		Type:              m.Type,
		Status:            m.Status,
		MiddlewareType:    "http",
		AddPrefix:         m.AddPrefix,
		StripPrefix:       m.StripPrefix,
		StripPrefixRegex:  m.StripPrefixRegex,
		ReplacePath:       m.ReplacePath,
		ReplacePathRegex:  m.ReplacePathRegex,
		Chain:             m.Chain,
		IPAllowList:       m.IPAllowList,
		Headers:           m.Headers,
		Errors:            m.Errors,
		RateLimit:         m.RateLimit,
		RedirectRegex:     m.RedirectRegex,
		RedirectScheme:    m.RedirectScheme,
		BasicAuth:         m.BasicAuth,
		DigestAuth:        m.DigestAuth,
		ForwardAuth:       m.ForwardAuth,
		InFlightReq:       m.InFlightReq,
		Buffering:         m.Buffering,
		CircuitBreaker:    m.CircuitBreaker,
		Compress:          m.Compress,
		PassTLSClientCert: m.PassTLSClientCert,
		Retry:             m.Retry,
		GrpcWeb:           m.GrpcWeb,
		Plugin:            m.Plugin,
	}
}

func (m TCPMiddleware) ToMiddleware() *Middleware {
	var allowList *dynamic.IPAllowList
	if m.IPAllowList != nil {
		allowList = &dynamic.IPAllowList{
			SourceRange: m.IPAllowList.SourceRange,
		}
	}

	return &Middleware{
		Name:           m.Name,
		Provider:       m.Provider,
		Type:           m.Type,
		Status:         m.Status,
		MiddlewareType: "tcp",
		InFlightConn:   m.InFlightConn,
		IPAllowList:    allowList,
	}
}

func getMiddlewares[T Middlewareable](profile db.Profile, endpoint string) map[string]Middleware {
	body, err := fetch(profile, endpoint)
	if err != nil {
		slog.Error("Failed to get middlewares", "error", err)
		return nil
	}
	defer body.Close()

	var middlewareables []T
	if err := json.NewDecoder(body).Decode(&middlewareables); err != nil {
		slog.Error("Failed to decode middlewareables", "error", err)
		return nil
	}

	middlewares := make(map[string]Middleware, len(middlewareables))
	for _, m := range middlewareables {
		newMiddleware := m.ToMiddleware()
		if newMiddleware.Name == "" {
			continue
		}
		middlewares[newMiddleware.Name] = *newMiddleware
	}

	return middlewares
}

func GetTraefikConfig() {
	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
		return
	}

	for _, profile := range profiles {
		if profile.Url == "" {
			continue
		}

		config, err := db.Query.GetConfigByProfileID(context.Background(), profile.ID)
		if err != nil {
			slog.Error("Failed to get config", "error", err)
			return
		}

		data, err := DecodeFromDB(config)
		if err != nil {
			slog.Error("Failed to decode config", "error", err)
			return
		}

		// Fetch routers
		data.Routers = merge(
			data.Routers,
			getRouters[HTTPRouter](profile, HTTPRouterAPI),
			getRouters[TCPRouter](profile, TCPRouterAPI),
			getRouters[UDPRouter](profile, UDPRouterAPI),
		)

		// Fetch services
		data.Services = merge(
			data.Services,
			getServices[HTTPService](profile, HTTPServiceAPI),
			getServices[TCPService](profile, TCPServiceAPI),
			getServices[UDPService](profile, UDPServiceAPI),
		)

		// Fetch middlewares
		data.Middlewares = merge(
			data.Middlewares,
			getMiddlewares[HTTPMiddleware](profile, HTTPMiddlewaresAPI),
			getMiddlewares[TCPMiddleware](profile, TCPMiddlewaresAPI),
		)

		// Fetch overview
		overview, err := fetch(profile, OverviewAPI)
		if err != nil {
			slog.Error("Failed to get overview", "error", err)
			return
		}
		defer overview.Close()

		var dataOverview Overview
		if err = json.NewDecoder(overview).Decode(&dataOverview); err != nil {
			slog.Error("Failed to decode overview", "error", err)
			return
		}
		data.Overview = &dataOverview

		// Retrieve entrypoints
		entrypoints, err := fetch(profile, EntrypointsAPI)
		if err != nil {
			slog.Error("Failed to get entrypoints", "error", err)
			return
		}
		defer entrypoints.Close()

		var dataEntrypoints []Entrypoint
		if err = json.NewDecoder(entrypoints).Decode(&dataEntrypoints); err != nil {
			slog.Error("Failed to decode entrypoints", "error", err)
			return
		}
		data.Entrypoints = dataEntrypoints

		// Fetch version
		version, err := fetch(profile, VersionAPI)
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
		data.Version = v.Version

		// Verify config and write to db
		if _, err := EncodeToDB(config.ProfileID, data); err != nil {
			slog.Error("Failed to update config", "error", err)
			return
		}
	}

	// Broadcast the update to all clients
	util.Broadcast <- "profiles"
}

// Sync periodically syncs the Traefik configuration
func Sync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	GetTraefikConfig()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			GetTraefikConfig()
		}
	}
}

func merge[T any](local map[string]T, externals ...map[string]T) map[string]T {
	merged := make(map[string]T)

	// Add local provider ("http") and DNSProvider-preserving routers to merged
	for k, v := range local {
		switch item := any(v).(type) {
		case Router:
			if item.Provider == "http" || item.DNSProvider != nil {
				merged[k] = v
			}
		case Service:
			if item.Provider == "http" {
				merged[k] = v
			}
		case Middleware:
			if item.Provider == "http" {
				merged[k] = v
			}
		}
	}

	// Merge in external data without overwriting local "http" provider entries
	for _, external := range externals {
		for k, v := range external {
			if existing, found := merged[k]; found {
				// If exists, check and update for specific fields (e.g., DNSProvider)
				switch existingItem := any(existing).(type) {
				case Router:
					if newRouter, ok := any(v).(Router); ok {
						newRouter.DNSProvider = existingItem.DNSProvider
						newRouter.SSLError = existingItem.SSLError
						merged[k] = any(newRouter).(T)
					}
				default:
					// Services or Middleware might not need this DNSProvider logic
					merged[k] = v
				}
			} else {
				// Add non-http provider entries
				switch newItem := any(v).(type) {
				case Router:
					if newItem.Provider != "http" {
						merged[k] = v
					}
				case Service:
					if newItem.Provider != "http" {
						merged[k] = v
					}
				case Middleware:
					if newItem.Provider != "http" {
						merged[k] = v
					}
				default:
					merged[k] = v
				}
			}
		}
	}

	return merged
}

func fetch(profile db.Profile, endpoint string) (io.ReadCloser, error) {
	if profile.Url == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	client := http.Client{Timeout: time.Second * 10}
	if !profile.Tls {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("GET", profile.Url+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if *profile.Username != "" && *profile.Password != "" {
		req.SetBasicAuth(*profile.Username, *profile.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", profile.Url+endpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
