package schema

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type (
	HTTPRouter dynamic.Router
	TCPRouter  dynamic.TCPRouter
	UDPRouter  dynamic.UDPRouter
)

// Scanner --------------------------------------------------------------------

func (r *HTTPRouter) Scan(data any) error {
	return scanJSON(data, &r)
}

func (r *HTTPRouter) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *TCPRouter) Scan(data any) error {
	return scanJSON(data, &r)
}

func (r *TCPRouter) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *UDPRouter) Scan(data any) error {
	return scanJSON(data, &r)
}

func (r *UDPRouter) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Wrappers -------------------------------------------------------------------

func (r *HTTPRouter) ToDynamic() *dynamic.Router {
	return (*dynamic.Router)(r)
}

func (r *TCPRouter) ToDynamic() *dynamic.TCPRouter {
	return (*dynamic.TCPRouter)(r)
}

func (r *UDPRouter) ToDynamic() *dynamic.UDPRouter {
	return (*dynamic.UDPRouter)(r)
}

func WrapRouter(r *dynamic.Router) *HTTPRouter {
	return (*HTTPRouter)(r)
}

func WrapTCPRouter(r *dynamic.TCPRouter) *TCPRouter {
	return (*TCPRouter)(r)
}

func WrapUDPRouter(r *dynamic.UDPRouter) *UDPRouter {
	return (*UDPRouter)(r)
}
