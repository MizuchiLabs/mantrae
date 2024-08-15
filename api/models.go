package api

import (
	"math/big"

	"github.com/traefik/genconf/dynamic"
)

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
	RouterType  string                   `json:"routerType,omitempty"`
	Entrypoints []string                 `json:"entrypoints,omitempty"`
	Middlewares []string                 `json:"middlewares,omitempty"`
	Rule        string                   `json:"rule,omitempty"`
	Service     string                   `json:"service,omitempty"`
	Priority    big.Int                  `json:"priority,omitempty"`
	TLS         *dynamic.RouterTLSConfig `json:"tls,omitempty"`
	Status      string                   `json:"status,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Provider    string                   `json:"provider,omitempty"`
}

type Service struct {
	ServiceType  string                       `json:"serviceType,omitempty"`
	LoadBalancer *dynamic.ServersLoadBalancer `json:"loadBalancer,omitempty"`
	Weighted     *dynamic.WeightedRoundRobin  `json:"weighted,omitempty"`
	Mirroring    *dynamic.Mirroring           `json:"mirroring,omitempty"`
	Failover     *dynamic.Failover            `json:"failover,omitempty"`
	ServerStatus map[string]string            `json:"serverStatus,omitempty"`
	Status       string                       `json:"status,omitempty"`
	Name         string                       `json:"name,omitempty"`
	Provider     string                       `json:"provider,omitempty"`
	Type         string                       `json:"type,omitempty"`
}

type HTTPMiddleware struct {
	AddPrefix         *dynamic.AddPrefix            `json:"addPrefix,omitempty"`
	StripPrefix       *dynamic.StripPrefix          `json:"stripPrefix,omitempty"`
	StripPrefixRegex  *dynamic.StripPrefixRegex     `json:"stripPrefixRegex,omitempty"`
	ReplacePath       *dynamic.ReplacePath          `json:"replacePath,omitempty"`
	ReplacePathRegex  *dynamic.ReplacePathRegex     `json:"replacePathRegex,omitempty"`
	Chain             *dynamic.Chain                `json:"chain,omitempty"`
	IPWhiteList       *dynamic.IPWhiteList          `json:"ipWhiteList,omitempty"`
	IPAllowList       *dynamic.IPAllowList          `json:"ipAllowList,omitempty"`
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
	ContentType       *dynamic.ContentType          `json:"contentType,omitempty"`
	Plugin            map[string]dynamic.PluginConf `json:"plugin,omitempty"`
	Status            string                        `json:"status,omitempty"`
	Name              string                        `json:"name,omitempty"`
	Provider          string                        `json:"provider,omitempty"`
	Type              string                        `json:"type,omitempty"`
}

type TCPMiddleware struct {
	InFlightConn *dynamic.TCPInFlightConn `json:"inFlightConn,omitempty"`
	IPWhiteList  *dynamic.TCPIPWhiteList  `json:"ipWhiteList,omitempty"`
	IPAllowList  *dynamic.TCPIPAllowList  `json:"ipAllowList,omitempty"`
	Status       string                   `json:"status,omitempty"`
	Name         string                   `json:"name,omitempty"`
	Provider     string                   `json:"provider,omitempty"`
	Type         string                   `json:"type,omitempty"`
}

type Dynamic struct {
	Entrypoints     []Entrypoint     `json:"entrypoints,omitempty"`
	Routers         []Router         `json:"routers,omitempty"`
	Services        []Service        `json:"services,omitempty"`
	HTTPMiddlewares []HTTPMiddleware `json:"httpmiddlewares,omitempty"`
	TCPMiddlewares  []TCPMiddleware  `json:"tcpmiddlewares,omitempty"`
	Version         string           `json:"version,omitempty"`
}

type Instance struct {
	URL      string  `json:"url"`
	Username string  `json:"username,omitempty"`
	Password string  `json:"password,omitempty"`
	Dynamic  Dynamic `json:"dynamic,omitempty"`
}

type Profile struct {
	Name     string   `json:"name"`
	Instance Instance `json:"instance,omitempty"`
}
