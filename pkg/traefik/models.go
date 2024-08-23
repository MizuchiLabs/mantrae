// Package traefik provides a client for the Traefik API
// Here are all the models used to convert between the API and the UI
package traefik

import (
	"sync"

	"github.com/traefik/genconf/dynamic"
)

type Profiles struct {
	Profiles map[string]Profile `json:"profiles,omitempty"`
	mu       sync.RWMutex
}

type Profile struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Username string  `json:"username,omitempty"`
	Password string  `json:"password,omitempty"`
	Dynamic  Dynamic `json:"dynamic,omitempty"`
}

type Dynamic struct {
	Entrypoints []Entrypoint          `json:"entrypoints,omitempty"`
	Routers     map[string]Router     `json:"routers,omitempty"`
	Services    map[string]Service    `json:"services,omitempty"`
	Middlewares map[string]Middleware `json:"middlewares,omitempty"`
	Version     string                `json:"version,omitempty"`
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
	Name       string `json:"name,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Status     string `json:"status,omitempty"`
	RouterType string `json:"routerType,omitempty"`

	Entrypoints []string `json:"entrypoints,omitempty"` // http, tcp, udp
	Middlewares []string `json:"middlewares,omitempty"` // http, tcp
	Rule        string   `json:"rule,omitempty"`        // http, tcp
	Service     string   `json:"service,omitempty"`     // http, tcp, udp
	// Priority    *big.Int                    `json:"priority,omitempty"`
	TLS    *dynamic.RouterTLSConfig    `json:"tls,omitempty"`    // http
	TCPTLS *dynamic.RouterTCPTLSConfig `json:"tcpTLS,omitempty"` // tcp
}

type Service struct {
	// Common fields
	Name         string            `json:"name,omitempty"`
	Provider     string            `json:"provider,omitempty"`
	Type         string            `json:"type,omitempty"`
	Status       string            `json:"status,omitempty"`
	ServiceType  string            `json:"serviceType,omitempty"` // "http" or "tcp" or "udp"
	ServerStatus map[string]string `json:"serverStatus,omitempty"`

	// HTTP-specific fields
	LoadBalancer *dynamic.ServersLoadBalancer `json:"loadBalancer,omitempty"`
	Weighted     *dynamic.WeightedRoundRobin  `json:"weighted,omitempty"`
	Mirroring    *dynamic.Mirroring           `json:"mirroring,omitempty"`
	Failover     *dynamic.Failover            `json:"failover,omitempty"`

	// TCP-specific fields
	TCPLoadBalancer *dynamic.TCPServersLoadBalancer `json:"tcpLoadBalancer,omitempty"`
	TCPWeighted     *dynamic.TCPWeightedRoundRobin  `json:"tcpWeighted,omitempty"`

	// UDP-specific fields
	UDPLoadBalancer *dynamic.UDPServersLoadBalancer `json:"udpLoadBalancer,omitempty"`
	UDPWeighted     *dynamic.UDPWeightedRoundRobin  `json:"udpWeighted,omitempty"`
}

type Middleware struct {
	// Common fields
	Name           string `json:"name,omitempty"`
	Provider       string `json:"provider,omitempty"`
	Type           string `json:"type,omitempty"`
	Status         string `json:"status,omitempty"`
	MiddlewareType string `json:"middlewareType,omitempty"` // "http" or "tcp"

	// HTTP-specific fields
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

	// TCP-specific fields
	InFlightConn   *dynamic.TCPInFlightConn `json:"inFlightConn,omitempty"`
	TCPIPAllowList *dynamic.TCPIPAllowList  `json:"tcpIpAllowList,omitempty"`
}
