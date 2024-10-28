package dns

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
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
		data, err := traefik.DecodeFromDB(profile.ID)
		if err != nil {
			slog.Error("Failed to decode config", "error", err)
		}

		for i, router := range data.Routers {
			if router.DNSProvider != nil {
				provider, err := getProvider(router.DNSProvider)
				if err != nil {
					slog.Error("Failed to get provider", "error", err)

					// Delete provider from router
					router.DNSProvider = nil
					data.Routers[i] = router
					continue
				}

				domain, err := util.ExtractDomainFromRule(router.Rule)
				if err != nil {
					slog.Error("Failed to extract domain from rule", "error", err)
					continue
				}

				if err := provider.UpsertRecord(domain); err != nil {
					slog.Error("Failed to upsert record", "error", err)
					router.ErrorState.DNS = err.Error()
					data.Routers[i] = router
				} else {
					router.ErrorState.DNS = ""
					data.Routers[i] = router
				}
			}
		}

		if _, err := traefik.EncodeToDB(data); err != nil {
			slog.Error("Failed to update config", "error", err)
		}
	}
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func DeleteDNS(router traefik.Router) {
	provider, err := getProvider(router.DNSProvider)
	if err != nil {
		slog.Error("Failed to get provider", "error", err)
		return
	}

	subdomain, err := util.ExtractDomainFromRule(router.Rule)
	if err != nil {
		slog.Error("Failed to extract domain from rule", "error", err)
		return
	}

	if err := provider.DeleteRecord(subdomain); err != nil {
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
