package dns

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/publicsuffix"
)

type Providers struct {
	Providers map[string]Provider `json:"providers,omitempty"`
	mu        sync.RWMutex
}

type Provider struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	ExternalIP string `json:"externalIP,omitempty"`
	APIKey     string `json:"key,omitempty"`
	APIURL     string `json:"url,omitempty"`
}

func (p *Providers) Load() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(cwd, "provider.json")
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open providers file: %w", err)
	}

	p.Providers = make(map[string]Provider)
	if err = json.NewDecoder(file).Decode(&p.Providers); err != nil {
		return fmt.Errorf("failed to decode providers: %w", err)
	}

	return nil
}

func (p *Providers) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	providerPath := filepath.Join(cwd, "provider.json")

	providerBytes, err := json.Marshal(p.Providers)
	if err != nil {
		return fmt.Errorf("failed to marshal providers: %w", err)
	}

	if err := os.WriteFile(providerPath, providerBytes, 0600); err != nil {
		return fmt.Errorf("failed to write providers: %w", err)
	}

	return nil
}

func (p *Provider) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("provider type cannot be empty")
	}
	if p.ExternalIP == "" {
		return fmt.Errorf("provider external ip cannot be empty")
	}
	if p.APIKey == "" {
		return fmt.Errorf("provider api key cannot be empty")
	}

	return nil
}

func getBaseDomain(subdomain string) string {
	u, err := url.Parse(subdomain)
	if err != nil {
		log.Fatal(err)
	}
	// If the URL doesn't have a scheme, url.Parse might put the whole string in Path
	if u.Host == "" {
		u, err = url.Parse("http://" + subdomain)
		if err != nil {
			log.Fatal(err)
		}
	}

	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		log.Fatal(err)
	}

	return baseDomain
}
