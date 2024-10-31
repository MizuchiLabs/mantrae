// Package traefik provides a client for the Traefik API
// Here are all the models used to convert between the API and the UI
package traefik

import (
	"sync"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type Dynamic struct {
	ProfileID   int64        `json:"profile_id,omitempty"`
	Overview    *Overview    `json:"overview,omitempty"`
	Entrypoints []Entrypoint `json:"entrypoints,omitempty"`
	Version     string       `json:"version,omitempty"`
	Mutex       sync.Mutex   `json:"-"`
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

type Server struct {
	URL     string `json:"url,omitempty"`
	Address string `json:"address,omitempty"`
}

type ErrorState struct {
	SSL string `json:"ssl,omitempty"`
	DNS string `json:"dns,omitempty"`
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

type Plugin struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	Author        string        `json:"author"`
	Type          string        `json:"type"`
	Import        string        `json:"import"`
	Summary       string        `json:"summary"`
	IconUrl       string        `json:"iconUrl"`
	BannerUrl     string        `json:"bannerUrl"`
	Readme        string        `json:"readme"`
	LatestVersion string        `json:"latestVersion"`
	Versions      []string      `json:"versions"`
	Stars         int           `json:"stars"`
	Snippet       PluginSnippet `json:"snippet"`
	CreatedAt     string        `json:"createdAt"`
}

type PluginSnippet struct {
	Yaml string `json:"yaml"`
}
