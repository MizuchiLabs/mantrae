package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// JSONType is a generic wrapper for any type that needs JSON serialization
type JSONType[T any] struct {
	Data *T
}

func (j *JSONType[T]) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		if str, ok := value.(string); ok {
			bytes = []byte(str)
		} else {
			return fmt.Errorf("unsupported type: %T", value)
		}
	}
	return json.Unmarshal(bytes, &j.Data)
}

func (j JSONType[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j JSONType[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// Type aliases for cleaner usage
type (
	RouterConfig              = JSONType[dynamic.Router]
	TCPRouterConfig           = JSONType[dynamic.TCPRouter]
	UDPRouterConfig           = JSONType[dynamic.UDPRouter]
	ServiceConfig             = JSONType[dynamic.Service]
	TCPServiceConfig          = JSONType[dynamic.TCPService]
	UDPServiceConfig          = JSONType[dynamic.UDPService]
	MiddlewareConfig          = JSONType[dynamic.Middleware]
	TCPMiddlewareConfig       = JSONType[dynamic.TCPMiddleware]
	ServersTransportConfig    = JSONType[dynamic.ServersTransport]
	TCPServersTransportConfig = JSONType[dynamic.TCPServersTransport]
	DNSProviderConfig         = JSONType[mantraev1.DNSProviderConfig]
)
