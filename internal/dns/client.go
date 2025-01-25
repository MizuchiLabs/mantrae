package dns

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"golang.org/x/net/publicsuffix"
)

type DNSProvider interface {
	UpsertRecord(subdomain string) error
	DeleteRecord(subdomain string) error
	ListRecords(subdomain string) ([]DNSRecord, error)
}

type DNSRecord struct {
	ID      string
	Name    string
	Type    string
	Content string
}

var (
	DNSProviders = []string{"cloudflare", "powerdns", "technitium"}
	ZoneTypes    = []string{"primary", "forwarder"}
	managedTXT   = "\"managed-by=mantrae\""
)

func getProvider(id int64, q *db.Queries) (DNSProvider, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid provider id")
	}

	provider, err := q.GetDNSProvider(context.Background(), id)
	if err != nil {
		return nil, err
	}

	var d db.DNSProviderConfig
	switch provider.Type {
	case "cloudflare":
		return NewCloudflareProvider(&d), nil
	case "powerdns":
		return NewPowerDNSProvider(&d), nil
	case "technitium":
		return NewTechnitiumProvider(&d), nil
	default:
		return nil, fmt.Errorf("invalid provider type")
	}
}

// UpdateDNS updates the DNS records for all locally managed domains
func UpdateDNS(q *db.Queries) (err error) {
	traefikIDs, err := q.ListTraefikIDs(context.Background())
	if err != nil {
		return err
	}

	for _, id := range traefikIDs {
		rdps, err := q.ListRouterDNSProvidersByTraefikID(context.Background(), id)
		if err != nil {
			continue
		}

		for _, rdp := range rdps {
			provider, err := getProvider(rdp.ProviderID, q)
			if err != nil {
				slog.Error("Failed to get provider", "error", err)
				continue
			}
			config, err := q.GetTraefikConfig(context.Background(), rdp.TraefikID)
			if err != nil {
				slog.Error("Failed to get traefik config", "error", err)
				continue
			}
			router := config.Config.Routers[rdp.RouterName]
			if router == nil || router.Rule == "" {
				continue
			}
			domains, err := util.ExtractDomainFromRule(router.Rule)
			if err != nil {
				slog.Error("Failed to extract domain from rule", "error", err)
				continue
			}
			for _, domain := range domains {
				if err := provider.UpsertRecord(domain); err != nil {
					slog.Error("Failed to upsert record", "error", err)
				}
			}
		}
	}
	return nil
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func DeleteDNS(q *db.Queries, traefikID int64, routerName string) error {
	// Get DNS provider mapping before deletion
	params := db.GetRouterDNSProviderParams{
		TraefikID:  traefikID,
		RouterName: routerName,
	}

	rdp, err := q.GetRouterDNSProvider(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to get router DNS provider: %w", err)
	}

	// Get traefik config to extract domains
	config, err := q.GetTraefikConfig(context.Background(), traefikID)
	if err != nil {
		return fmt.Errorf("failed to get traefik config: %w", err)
	}

	router := config.Config.Routers[routerName]
	if router == nil {
		return fmt.Errorf("router not found: %s", routerName)
	}

	// Get domains from router rule
	domains, err := util.ExtractDomainFromRule(router.Rule)
	if err != nil {
		return fmt.Errorf("failed to extract domains: %w", err)
	}

	// Get DNS provider and delete records
	provider, err := getProvider(rdp.ProviderID, q)
	if err != nil {
		return fmt.Errorf("failed to get DNS provider: %w", err)
	}

	for _, domain := range domains {
		if err := provider.DeleteRecord(domain); err != nil {
			slog.Error("Failed to delete DNS record",
				"domain", domain,
				"error", err,
			)
		}
	}

	return nil
}

func getBaseDomain(domain string) (string, error) {
	// Ensure the domain doesn't contain a scheme
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	return publicsuffix.EffectiveTLDPlusOne(domain)
}

func verifyRecords(records []DNSRecord, subdomain string, content string) bool {
	for _, record := range records {
		switch record.Type {
		case "A":
			if record.Content != content {
				return true
			}
		case "AAAA":
			if record.Content != content {
				return true
			}
		case "TXT":
			if record.Name != "_mantrae-"+subdomain {
				return true
			}
		default:
			return false
		}
	}

	return false
}
