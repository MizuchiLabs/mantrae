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
	"github.com/google/uuid"
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
	Protocol string `json:"protocol,omitempty"`
}

// Extended routers
type HTTPRouter struct {
	BaseFields
	dynamic.Router
}

type TCPRouter struct {
	BaseFields
	dynamic.TCPRouter
}

type UDPRouter struct {
	BaseFields
	dynamic.UDPRouter
}

type Routerable interface {
	ToRouter() *db.Router
}

func (r HTTPRouter) ToRouter() *db.Router {
	var status *string
	if r.Status != "" {
		status = &r.Status
	}

	var rulesyntax *string
	if r.RuleSyntax != "" {
		rulesyntax = &r.RuleSyntax
	}

	var priority *int64
	if r.Priority != 0 {
		p := int64(r.Priority)
		priority = &p
	}

	var tlsConfig *dynamic.RouterTCPTLSConfig
	if r.TLS != nil {
		tlsConfig = &dynamic.RouterTCPTLSConfig{
			Options:      r.TLS.Options,
			CertResolver: r.TLS.CertResolver,
			Domains:      r.TLS.Domains,
		}
	}

	return &db.Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      status,
		Protocol:    "http",
		EntryPoints: r.EntryPoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		RuleSyntax:  rulesyntax,
		Service:     r.Service,
		Priority:    priority,
		Tls:         tlsConfig,
	}
}

func (r TCPRouter) ToRouter() *db.Router {
	var status *string
	if r.Status != "" {
		status = &r.Status
	}

	var rulesyntax *string
	if r.RuleSyntax != "" {
		rulesyntax = &r.RuleSyntax
	}

	var priority *int64
	if r.Priority != 0 {
		p := int64(r.Priority)
		priority = &p
	}

	var tlsConfig *dynamic.RouterTCPTLSConfig
	if r.TLS != nil {
		tlsConfig = &dynamic.RouterTCPTLSConfig{
			Options:      r.TLS.Options,
			CertResolver: r.TLS.CertResolver,
			Domains:      r.TLS.Domains,
		}
	}

	return &db.Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      status,
		Protocol:    "tcp",
		EntryPoints: r.EntryPoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		RuleSyntax:  rulesyntax,
		Service:     r.Service,
		Priority:    priority,
		Tls:         tlsConfig,
	}
}

func (r UDPRouter) ToRouter() *db.Router {
	var status *string
	if r.Status != "" {
		status = &r.Status
	}

	return &db.Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      status,
		Protocol:    "udp",
		EntryPoints: r.EntryPoints,
		Service:     r.Service,
	}
}

func getRouters[T Routerable](profile db.Profile, endpoint string) error {
	body, err := fetch(profile, endpoint)
	if err != nil {
		return fmt.Errorf("failed to get routers: %w", err)
	}
	defer body.Close()
	if body == nil {
		return nil
	}
	var routerables []T
	if err := json.NewDecoder(body).Decode(&routerables); err != nil {
		return fmt.Errorf("failed to decode routers: %w", err)
	}

	// Current routers
	dbRouters, err := db.Query.ListRoutersByProfileID(context.Background(), profile.ID)
	if err != nil {
		return fmt.Errorf("failed to list routers: %w", err)
	}
	// fmt.Printf("dbRouters: %+v\n", dbRouters)

	// Create a map to quickly look up existing routers by name
	existingRouters := make(map[string]string, len(dbRouters)) // name to ID mapping
	for _, dbRouter := range dbRouters {
		if dbRouter.Name == "" {
			continue
		}
		existingRouters[dbRouter.Name] = dbRouter.ID
	}
	fmt.Printf("existingRouters: %v\n", existingRouters)

	var currentProtocol string // of T
	routers := make(map[string]db.Router, len(routerables))
	for _, r := range routerables {
		newRouter := r.ToRouter()
		if newRouter.Name == "" {
			continue
		}
		routers[newRouter.Name] = *newRouter
		currentProtocol = newRouter.Protocol

		data := db.UpsertRouterParams{
			ProfileID:  profile.ID,
			Provider:   newRouter.Provider,
			Protocol:   newRouter.Protocol,
			Status:     newRouter.Status,
			Rule:       newRouter.Rule,
			RuleSyntax: newRouter.RuleSyntax,
			Service:    newRouter.Service,
			Priority:   newRouter.Priority,
		}
		// Get existing ID or generate a new one
		data.ID = existingRouters[newRouter.Name]
		if data.ID == "" {
			data.ID = uuid.New().String()
			data.Name = newRouter.Name
		}

		data.EntryPoints, _ = json.Marshal(newRouter.EntryPoints)
		data.Middlewares, _ = json.Marshal(newRouter.Middlewares)
		data.Tls, _ = json.Marshal(newRouter.Tls)
		if _, err := db.Query.UpsertRouter(context.Background(), data); err != nil {
			slog.Error("Failed to upsert router", "error", err)
			continue
		}
	}

	// Cleanup if router doesn't exist locally (except our provider)
	for _, r := range dbRouters {
		if r.Protocol != currentProtocol {
			continue
		}

		if _, ok := routers[r.Name]; !ok && r.Provider != "http" {
			if err := db.Query.DeleteRouterByID(context.Background(), r.ID); err != nil {
				slog.Error("failed to delete router", "error", err)
				continue
			}
		}
	}
	return nil
}

