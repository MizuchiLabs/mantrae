package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	ServersTransport    dynamic.ServersTransport
	TCPServersTransport dynamic.TCPServersTransport
)

func (s *ServersTransport) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *ServersTransport) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *TCPServersTransport) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *TCPServersTransport) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Wrappers
func (s *ServersTransport) ToDynamic() *dynamic.ServersTransport {
	return (*dynamic.ServersTransport)(s)
}

func (s *TCPServersTransport) ToDynamic() *dynamic.TCPServersTransport {
	return (*dynamic.TCPServersTransport)(s)
}

func WrapServersTransport(s *dynamic.ServersTransport) *ServersTransport {
	return (*ServersTransport)(s)
}

func WrapTCPServersTransport(s *dynamic.TCPServersTransport) *TCPServersTransport {
	return (*TCPServersTransport)(s)
}
