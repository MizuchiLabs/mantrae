package schema

import (
	"github.com/traefik/traefik/v3/pkg/config/runtime"
	"github.com/traefik/traefik/v3/pkg/types"
)

// EntryPoint -----------------------------------------------------------------
type EntryPoint struct {
	Name            string `json:"name,omitempty"`
	Address         string `json:"address,omitempty"`
	AllowACMEByPass bool   `json:"allow_acme_by_pass,omitempty"`
	ReusePort       bool   `json:"reuse_port,omitempty"`
	AsDefault       bool   `json:"as_default,omitempty"`
	// Transport        *EntryPointsTransport `json:"transport,omitempty"`
	ProxyProtocol    *ProxyProtocol    `json:"proxy_protocol,omitempty"`
	ForwardedHeaders *ForwardedHeaders `json:"forwarded_headers,omitempty"`
	HTTP             HTTPConfig        `json:"http"`
	HTTP2            *HTTP2Config      `json:"http_2,omitempty"`
	HTTP3            *HTTP3Config      `json:"http_3,omitempty"`
	// UDP              *UDPConfig           `json:"udp,omitempty"`
	Observability *ObservabilityConfig `json:"observability,omitempty"`
}

type TLSConfig struct {
	Options      string         `json:"options,omitempty"`
	CertResolver string         `json:"cert_resolver,omitempty"`
	Domains      []types.Domain `json:"domains,omitempty"`
}

type ForwardedHeaders struct {
	Insecure   bool     `json:"insecure,omitempty"`
	TrustedIPs []string `json:"trusted_i_ps,omitempty"`
	Connection []string `json:"connection,omitempty"`
}

type HTTPConfig struct {
	Redirections          *Redirections `json:"redirections,omitempty"`
	Middlewares           []string      `json:"middlewares,omitempty"`
	TLS                   *TLSConfig    `json:"tls,omitempty"`
	EncodeQuerySemicolons bool          `json:"encode_query_semicolons,omitempty"`
	SanitizePath          *bool         `json:"sanitize_path,omitempty"`
	MaxHeaderBytes        int           `json:"max_header_bytes,omitempty"`
}

type HTTP2Config struct {
	MaxConcurrentStreams int32 `json:"max_concurrent_streams,omitempty"`
}

type HTTP3Config struct {
	AdvertisedPort int `json:"advertised_port,omitempty"`
}

// type UDPConfig struct {
// 	Timeout ptypes.Duration `json:"timeout,omitempty"`
// }

type ObservabilityConfig struct {
	AccessLogs *bool `json:"access_logs,omitempty"`
	Tracing    *bool `json:"tracing,omitempty"`
	Metrics    *bool `json:"metrics,omitempty"`
}

type Redirections struct {
	EntryPoint *RedirectEntryPoint `json:"entry_point,omitempty"`
}

type RedirectEntryPoint struct {
	To        string `json:"to,omitempty"`
	Scheme    string `json:"scheme,omitempty"`
	Permanent bool   `json:"permanent,omitempty"`
	Priority  int    `json:"priority,omitempty"`
}

type ProxyProtocol struct {
	Insecure   bool     `json:"insecure,omitempty"`
	TrustedIPs []string `json:"trusted_i_ps,omitempty"`
}

// type EntryPointsTransport struct {
// 	// LifeCycle            *LifeCycle
// 	// RespondingTimeouts   *RespondingTimeouts
// 	KeepAliveMaxTime     ptypes.Duration
// 	KeepAliveMaxRequests int
// }

// Overview -------------------------------------------------------------------
type Overview struct {
	HTTP      SchemeOverview `json:"http"`
	TCP       SchemeOverview `json:"tcp"`
	UDP       SchemeOverview `json:"udp"`
	Features  Features       `json:"features"`
	Providers []string       `json:"providers,omitempty"`
}
type SchemeOverview struct {
	Routers     *Section `json:"routers,omitempty"`
	Services    *Section `json:"services,omitempty"`
	Middlewares *Section `json:"middlewares,omitempty"`
}
type Section struct {
	Total    int `json:"total"`
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
}
type Features struct {
	Tracing   string `json:"tracing"`
	Metrics   string `json:"metrics"`
	AccessLog bool   `json:"accessLog"`
}

// Configuration --------------------------------------------------------------
type Configuration struct {
	Routers        map[string]*runtime.RouterInfo        `json:"routers,omitempty"`
	Middlewares    map[string]*runtime.MiddlewareInfo    `json:"middlewares,omitempty"`
	Services       map[string]*serviceInfoRepresentation `json:"services,omitempty"`
	TCPRouters     map[string]*runtime.TCPRouterInfo     `json:"tcpRouters,omitempty"`
	TCPMiddlewares map[string]*runtime.TCPMiddlewareInfo `json:"tcpMiddlewares,omitempty"`
	TCPServices    map[string]*runtime.TCPServiceInfo    `json:"tcpServices,omitempty"`
	UDPRouters     map[string]*runtime.UDPRouterInfo     `json:"udpRouters,omitempty"`
	UDPServices    map[string]*runtime.UDPServiceInfo    `json:"udpServices,omitempty"`
}
type serviceInfoRepresentation struct {
	*runtime.ServiceInfo
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
}
