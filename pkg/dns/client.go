package dns

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"regexp"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"golang.org/x/net/publicsuffix"
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

type DomainProvider struct {
	Domain   string
	Provider DNSProvider
}

func getProvider(id *int64) DNSProvider {
	if id == nil {
		return nil
	}

	provider, err := db.Query.GetProviderByID(context.Background(), *id)
	if err != nil {
		slog.Error("Failed to get providers", "error", err)
		return nil
	}

	switch provider.Type {
	case "cloudflare":
		return NewCloudflareProvider(provider.ApiKey, provider.ExternalIp, provider.Proxied)
	case "powerdns":
		return NewPowerDNSProvider(*provider.ApiUrl, provider.ApiKey, provider.ExternalIp)
	default:
		slog.Error("Unknown provider type", "type", provider.Type)
	}
	return nil
}

func getDomainProviderMap() map[string]DomainProvider {
	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
		return nil
	}

	// Get all local
	domainProviderMap := make(map[string]DomainProvider)
	for _, profile := range profiles {
		config, err := db.Query.GetConfigByProfileID(context.Background(), profile.ID)
		if err != nil {
			slog.Error("Failed to get config", "error", err)
			return nil
		}

		data, err := traefik.DecodeConfig(config)
		if err != nil {
			slog.Error("Failed to decode config", "error", err)
			return nil
		}

		for _, router := range data.Routers {
			if router.DNSProvider != nil {
				provider := getProvider(router.DNSProvider)
				if provider == nil {
					continue
				}

				domain, err := extractDomainFromRule(router.Rule)
				if err != nil {
					slog.Error("Failed to extract domain from rule", "error", err)
					continue
				}
				domainProviderMap[domain] = DomainProvider{
					Domain:   domain,
					Provider: provider,
				}
			}
		}
	}

	return domainProviderMap
}

// UpdateDNS updates the DNS records for all locally managed domains
func UpdateDNS() {
	dpMap := getDomainProviderMap()
	if dpMap == nil {
		return
	}

	// Update all DNS records
	for _, dp := range dpMap {
		if err := dp.Provider.UpsertRecord(dp.Domain); err != nil {
			slog.Error("Failed to upsert record", "error", err)
		}
	}
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func DeleteDNS(router traefik.Router) {
	dnsProvider := getProvider(router.DNSProvider)
	if dnsProvider == nil {
		return
	}

	subdomain, err := extractDomainFromRule(router.Rule)
	if err != nil {
		slog.Error("Failed to extract domain from rule", "error", err)
		return
	}

	if err := dnsProvider.DeleteRecord(subdomain); err != nil {
		slog.Error("Failed to delete record", "error", err)
		return
	}
}

// Sync periodically syncs the DNS records
func Sync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 300)
	defer ticker.Stop()

	UpdateDNS()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			UpdateDNS()
		}
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

func verifyRecords(records []DNSRecord, subdomain string, content string) bool {
	update := false

L:
	for _, record := range records {
		switch record.Type {
		case "A":
			if record.Content != content {
				update = true
				break L
			}
		case "AAAA":
			if record.Content != content {
				update = true
				break L
			}
		case "TXT":
			if record.Name != "_mantrae-"+subdomain {
				update = true
				break L
			}
		default:
			return false
		}
	}
	return update
}
