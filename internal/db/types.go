package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/traefik/traefik/v3/pkg/config/runtime"
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

type TraefikVersion struct {
	Version   string `json:"version,omitempty"`
	Codename  string `json:"codename,omitempty"`
	StartDate string `json:"startDate,omitempty"`
}

type ServiceInfo struct {
	*runtime.ServiceInfo
	ServerStatus map[string]string `json:"serverStatus,omitempty"`
}

type TraefikConfiguration struct {
	Routers        map[string]*runtime.RouterInfo        `json:"routers,omitempty"`
	Middlewares    map[string]*runtime.MiddlewareInfo    `json:"middlewares,omitempty"`
	Services       map[string]*ServiceInfo               `json:"services,omitempty"`
	TCPRouters     map[string]*runtime.TCPRouterInfo     `json:"tcpRouters,omitempty"`
	TCPMiddlewares map[string]*runtime.TCPMiddlewareInfo `json:"tcpMiddlewares,omitempty"`
	TCPServices    map[string]*runtime.TCPServiceInfo    `json:"tcpServices,omitempty"`
	UDPRouters     map[string]*runtime.UDPRouterInfo     `json:"udpRouters,omitempty"`
	UDPServices    map[string]*runtime.UDPServiceInfo    `json:"udpServices,omitempty"`
}

type DNSProviderConfig struct {
	APIKey     string `json:"apiKey"`
	APIUrl     string `json:"apiUrl"`
	TraefikIP  string `json:"traefikIp"`
	Proxied    bool   `json:"proxied"`
	AutoUpdate bool   `json:"autoUpdate"`
	ZoneType   string `json:"zoneType"`
}

type AgentPrivateIPs struct {
	IPs []string `json:"privateIps,omitempty"`
}

type AgentContainer struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Image   string            `json:"image,omitempty"`
	Portmap map[int32]int32   `json:"portmap,omitempty"`
	Status  string            `json:"status,omitempty"`
	Created time.Time         `json:"created,omitempty"`
}

type AgentContainers []AgentContainer

// Handles the JSON marshalling and unmarshalling of the TraefikEntryPoints type
func (e *TraefikEntryPoints) Scan(value any) error {
	return scanJSON(value, e)
}

func (e TraefikEntryPoints) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (o *TraefikOverview) Scan(value any) error {
	return scanJSON(value, o)
}

func (o TraefikOverview) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (t *TraefikConfiguration) Scan(value any) error {
	return scanJSON(value, t)
}

func (t TraefikConfiguration) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (d *DNSProviderConfig) Scan(value any) error {
	return scanJSON(value, d)
}

func (d DNSProviderConfig) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (a *AgentPrivateIPs) Scan(value any) error {
	return scanJSON(value, a)
}

func (a AgentPrivateIPs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AgentContainers) Scan(value any) error {
	return scanJSON(value, a)
}

func (a AgentContainers) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// scanJSON is a helper function to unmarshal a JSON value into a struct
func scanJSON[T any](value any, out *T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during Scan: %v", r)
		}
	}()

	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, out)
}

// Additional conversion helpers

// Convert from runtime.ServiceInfo to your ServiceInfo
func FromRuntimeServiceInfo(ri *runtime.ServiceInfo) *ServiceInfo {
	if ri == nil {
		return nil
	}
	return &ServiceInfo{
		ServiceInfo: &runtime.ServiceInfo{
			Service: ri.Service,
			Err:     ri.Err,
			Status:  ri.Status,
			UsedBy:  ri.UsedBy,
		},
		ServerStatus: make(map[string]string),
	}
}

// Convert from your ServiceInfo to runtime.ServiceInfo
func (si *ServiceInfo) ToRuntimeServiceInfo() *runtime.ServiceInfo {
	if si == nil {
		return nil
	}
	return &runtime.ServiceInfo{
		Service: si.Service,
		Err:     si.Err,
		Status:  si.Status,
		UsedBy:  si.UsedBy,
	}
}
