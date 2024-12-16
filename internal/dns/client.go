package dns

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
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

func getProvider(id *int64) (DNSProvider, error) {
	if id == nil || *id == 0 {
		return nil, fmt.Errorf("invalid provider id")
	}

	provider, err := db.Query.GetProviderByID(context.Background(), *id)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(DNSProviders, provider.Type) {
		return nil, fmt.Errorf("invalid provider type")
	}

	switch provider.Type {
	case "cloudflare":
		return NewCloudflareProvider(provider.ApiKey, provider.ExternalIp, provider.Proxied), nil
	case "powerdns":
		return NewPowerDNSProvider(*provider.ApiUrl, provider.ApiKey, provider.ExternalIp), nil
	case "technitium":
		return NewTechnitiumProvider(
			*provider.ApiUrl,
			provider.ApiKey,
			provider.ExternalIp,
			*provider.ZoneType,
		), nil
	default:
		return nil, fmt.Errorf("invalid provider type")
	}
}

// UpdateDNS updates the DNS records for all locally managed domains
func UpdateDNS() {
	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		slog.Error("Failed to get profiles", "error", err)
	}

	// Get all local
	for _, profile := range profiles {
		routers, err := db.Query.ListRoutersByProfileID(context.Background(), profile.ID)
		if err != nil {
			slog.Error("Failed to get routers", "error", err)
			continue
		}

		for i, router := range routers {
			if router.DnsProvider != nil {
				provider, err := getProvider(router.DnsProvider)
				if err != nil {
					slog.Error("Failed to get provider", "error", err)

					// Delete provider from router
					router.DnsProvider = nil
					routers[i] = router
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
						router.UpdateError("dns", err.Error())
					} else {
						router.UpdateError("dns", "")
					}
				}

			}
			// Update routers
			if _, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams{
				ID:          router.ID,
				ProfileID:   router.ProfileID,
				Name:        router.Name,
				DnsProvider: router.DnsProvider,
			}); err != nil {
				slog.Error("Failed to update routers", "error", err)
			}
		}
	}
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func DeleteDNS(router db.Router) {
	provider, err := getProvider(router.DnsProvider)
	if err != nil {
		slog.Error("Failed to get provider", "error", err)
		return
	}

	domains, err := util.ExtractDomainFromRule(router.Rule)
	if err != nil {
		slog.Error("Failed to extract domain from rule", "error", err)
		return
	}

	for _, domain := range domains {
		if err := provider.DeleteRecord(domain); err != nil {
			slog.Error("Failed to delete record", "error", err)
			return
		}
	}
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
