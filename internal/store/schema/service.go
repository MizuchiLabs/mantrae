package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	HTTPService dynamic.Service
	TCPService  dynamic.TCPService
	UDPService  dynamic.UDPService
)

// Scanner --------------------------------------------------------------------

func (s *HTTPService) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *HTTPService) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *TCPService) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *TCPService) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *UDPService) Scan(data any) error {
	return scanJSON(data, &s)
}

func (s *UDPService) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Wrappers -------------------------------------------------------------------

func (s *HTTPService) ToDynamic() *dynamic.Service {
	return (*dynamic.Service)(s)
}

func (s *TCPService) ToDynamic() *dynamic.TCPService {
	return (*dynamic.TCPService)(s)
}

func (s *UDPService) ToDynamic() *dynamic.UDPService {
	return (*dynamic.UDPService)(s)
}

func WrapService(s *dynamic.Service) *HTTPService {
	return (*HTTPService)(s)
}

func WrapTCPService(s *dynamic.TCPService) *TCPService {
	return (*TCPService)(s)
}

func WrapUDPService(s *dynamic.UDPService) *UDPService {
	return (*UDPService)(s)
}
