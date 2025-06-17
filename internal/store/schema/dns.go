package schema

type DNSProviderConfig struct {
	APIKey     string `json:"apiKey"`
	APIUrl     string `json:"apiUrl"`
	TraefikIP  string `json:"traefikIp"`
	Proxied    bool   `json:"proxied"`
	AutoUpdate bool   `json:"autoUpdate"`
	ZoneType   string `json:"zoneType"`
}
