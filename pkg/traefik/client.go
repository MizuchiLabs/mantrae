package traefik

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
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
	var dbRouter *db.Router
	dbBytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(dbBytes, &dbRouter); err != nil {
		return nil
	}

	dbRouter.Protocol = "http"
	return dbRouter
}

func (r TCPRouter) ToRouter() *db.Router {
	var dbRouter *db.Router
	dbBytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(dbBytes, &dbRouter); err != nil {
		return nil
	}

	dbRouter.Protocol = "tcp"
	return dbRouter
}

func (r UDPRouter) ToRouter() *db.Router {
	var dbRouter *db.Router
	dbBytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(dbBytes, &dbRouter); err != nil {
		return nil
	}

	dbRouter.Protocol = "udp"
	return dbRouter
}

func getRouters[T Routerable](profile db.Profile, endpoint string) error {
	typeName := reflect.TypeOf((*T)(nil)).Elem().Name() // Get the name of the type T
	var protocol string
	switch typeName {
	case "HTTPRouter":
		protocol = "http"
	case "TCPRouter":
		protocol = "tcp"
	case "UDPRouter":
		protocol = "udp"
	}

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

	routers := make(map[string]bool, len(routerables))
	for _, r := range routerables {
		newRouter := r.ToRouter()
		if newRouter.Name == "" || newRouter.Provider == "http" {
			continue
		}
		routers[newRouter.Name] = true

		data := db.UpsertRouterParams{
			ID:         uuid.New().String(),
			ProfileID:  profile.ID,
			Name:       newRouter.Name,
			Provider:   newRouter.Provider,
			Protocol:   newRouter.Protocol,
			Status:     newRouter.Status,
			Rule:       newRouter.Rule,
			RuleSyntax: newRouter.RuleSyntax,
			Service:    newRouter.Service,
			Priority:   newRouter.Priority,
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
		if r.Protocol != protocol || r.Provider == "http" {
			continue
		}

		if _, ok := routers[r.Name]; !ok {
			slog.Info("Removing router", "name", r.Name, "id", r.ID)
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
	var dbService *db.Service
	sBytes, err := json.Marshal(s)
	if err != nil {
		slog.Error("Failed to marshal service", "error", err)
		return nil
	}

	if err := json.Unmarshal(sBytes, &dbService); err != nil {
		slog.Error("Failed to unmarshal service", "error", err)
		return nil
	}
	dbService.Protocol = "http"
	return dbService
}

func (s TCPService) ToService() *db.Service {
	var dbService *db.Service
	sBytes, err := json.Marshal(s)
	if err != nil {
		slog.Error("Failed to marshal service", "error", err)
		return nil
	}

	if err := json.Unmarshal(sBytes, &dbService); err != nil {
		slog.Error("Failed to unmarshal service", "error", err)
		return nil
	}

	dbService.Protocol = "tcp"
	return dbService
}

func (s UDPService) ToService() *db.Service {
	var dbService *db.Service
	sBytes, err := json.Marshal(s)
	if err != nil {
		slog.Error("Failed to marshal service", "error", err)
		return nil
	}

	if err := json.Unmarshal(sBytes, &dbService); err != nil {
		slog.Error("Failed to unmarshal service", "error", err)
		return nil
	}

	dbService.Protocol = "udp"
	return dbService
}

func getServices[T Serviceable](profile db.Profile, endpoint string) error {
	typeName := reflect.TypeOf((*T)(nil)).Elem().Name() // Get the name of the type T
	var protocol string
	switch typeName {
	case "HTTPService":
		protocol = "http"
	case "TCPService":
		protocol = "tcp"
	case "UDPService":
		protocol = "udp"
	}

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

	services := make(map[string]db.Service, len(serviceables))
	for _, s := range serviceables {
		newService := s.ToService()
		if newService.Name == "" {
			continue
		}
		services[newService.Name] = *newService

		_, err := db.Query.GetServiceByName(context.Background(), db.GetServiceByNameParams{
			ProfileID: profile.ID,
			Name:      newService.Name,
		})
		if newService.Provider == "http" && err == nil {
			data := db.UpsertServiceParams{
				ID:        uuid.New().String(),
				ProfileID: profile.ID,
				Name:      newService.Name,
			}
			data.ServerStatus, _ = json.Marshal(newService.ServerStatus)
			if _, err := db.Query.UpsertService(context.Background(), data); err != nil {
				slog.Error("Failed to upsert service", "error", err)
				continue
			}
		}

		if newService.Provider != "http" {
			data := db.UpsertServiceParams{
				ID:        uuid.New().String(),
				ProfileID: profile.ID,
				Name:      newService.Name,
				Provider:  newService.Provider,
				Type:      newService.Type,
				Protocol:  newService.Protocol,
				Status:    newService.Status,
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
	}

	// Cleanup if router doesn't exist locally (except our provider)
	for _, s := range dbServices {
		if s.Protocol != protocol || s.Provider == "http" {
			continue
		}

		if _, ok := services[s.Name]; !ok {
			slog.Info("Removing service", "name", s.Name)
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
		if err := getRouters[HTTPRouter](profile, HTTPRouterAPI); err != nil {
			slog.Error("Failed to get routers", "error", err)
		}
		if err := getRouters[TCPRouter](profile, TCPRouterAPI); err != nil {
			slog.Error("Failed to get routers", "error", err)
		}
		if err := getRouters[UDPRouter](profile, UDPRouterAPI); err != nil {
			slog.Error("Failed to get routers", "error", err)
		}

		// Fetch services
		if err := getServices[HTTPService](profile, HTTPServiceAPI); err != nil {
			slog.Error("Failed to get services", "error", err)
		}
		if err := getServices[TCPService](profile, TCPServiceAPI); err != nil {
			slog.Error("Failed to get services", "error", err)
		}
		if err := getServices[UDPService](profile, UDPServiceAPI); err != nil {
			slog.Error("Failed to get services", "error", err)
		}

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
