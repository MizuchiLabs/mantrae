package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type HTTPRouters map[string]*dynamic.Router
type TCPRouters map[string]*dynamic.TCPRouter
type UDPRouters map[string]*dynamic.UDPRouter

type EntryPointAPI struct {
	Name            string       `json:"name,omitempty"`
	Address         string       `json:"address,omitempty"`
	AllowACMEByPass bool         `json:"allowACMEByPass,omitempty"`
	ReusePort       bool         `json:"reusePort,omitempty"`
	AsDefault       bool         `json:"asDefault,omitempty"`
	HTTP            HTTPConfig   `json:"http,omitempty"`
	HTTP2           *HTTP2Config `json:"http2,omitempty"`
	HTTP3           *HTTP3Config `json:"http3,omitempty"`
}
type HTTPConfig struct {
	Middlewares           []string   `json:"middlewares,omitempty"`
	TLS                   *TLSConfig `json:"tls,omitempty"`
	EncodeQuerySemicolons bool       `json:"encodeQuerySemicolons,omitempty"`
	MaxHeaderBytes        int        `json:"maxHeaderBytes,omitempty"`
}
type HTTP2Config struct {
	MaxConcurrentStreams int32 `json:"maxConcurrentStreams,omitempty"`
}
type HTTP3Config struct {
	AdvertisedPort int `json:"advertisedPort,omitempty"`
}
type TLSConfig struct {
	Options      string `json:"options,omitempty"`
	CertResolver string `json:"certResolver,omitempty"`
	// Domains      []types.Domain `json:"domains,omitempty"`
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
