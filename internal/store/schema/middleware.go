package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	Middleware    dynamic.Middleware
	TCPMiddleware dynamic.TCPMiddleware
)

func (m *Middleware) Scan(data any) error {
	return scanJSON(data, &m)
}

func (m *Middleware) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *TCPMiddleware) Scan(data any) error {
	return scanJSON(data, &m)
}

func (m *TCPMiddleware) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Wrappers
func (m *Middleware) ToDynamic() *dynamic.Middleware {
	return (*dynamic.Middleware)(m)
}

func (m *TCPMiddleware) ToDynamic() *dynamic.TCPMiddleware {
	return (*dynamic.TCPMiddleware)(m)
}

func WrapMiddleware(m *dynamic.Middleware) *Middleware {
	return (*Middleware)(m)
}

func WrapTCPMiddleware(m *dynamic.TCPMiddleware) *TCPMiddleware {
	return (*TCPMiddleware)(m)
}
