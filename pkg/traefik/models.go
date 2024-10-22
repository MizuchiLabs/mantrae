// Package traefik provides a client for the Traefik API
// Here are all the models used to convert between the API and the UI
package traefik

import (
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type Dynamic struct {
	ProfileID   int64                     `json:"profile_id,omitempty"`
	Overview    *Overview                 `json:"overview,omitempty"`
	Entrypoints []Entrypoint              `json:"entrypoints,omitempty"`
	Routers     map[string]Router         `json:"routers,omitempty"`
	Services    map[string]Service        `json:"services,omitempty"`
	Middlewares map[string]Middleware     `json:"middlewares,omitempty"`
	TLS         *dynamic.TLSConfiguration `json:"tls,omitempty"`
	Version     string                    `json:"version,omitempty"`
}

type Entrypoint struct {
	Name      string `json:"name,omitempty"`
	Address   string `json:"address,omitempty"`
	AsDefault bool   `json:"asDefault,omitempty"`
	HTTP      struct {
		Middlewares []string                 `json:"middlewares,omitempty"`
		TLS         *dynamic.RouterTLSConfig `json:"tls,omitempty"`
	} `json:"http,omitempty"`
}

type Router struct {
	// Common fields
	Name        string `json:"name,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Status      string `json:"status,omitempty"`
	RouterType  string `json:"routerType"`
	DNSProvider *int64 `json:"dnsProvider,omitempty"`
	SSLError    string `json:"sslError,omitempty"`
	AgentID     string `json:"agentID,omitempty"`

	Entrypoints []string                    `json:"entrypoints,omitempty"` // http, tcp, udp
	Middlewares []string                    `json:"middlewares,omitempty"` // http, tcp
	Rule        string                      `json:"rule,omitempty"`        // http, tcp
	RuleSyntax  string                      `json:"ruleSyntax,omitempty"`
	Service     string                      `json:"service,omitempty"` // http, tcp, udp
	Priority    int                         `json:"priority,omitempty"`
	TLS         *dynamic.RouterTCPTLSConfig `json:"tls,omitempty"` // Merge tcp and http
}

type Service struct {
	// Common fields
	Name         string            `json:"name,omitempty"`
	Provider     string            `json:"provider,omitempty"`
	Type         string            `json:"type,omitempty"`
	Status       string            `json:"status,omitempty"`
	ServiceType  string            `json:"serviceType,omitempty"` // "http" or "tcp" or "udp"
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
	AgentID      string            `json:"agentID,omitempty"`

	LoadBalancer *LoadBalancer       `json:"loadBalancer,omitempty"`
	Weighted     *WeightedRoundRobin `json:"weighted,omitempty"`

	// HTTP-specific fields
	Mirroring *dynamic.Mirroring `json:"mirroring,omitempty"`
	Failover  *dynamic.Failover  `json:"failover,omitempty"`
}

type LoadBalancer struct {
	Servers []Server `json:"servers,omitempty"`

	// HTTP-specific fields
	Sticky             *dynamic.Sticky             `json:"sticky,omitempty"`
	PassHostHeader     *bool                       `json:"passHostHeader,omitempty"`
	HealthCheck        *dynamic.ServerHealthCheck  `json:"healthCheck,omitempty"`
	ResponseForwarding *dynamic.ResponseForwarding `json:"responseForwarding,omitempty"`
	ServersTransport   string                      `json:"serversTransport,omitempty"`

	// TCP-specific fields
	ProxyProtocol    *dynamic.ProxyProtocol `json:"proxyProtocol,omitempty"`
	TerminationDelay *int                   `json:"terminationDelay,omitempty"`
}

type Server struct {
	URL     string `json:"url,omitempty"`
	Address string `json:"address,omitempty"`
}

type WeightedRoundRobin struct {
	Services []dynamic.WRRService `json:"services,omitempty"`

	// HTTP-specific fields
	Sticky      *dynamic.Sticky      `json:"sticky,omitempty"`
	HealthCheck *dynamic.HealthCheck `json:"healthCheck,omitempty"`
}

type Middleware struct {
	// Common fields
	Name           string `json:"name,omitempty"`
	Provider       string `json:"provider,omitempty"`
	Type           string `json:"type,omitempty"`
	Status         string `json:"status,omitempty"`
	MiddlewareType string `json:"middlewareType,omitempty"`
	AgentID        string `json:"agentID,omitempty"`

	// HTTP-specific fields
	AddPrefix         *dynamic.AddPrefix            `json:"addPrefix,omitempty"`
	StripPrefix       *dynamic.StripPrefix          `json:"stripPrefix,omitempty"`
	StripPrefixRegex  *dynamic.StripPrefixRegex     `json:"stripPrefixRegex,omitempty"`
	ReplacePath       *dynamic.ReplacePath          `json:"replacePath,omitempty"`
	ReplacePathRegex  *dynamic.ReplacePathRegex     `json:"replacePathRegex,omitempty"`
	Chain             *dynamic.Chain                `json:"chain,omitempty"`
	IPAllowList       *dynamic.IPAllowList          `json:"ipAllowList,omitempty"` // also for tcp
	Headers           *dynamic.Headers              `json:"headers,omitempty"`
	Errors            *dynamic.ErrorPage            `json:"errors,omitempty"`
	RateLimit         *dynamic.RateLimit            `json:"rateLimit,omitempty"`
	RedirectRegex     *dynamic.RedirectRegex        `json:"redirectRegex,omitempty"`
	RedirectScheme    *dynamic.RedirectScheme       `json:"redirectScheme,omitempty"`
	BasicAuth         *dynamic.BasicAuth            `json:"basicAuth,omitempty"`
	DigestAuth        *dynamic.DigestAuth           `json:"digestAuth,omitempty"`
	ForwardAuth       *dynamic.ForwardAuth          `json:"forwardAuth,omitempty"`
	InFlightReq       *dynamic.InFlightReq          `json:"inFlightReq,omitempty"`
	Buffering         *dynamic.Buffering            `json:"buffering,omitempty"`
	CircuitBreaker    *dynamic.CircuitBreaker       `json:"circuitBreaker,omitempty"`
	Compress          *dynamic.Compress             `json:"compress,omitempty"`
	PassTLSClientCert *dynamic.PassTLSClientCert    `json:"passTLSClientCert,omitempty"`
	Retry             *dynamic.Retry                `json:"retry,omitempty"`
	GrpcWeb           *dynamic.GrpcWeb              `json:"grpcWeb,omitempty"`
	Plugin            map[string]dynamic.PluginConf `json:"plugin,omitempty"`

	// TCP-specific fields
	InFlightConn *dynamic.TCPInFlightConn `json:"inFlightConn,omitempty"`
}

type Overview struct {
	HTTP     HTTPOverview `json:"http,omitempty"`
	TCP      TCPOverview  `json:"tcp,omitempty"`
	UDP      UDPOverview  `json:"udp,omitempty"`
	Features struct {
		Tracing   string `json:"tracing,omitempty"`
		Metrics   string `json:"metrics,omitempty"`
		AccessLog bool   `json:"accessLog,omitempty"`
	} `json:"features,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

type BasicOverview struct {
	Total    int `json:"total,omitempty"`
	Warnings int `json:"warnings,omitempty"`
	Errors   int `json:"errors,omitempty"`
}

type HTTPOverview struct {
	Routers    BasicOverview `json:"routers,omitempty"`
	Services   BasicOverview `json:"services,omitempty"`
	Middleware BasicOverview `json:"middlewares,omitempty"`
}

type TCPOverview struct {
	Routers    BasicOverview `json:"routers,omitempty"`
	Services   BasicOverview `json:"services,omitempty"`
	Middleware BasicOverview `json:"middlewares,omitempty"`
}

type UDPOverview struct {
	Routers  BasicOverview `json:"routers,omitempty"`
	Services BasicOverview `json:"services,omitempty"`
}
