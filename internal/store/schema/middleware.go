package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	HTTPMiddleware dynamic.Middleware
	TCPMiddleware  dynamic.TCPMiddleware
)

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
	Stars         int64         `json:"stars"`
	Snippet       PluginSnippet `json:"snippet"`
	CreatedAt     string        `json:"createdAt"`
}

type PluginSnippet struct {
	K8S  string `json:"k8s"`
	Yaml string `json:"yaml"`
	Toml string `json:"toml"`
}

// Scanner --------------------------------------------------------------------

func (m *HTTPMiddleware) Scan(data any) error {
	return scanJSON(data, &m)
}

func (m *HTTPMiddleware) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *TCPMiddleware) Scan(data any) error {
	return scanJSON(data, &m)
}

func (m *TCPMiddleware) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Wrappers -------------------------------------------------------------------

func (m *HTTPMiddleware) ToDynamic() *dynamic.Middleware {
	return (*dynamic.Middleware)(m)
}

func (m *TCPMiddleware) ToDynamic() *dynamic.TCPMiddleware {
	return (*dynamic.TCPMiddleware)(m)
}

func WrapMiddleware(m *dynamic.Middleware) *HTTPMiddleware {
	return (*HTTPMiddleware)(m)
}

func WrapTCPMiddleware(m *dynamic.TCPMiddleware) *TCPMiddleware {
	return (*TCPMiddleware)(m)
}
