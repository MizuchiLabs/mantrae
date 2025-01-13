package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
	"github.com/traefik/traefik/v3/pkg/config/static"
)

type TraefikConfig runtime.Configuration
type HTTPRouters map[string]*dynamic.Router
type TCPRouters map[string]*dynamic.TCPRouter
type UDPRouters map[string]*dynamic.UDPRouter

type EntryPointAPI struct {
	Name string `json:"name,omitempty"`
	static.EntryPoint
}
type TraefikEntryPoints []EntryPointAPI

type TraefikOverview struct {
	HTTP     SchemeOverview `json:"http,omitempty"`
	TCP      SchemeOverview `json:"tcp,omitempty"`
	UDP      SchemeOverview `json:"udp,omitempty"`
	Features struct {
		Tracing   string `json:"tracing,omitempty"`
		Metrics   string `json:"metrics,omitempty"`
		AccessLog bool   `json:"accessLog,omitempty"`
	} `json:"features,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

type Section struct {
	Total    int `json:"total,omitempty"`
	Warnings int `json:"warnings,omitempty"`
	Errors   int `json:"errors,omitempty"`
}

type SchemeOverview struct {
	Routers    Section `json:"routers,omitempty"`
	Services   Section `json:"services,omitempty"`
	Middleware Section `json:"middlewares,omitempty"`
}

func (c *TraefikConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected bytes, got %T", value)
	}
	return json.Unmarshal(bytes, (*runtime.Configuration)(c))
}

func (c TraefikConfig) Value() (driver.Value, error) {
	return json.Marshal(runtime.Configuration(c))
}

func (e *TraefikEntryPoints) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected bytes, got %T", value)
	}
	return json.Unmarshal(bytes, (*[]EntryPointAPI)(e))
}

func (e TraefikEntryPoints) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (o *TraefikOverview) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected bytes, got %T", value)
	}
	return json.Unmarshal(bytes, (*TraefikOverview)(o))
}

func (o TraefikOverview) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (r *HTTPRouters) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return fmt.Errorf("expected bytes or string, got %T", value)
		}
	}
	return json.Unmarshal(bytes, r)
}

func (r HTTPRouters) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *TCPRouters) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return fmt.Errorf("expected bytes or string, got %T", value)
		}
	}
	return json.Unmarshal(bytes, r)
}
func (r TCPRouters) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *UDPRouters) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return fmt.Errorf("expected bytes or string, got %T", value)
		}
	}
	return json.Unmarshal(bytes, r)
}
func (r UDPRouters) Value() (driver.Value, error) {
	return json.Marshal(r)
}
