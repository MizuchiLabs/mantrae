package schema

import (
	"database/sql/driver"
	"encoding/json"
)

type DNSProviderConfig struct {
	APIKey     string `json:"apiKey"`
	APIUrl     string `json:"apiUrl"`
	IP         string `json:"ip"`
	Proxied    bool   `json:"proxied"`
	AutoUpdate bool   `json:"autoUpdate"`
}

func (c *DNSProviderConfig) Scan(data any) error {
	return scanJSON(data, &c)
}

func (c *DNSProviderConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}