type HTTPService struct {
	BaseFields
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.Service
}

type TCPService struct {
	BaseFields
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.TCPService
}

type UDPService struct {
	BaseFields
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	dynamic.UDPService
}

type Serviceable interface {
	ToService() *db.Service
}

func (s HTTPService) ToService() *db.Service {
	var status *string
	if s.Status != "" {
		status = &s.Status
	}

	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				if server.URL == "" {
					continue
				}
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

	return &db.Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       status,
		Protocol:     "http",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
		Mirroring:    s.Mirroring,
		Failover:     s.Failover,
	}
}

func (s TCPService) ToService() *db.Service {
	var status *string
	if s.Status != "" {
		status = &s.Status
	}

	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				if server.Address == "" {
					continue
				}
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
			if service.Name == "" {
				continue
			}
			services = append(services, dynamic.WRRService{
				Name:   service.Name,
				Weight: service.Weight,
			})
		}

		weighted = &WeightedRoundRobin{Services: services}
	}
	return &db.Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       status,
		Protocol:     "tcp",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
	}
}

func (s UDPService) ToService() *db.Service {
	var status *string
	if s.Status != "" {
		status = &s.Status
	}

	var lb *LoadBalancer
	var servers []Server
	if s.LoadBalancer != nil && s.LoadBalancer.Servers != nil {
		if s.LoadBalancer.Servers != nil {
			servers = make([]Server, len(s.LoadBalancer.Servers))
			for i, server := range s.LoadBalancer.Servers {
				if server.Address == "" {
					continue
				}
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
			if service.Name == "" {
				continue
			}
			services = append(services, dynamic.WRRService{
				Name:   service.Name,
				Weight: service.Weight,
			})
		}

		weighted = &WeightedRoundRobin{Services: services}
	}

	return &db.Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       status,
		Protocol:     "udp",
		ServerStatus: s.ServerStatus,
		LoadBalancer: lb,
		Weighted:     weighted,
	}
}

func getServices[T Serviceable](profile db.Profile, endpoint string) error {
	body, err := fetch(profile, endpoint)
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}
	defer body.Close()

	var serviceables []T
	if err := json.NewDecoder(body).Decode(&serviceables); err != nil {
		return fmt.Errorf("failed to decode services: %w", err)
	}

	// Current services
	dbServices, err := db.Query.ListServicesByProfileID(context.Background(), profile.ID)
	if err != nil {
		return fmt.Errorf("failed to list routers: %w", err)
	}

	existingServices := make(map[string]string) // name to ID mapping
	for _, dbService := range dbServices {
		existingServices[dbService.Name] = dbService.ID
	}

	var currentProtocol string // of T
	services := make(map[string]db.Service, len(serviceables))
	for _, s := range serviceables {
		newService := s.ToService()
		if newService.Name == "" {
			continue
		}
		services[newService.Name] = *newService
		currentProtocol = newService.Protocol

		data := db.UpsertServiceParams{
			ProfileID: profile.ID,
			Provider:  newService.Provider,
			Type:      newService.Type,
			Protocol:  newService.Protocol,
			Status:    newService.Status,
		}
		data.ID = existingServices[newService.Name]
		if data.ID == "" {
			data.ID = uuid.New().String()
			data.Name = newService.Name
		}

		data.ServerStatus, _ = json.Marshal(newService.ServerStatus)
		data.LoadBalancer, _ = json.Marshal(newService.LoadBalancer)
		data.Weighted, _ = json.Marshal(newService.Weighted)
		data.Mirroring, _ = json.Marshal(newService.Mirroring)
		data.Failover, _ = json.Marshal(newService.Failover)
		if _, err := db.Query.UpsertService(context.Background(), data); err != nil {
			slog.Error("Failed to upsert service", "error", err)
			continue
		}
	}

	// Cleanup if router doesn't exist locally (except our provider)
	for _, s := range dbServices {
		if s.Protocol != currentProtocol {
			continue
		}

		if _, ok := services[s.Name]; !ok && s.Provider != "http" {
			if err := db.Query.DeleteRouterByID(context.Background(), s.ID); err != nil {
				slog.Error("failed to delete service", "error", err)
				continue
			}
		}
	}

	return nil
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

		data, err := DecodeFromDB(profile.ID)
		if err != nil {
			slog.Error("Failed to decode config", "error", err)
			return
		}

		// Fetch routers
		getRouters[HTTPRouter](profile, HTTPRouterAPI)
		getRouters[TCPRouter](profile, TCPRouterAPI)
		getRouters[UDPRouter](profile, UDPRouterAPI)

		// Fetch services
		getServices[HTTPService](profile, HTTPServiceAPI)
		getServices[TCPService](profile, TCPServiceAPI)
		getServices[UDPService](profile, UDPServiceAPI)

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

		VerifyConfig(data)

		// Write to db
		if _, err := EncodeToDB(data); err != nil {
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
				switch existingItem := any(existing).(type) {
				case Router:
					if newRouter, ok := any(v).(Router); ok {
						newRouter.DNSProvider = existingItem.DNSProvider
						newRouter.ErrorState = existingItem.ErrorState
						merged[k] = any(newRouter).(T)
					}
				default:
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
