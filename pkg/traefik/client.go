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
	ToRouter() Router
}

func (r HTTPRouter) ToRouter() Router {
	return Router{
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

func (r TCPRouter) ToRouter() Router {
	return Router{
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

func (r UDPRouter) ToRouter() Router {
	return Router{
		Name:        r.Name,
		Provider:    r.Provider,
		Status:      r.Status,
		RouterType:  "udp",
		Entrypoints: r.Entrypoints,
		Service:     r.Service,
	}
}

func getRouters[T Routerable](c Client, endpoint string) []Router {
	body, err := c.fetch(endpoint)
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

	var routers []Router
	for _, r := range routerables {
		routers = append(routers, r.ToRouter())
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
	ToService() Service
}

func (s HTTPService) ToService() Service {
	return Service{
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

func (s TCPService) ToService() Service {
	return Service{
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

func (s UDPService) ToService() Service {
	return Service{
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

func getServices[T Serviceable](c Client, endpoint string) []Service {
	body, err := c.fetch(endpoint)
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

	var services []Service
	for _, s := range serviceables {
		services = append(services, s.ToService())
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
	ContentType       *dynamic.ContentType       `json:"contentType,omitempty"`
}

type TCPMiddleware struct {
	BaseFields
	MiddlewareType string                   `json:"middlewareType,omitempty"`
	InFlightConn   *dynamic.TCPInFlightConn `json:"inFlightConn,omitempty"`
	IPAllowList    *dynamic.TCPIPAllowList  `json:"ipAllowList,omitempty"`
}

type Middlewareable interface {
	ToMiddleware() Middleware
}

func (m HTTPMiddleware) ToMiddleware() Middleware {
	return Middleware{
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
		ContentType:       m.ContentType,
	}
}

func (m TCPMiddleware) ToMiddleware() Middleware {
	return Middleware{
		Name:           m.Name,
		Provider:       m.Provider,
		Type:           m.Type,
		Status:         m.Status,
		MiddlewareType: "tcp",
		InFlightConn:   m.InFlightConn,
		TCPIPAllowList: m.IPAllowList,
	}
}

func getMiddlewares[T Middlewareable](c Client, endpoint string) []Middleware {
	body, err := c.fetch(endpoint)
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

	var middlewares []Middleware
	for _, m := range middlewareables {
		middlewares = append(middlewares, m.ToMiddleware())
	}

	return middlewares
}

func GetTraefikConfig() {
	var p Profiles
	if err := p.Load(); err != nil {
		slog.Error("Failed to load profiles", "error", err)
		return
	}

	for _, profile := range p.Profiles {
		if profile.Client.URL == "" {
			continue
		}

		c := profile.Client
		d := Dynamic{
			Routers:     make(map[string]Router),
			Services:    make(map[string]Service),
			Middlewares: make(map[string]Middleware),
		}

		// Retrieve routers
		var tRouter []Router
		tRouter = append(tRouter, getRouters[HTTPRouter](c, HTTPRouterAPI)...)
		tRouter = append(tRouter, getRouters[TCPRouter](c, TCPRouterAPI)...)
		tRouter = append(tRouter, getRouters[UDPRouter](c, UDPRouterAPI)...)
		for _, r := range tRouter {
			d.Routers[r.Name] = r
		}
		for _, r := range profile.Client.Dynamic.Routers {
			d.Routers[r.Name] = r
		}

		// Retrieve services
		var tServices []Service
		tServices = append(tServices, getServices[HTTPService](c, HTTPServiceAPI)...)
		tServices = append(tServices, getServices[TCPService](c, TCPServiceAPI)...)
		tServices = append(tServices, getServices[UDPService](c, UDPServiceAPI)...)
		for _, s := range tServices {
			d.Services[s.Name] = s
		}
		for _, s := range profile.Client.Dynamic.Services {
			d.Services[s.Name] = s
		}

		// Fetch middlewares
		var tMiddlewares []Middleware
		tMiddlewares = append(
			tMiddlewares,
			getMiddlewares[HTTPMiddleware](c, HTTPMiddlewaresAPI)...)
		tMiddlewares = append(
			tMiddlewares,
			getMiddlewares[TCPMiddleware](c, TCPMiddlewaresAPI)...)
		for _, m := range tMiddlewares {
			d.Middlewares[m.Name] = m
		}
		for _, m := range profile.Client.Dynamic.Middlewares {
			d.Middlewares[m.Name] = m
		}

		// Retrieve entrypoints
		entrypoints, err := c.fetch(EntrypointsAPI)
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
		version, err := c.fetch(VersionAPI)
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

		profile.Client.Dynamic = d
	}

	if err := p.Save(); err != nil {
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

func (c Client) fetch(endpoint string) (io.ReadCloser, error) {
	if c.URL == "" || endpoint == "" {
		return nil, fmt.Errorf("invalid URL or endpoint")
	}

	apiURL := c.URL + endpoint
	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
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
