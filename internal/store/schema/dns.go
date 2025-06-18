package schema

type DNSProviderConfig struct {
	APIKey     string `json:"apiKey"`
	APIUrl     string `json:"apiUrl"`
	IP         string `json:"ip"`
	Proxied    bool   `json:"proxied"`
	AutoUpdate bool   `json:"autoUpdate"`
	ZoneType   string `json:"zoneType"`
}
