package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

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

type TraefikConfiguration struct {
	*dynamic.Configuration
}

// Handles the JSON marshalling and unmarshalling of the TraefikEntryPoints type
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

// Handles the JSON marshalling and unmarshalling of the TraefikOverview type
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

// Handles the JSON marshalling and unmarshalling of the ConfigurationWrapper type
func (c *TraefikConfiguration) Scan(value interface{}) error {
	if value == nil {
		c.Configuration = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Configuration: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, &c.Configuration)
}

// Value implements driver.Valuer
func (c TraefikConfiguration) Value() (driver.Value, error) {
	if c.Configuration == nil {
		return nil, nil
	}
	return json.Marshal(c.Configuration)
}
