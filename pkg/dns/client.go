package dns

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
)

type DNSProvider interface {
	UpsertRecord(subdomain string) error
	DeleteRecord(subdomain string) error
	ListRecords(subdomain string) ([]DNSRecord, error)
	CheckRecord(subdomain string) (bool, error)
}

type DNSRecord struct {
	ID      string
	Name    string
	Type    string
	Content string
}

func getProvider() DNSProvider {
	var p Providers
	if err := p.Load(); err != nil {
		slog.Error("Failed to load providers", "error", err)
	}

	if p.Providers == nil {
		return nil
	}

	for _, provider := range p.Providers {
		switch provider.Type {
		case "cloudflare":
			return NewCloudflareProvider(provider.APIKey, provider.ExternalIP)
		case "powerdns":
			return NewPowerDNSProvider(provider.APIURL, provider.APIKey, provider.ExternalIP)
		default:
			slog.Error("Unknown provider type", "type", provider.Type)
		}
	}
	return nil
}

// UpdateDNS updates the DNS records for all locally managed domains
func UpdateDNS() {
	dnsProvider := getProvider()
	if dnsProvider == nil {
		return
	}

	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
		return
	}

	// Get all local
	domains := make(map[string]string)
	for _, profile := range profiles {
		config, err := db.Query.GetConfigByProfileID(context.Background(), profile.ID)
		if err != nil {
			slog.Error("Failed to get config", "error", err)
			return
		}

		data, err := traefik.DecodeConfig(config)
		if err != nil {
			slog.Error("Failed to decode config", "error", err)
			return
		}

		for _, router := range data.Routers {
			if router.Provider == "http" {
				domain, err := extractDomainFromRule(router.Rule)
				if err != nil {
					slog.Error("Failed to extract domain from rule", "error", err)
					continue
				}
				domains[router.Name] = domain
			}
		}
	}
	for _, domain := range domains {
		if err := dnsProvider.UpsertRecord(domain); err != nil {
			slog.Error("Failed to upsert record", "error", err)
		}
	}
}

func DeleteDNS(router traefik.Router) {
	dnsProvider := getProvider()
	if dnsProvider == nil {
		slog.Error("No DNS provider found")
		return
	}

	subdomain, err := extractDomainFromRule(router.Rule)
	if err != nil {
		slog.Error("Failed to extract domain from rule", "error", err)
		return
	}

	if err := dnsProvider.DeleteRecord(subdomain); err != nil {
		slog.Error("Failed to delete record", "error", err)
	}
}

// Sync periodically syncs the DNS records
func Sync() {
	ticker := time.NewTicker(time.Second * 300)
	defer ticker.Stop()

	for range ticker.C {
		UpdateDNS()
	}
}

func extractDomainFromRule(rule string) (string, error) {
	// Regular expression to match the domain inside a Host(`domain.com`) rule
	re := regexp.MustCompile(`Host\(` + "`" + `([^` + "`" + `]+)` + "`" + `\)`)
	matches := re.FindStringSubmatch(rule)
	if len(matches) < 2 {
		return "", fmt.Errorf("no domain found in rule")
	}
	return matches[1], nil
}
