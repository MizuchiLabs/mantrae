package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	HTTPServersTransport dynamic.ServersTransport
	TCPServersTransport  dynamic.TCPServersTransport
)

// Scanner --------------------------------------------------------------------

func (s *HTTPServersTransport) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *HTTPServersTransport) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *TCPServersTransport) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *TCPServersTransport) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Wrappers -------------------------------------------------------------------

func (s *HTTPServersTransport) ToDynamic() *dynamic.ServersTransport {
	return (*dynamic.ServersTransport)(s)
}

func (s *TCPServersTransport) ToDynamic() *dynamic.TCPServersTransport {
	return (*dynamic.TCPServersTransport)(s)
}

func WrapServersTransport(s *dynamic.ServersTransport) *HTTPServersTransport {
	return (*HTTPServersTransport)(s)
}

func WrapTCPServersTransport(s *dynamic.TCPServersTransport) *TCPServersTransport {
	return (*TCPServersTransport)(s)
}
