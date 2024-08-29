package traefik

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/traefik/genconf/dynamic"
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
	EntrypointsAPI     = "/api/entrypoints"
	VersionAPI         = "/api/version"
)

type BaseFields struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Status   string `json:"status,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type HTTPRouter struct {
	BaseFields
	RouterType  string                   `json:"routerType,omitempty"`
	Entrypoints []string                 `json:"entrypoints,omitempty"`
	Middlewares []string                 `json:"middlewares,omitempty"`
	Rule        string                   `json:"rule,omitempty"`
	Service     string                   `json:"service,omitempty"`
	Priority    *big.Int                 `json:"priority,omitempty"`
	TLS         *dynamic.RouterTLSConfig `json:"tls,omitempty"`
}

type TCPRouter struct {
	BaseFields
	RouterType  string                      `json:"routerType,omitempty"`
	Entrypoints []string                    `json:"entrypoints,omitempty"`
	Middlewares []string                    `json:"middlewares,omitempty"`
	Rule        string                      `json:"rule,omitempty"`
	Service     string                      `json:"service,omitempty"`
	Priority    *big.Int                    `json:"priority,omitempty"`
	TLS         *dynamic.RouterTCPTLSConfig `json:"tls,omitempty"`
}

type UDPRouter struct {
	BaseFields
	RouterType  string   `json:"routerType,omitempty"`
	Entrypoints []string `json:"entrypoints,omitempty"`
	Service     string   `json:"service,omitempty"`
}

type Routerable interface {
	ToRouter() *Router
}

func (r HTTPRouter) ToRouter() *Router {
	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "http",
		Entrypoints: r.Entrypoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		Service:     r.Service,
		// Priority:    r.Priority,
		TLS: r.TLS,
	}
}

func (r TCPRouter) ToRouter() *Router {
	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "tcp",
		Entrypoints: r.Entrypoints,
		Middlewares: r.Middlewares,
		Rule:        r.Rule,
		Service:     r.Service,
		// Priority:    r.Priority,
		TCPTLS: r.TLS,
	}
}

func (r UDPRouter) ToRouter() *Router {
	return &Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "udp",
		Entrypoints: r.Entrypoints,
		Service:     r.Service,
	}
}

func getRouters[T Routerable](p Profile, endpoint string) map[string]Router {
	body, err := p.fetch(endpoint)
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
		routers[newRouter.Name] = *newRouter
	}
	return routers
}

type HTTPService struct {
	BaseFields
	ServiceType  string                       `json:"serviceType,omitempty"`
	ServerStatus map[string]string            `json:"serverStatus,omitempty"`
	LoadBalancer *dynamic.ServersLoadBalancer `json:"loadBalancer,omitempty"`
	Weighted     *dynamic.WeightedRoundRobin  `json:"weighted,omitempty"`
	Mirroring    *dynamic.Mirroring           `json:"mirroring,omitempty"`
	Failover     *dynamic.Failover            `json:"failover,omitempty"`
}

type TCPService struct {
	BaseFields
	ServiceType  string                          `json:"serviceType,omitempty"`
	ServerStatus map[string]string               `json:"serverStatus,omitempty"`
	LoadBalancer *dynamic.TCPServersLoadBalancer `json:"loadBalancer,omitempty"`
	Weighted     *dynamic.TCPWeightedRoundRobin  `json:"weighted,omitempty"`
}

type UDPService struct {
	BaseFields
	ServiceType  string                          `json:"serviceType,omitempty"`
	ServerStatus map[string]string               `json:"serverStatus,omitempty"`
	LoadBalancer *dynamic.UDPServersLoadBalancer `json:"loadBalancer,omitempty"`
	Weighted     *dynamic.UDPWeightedRoundRobin  `json:"weighted,omitempty"`
}

type Serviceable interface {
	ToService() *Service
}

func (s HTTPService) ToService() *Service {
	return &Service{
		Name:         s.Name,
		Provider:     s.Provider,
		Type:         s.Type,
		Status:       s.Status,
		ServiceType:  "http",
		ServerStatus: s.ServerStatus,
		LoadBalancer: s.LoadBalancer,
		Weighted:     s.Weighted,
		Mirroring:    s.Mirroring,
		Failover:     s.Failover,
	}
}

func (s TCPService) ToService() *Service {
	return &Service{
		Name:            s.Name,
		Provider:        s.Provider,
		Type:            s.Type,
		Status:          s.Status,
		ServiceType:     "tcp",
		ServerStatus:    s.ServerStatus,
		TCPLoadBalancer: s.LoadBalancer,
		TCPWeighted:     s.Weighted,
	}
}

func (s UDPService) ToService() *Service {
	return &Service{
		Name:            s.Name,
		Provider:        s.Provider,
		Type:            s.Type,
		Status:          s.Status,
		ServiceType:     "udp",
		ServerStatus:    s.ServerStatus,
		UDPLoadBalancer: s.LoadBalancer,
		UDPWeighted:     s.Weighted,
	}
}

func getServices[T Serviceable](p Profile, endpoint string) map[string]Service {
	body, err := p.fetch(endpoint)
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
		services[newService.Name] = *newService
	}

	return services
}

type HTTPMiddleware struct {
	BaseFields
	MiddlewareType    string                     `json:"middlewareType,omitempty"`
	AddPrefix         *dynamic.AddPrefix         `json:"addPrefix,omitempty"`
	StripPrefix       *dynamic.StripPrefix       `json:"stripPrefix,omitempty"`
	StripPrefixRegex  *dynamic.StripPrefixRegex  `json:"stripPrefixRegex,omitempty"`
	ReplacePath       *dynamic.ReplacePath       `json:"replacePath,omitempty"`
	ReplacePathRegex  *dynamic.ReplacePathRegex  `json:"replacePathRegex,omitempty"`
	Chain             *dynamic.Chain             `json:"chain,omitempty"`
	IPAllowList       *dynamic.IPAllowList       `json:"ipAllowList,omitempty"`
	Headers           *dynamic.Headers           `json:"headers,omitempty"`
	Errors            *dynamic.ErrorPage         `json:"errors,omitempty"`
	RateLimit         *dynamic.RateLimit         `json:"rateLimit,omitempty"`
	RedirectRegex     *dynamic.RedirectRegex     `json:"redirectRegex,omitempty"`
	RedirectScheme    *dynamic.RedirectScheme    `json:"redirectScheme,omitempty"`
	BasicAuth         *dynamic.BasicAuth         `json:"basicAuth,omitempty"`
	DigestAuth        *dynamic.DigestAuth        `json:"digestAuth,omitempty"`
	ForwardAuth       *dynamic.ForwardAuth       `json:"forwardAuth,omitempty"`
	InFlightReq       *dynamic.InFlightReq       `json:"inFlightReq,omitempty"`
	Buffering         *dynamic.Buffering         `json:"buffering,omitempty"`
	CircuitBreaker    *dynamic.CircuitBreaker    `json:"circuitBreaker,omitempty"`
	Compress          *dynamic.Compress          `json:"compress,omitempty"`
	PassTLSClientCert *dynamic.PassTLSClientCert `json:"passTLSClientCert,omitempty"`
	Retry             *dynamic.Retry             `json:"retry,omitempty"`
}

type TCPMiddleware struct {
	BaseFields
	MiddlewareType string                   `json:"middlewareType,omitempty"`
	InFlightConn   *dynamic.TCPInFlightConn `json:"inFlightConn,omitempty"`
	IPAllowList    *dynamic.TCPIPAllowList  `json:"ipAllowList,omitempty"`
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
	}
}

func (m TCPMiddleware) ToMiddleware() *Middleware {
	return &Middleware{
		Name:           m.Name,
		Provider:       m.Provider,
		Type:           m.Type,
		Status:         m.Status,
		MiddlewareType: "tcp",
		InFlightConn:   m.InFlightConn,
		TCPIPAllowList: m.IPAllowList,
	}
}

func getMiddlewares[T Middlewareable](p Profile, endpoint string) map[string]Middleware {
	body, err := p.fetch(endpoint)
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
		middlewares[newMiddleware.Name] = *newMiddleware
	}

	return middlewares
}

func GetTraefikConfig() {
	for i, profile := range ProfileData.Profiles {
		if profile.URL == "" {
			continue
		}

		d := Dynamic{
			Entrypoints: make([]Entrypoint, 0),
			Routers:     make(map[string]Router),
			Services:    make(map[string]Service),
			Middlewares: make(map[string]Middleware),
		}

		// Retrieve routers
		d.Routers = merge(
			getRouters[HTTPRouter](profile, HTTPRouterAPI),
			getRouters[TCPRouter](profile, TCPRouterAPI),
			getRouters[UDPRouter](profile, UDPRouterAPI),
			filterByLocalProvider(
				profile.Dynamic.Routers,
				func(r Router) string { return r.Provider },
			),
		)

		// Retrieve services
		d.Services = merge(
			getServices[HTTPService](profile, HTTPServiceAPI),
			getServices[TCPService](profile, TCPServiceAPI),
			getServices[UDPService](profile, UDPServiceAPI),
			filterByLocalProvider(
				profile.Dynamic.Services,
				func(s Service) string { return s.Provider },
			),
		)

		// Fetch middlewares
		d.Middlewares = merge(
			getMiddlewares[HTTPMiddleware](profile, HTTPMiddlewaresAPI),
			getMiddlewares[TCPMiddleware](profile, TCPMiddlewaresAPI),
			filterByLocalProvider(
				profile.Dynamic.Middlewares,
				func(m Middleware) string { return m.Provider },
			),
		)

		// Retrieve entrypoints
		entrypoints, err := profile.fetch(EntrypointsAPI)
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
		version, err := profile.fetch(VersionAPI)
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

		profile.Dynamic = d
		ProfileData.Profiles[i] = profile
	}

	if err := ProfileData.Save(); err != nil {
		slog.Error("Failed to save profiles", "error", err)
	}
}

// Sync periodically syncs the Traefik configuration
func Sync() {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for range ticker.C {
		GetTraefikConfig()
	}
}

// Filter http provider
func filterByLocalProvider[T any](items map[string]T, getProvider func(T) string) map[string]T {
	filteredItems := make(map[string]T)
	for key, item := range items {
		if getProvider(item) == "http" {
			filteredItems[key] = item
		}
	}
	return filteredItems
}

func merge[T any](maps ...map[string]T) map[string]T {
	merged := make(map[string]T)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func (p Profile) fetch(endpoint string) (io.ReadCloser, error) {
	if p.URL == "" || endpoint == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	apiURL := p.URL + endpoint
	client := http.Client{Timeout: time.Second * 10}
	if !p.TLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.Username != "" && p.Password != "" {
		req.SetBasicAuth(p.Username, p.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", apiURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
