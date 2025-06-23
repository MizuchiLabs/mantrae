package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
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

type DNSRouterInfo struct {
	RouterName  string
	ProfileName string
	Provider    DNSProvider
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

	provider, err := q.GetDnsProvider(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if provider.Config.APIKey == "" {
		return nil, fmt.Errorf("invalid provider config")
	}
	if provider.Config.AutoUpdate {
		machineIPs, err := util.GetPublicIPsCached()
		if err != nil {
			return nil, err
		}
		if machineIPs.IPv4 != "" {
			provider.Config.IP = machineIPs.IPv4
		} else if machineIPs.IPv6 != "" {
			provider.Config.IP = machineIPs.IPv6
		}
	}

	var dnsProvider DNSProvider
	switch provider.Type {
	case "cloudflare":
		dnsProvider = NewCloudflareProvider(provider.Config)
	case "powerdns":
		dnsProvider = NewPowerDNSProvider(provider.Config)
	case "technitium":
		dnsProvider = NewTechnitiumProvider(provider.Config)
	default:
		return nil, fmt.Errorf("invalid provider type")
	}

	if dnsProvider == nil {
		return nil, fmt.Errorf("failed to initialize %s provider", provider.Type)
	}

	return dnsProvider, nil
}

// UpdateDNS updates the DNS records for all locally managed domains
func UpdateDNS(ctx context.Context, q *db.Queries) (err error) {
	domainMap, err := getDomainConfig(ctx, q)
	if err != nil {
		return err
	}

	for domain, entries := range domainMap {
		for _, entry := range entries {
			err := entry.Provider.UpsertRecord(domain)
			if err != nil {
				slog.Error("Failed to upsert DNS record",
					"domain", domain,
					"router", entry.RouterName,
					"profile", entry.ProfileName,
					"error", err,
				)
			}
		}
	}

	return nil
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func DeleteDNS(ctx context.Context, q *db.Queries, proto string, routerID int64) error {
	switch proto {
	case "http":
		router, err := q.GetHttpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get traefik config: %w", err)
		}
		domains, err := util.ExtractDomainFromRule(router.Config.Rule)
		if err != nil {
			return fmt.Errorf("failed to extract domains: %w", err)
		}

		providers, err := q.GetDnsProvidersByHttpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get DNS provider: %w", err)
		}

		for _, p := range providers {
			provider, err := getProvider(p.ID, q)
			if err != nil {
				slog.Error("Failed to get provider", "error", err, "provider", p)
				continue
			}
			for _, domain := range domains {
				if err := provider.DeleteRecord(domain); err != nil {
					slog.Error("Failed to delete DNS record",
						"domain", domain,
						"error", err,
					)
				}
			}
		}
	case "tcp":
		router, err := q.GetTcpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get traefik config: %w", err)
		}
		domains, err := util.ExtractDomainFromRule(router.Config.Rule)
		if err != nil {
			return fmt.Errorf("failed to extract domains: %w", err)
		}

		providers, err := q.GetDnsProvidersByTcpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get DNS provider: %w", err)
		}

		for _, p := range providers {
			provider, err := getProvider(p.ID, q)
			if err != nil {
				slog.Error("Failed to get provider", "error", err, "provider", p)
				continue
			}
			for _, domain := range domains {
				if err := provider.DeleteRecord(domain); err != nil {
					slog.Error("Failed to delete DNS record",
						"domain", domain,
						"error", err,
					)
				}
			}
		}

	default:
		return fmt.Errorf("invalid protocol: %s", proto)
	}

	return nil
}

// Result: map from domain → slice of provider info
func getDomainConfig(ctx context.Context, q *db.Queries) (map[string][]DNSRouterInfo, error) {
	domainMap := make(map[string][]DNSRouterInfo)

	process := func(
		routerName string,
		profileName string,
		rule string,
		providerID int64,
		q *db.Queries,
	) error {
		provider, err := getProvider(providerID, q)
		if err != nil {
			slog.Warn("Unable to load provider", "id", providerID, "err", err)
			return nil // soft fail
		}

		domains, err := util.ExtractDomainFromRule(rule)
		if err != nil {
			return fmt.Errorf("failed to extract domain from rule '%s': %w", rule, err)
		}

		for _, domain := range domains {
			domainMap[domain] = append(domainMap[domain], DNSRouterInfo{
				RouterName:  routerName,
				ProfileName: profileName,
				Provider:    provider,
			})

			slog.Debug("Mapped domain to provider",
				"domain", domain,
				"router", routerName,
				"profile", profileName,
			)
		}
		return nil
	}

	// HTTP
	if httpRouters, err := q.GetHttpRouterDomains(ctx); err != nil {
		return nil, err
	} else {
		for _, r := range httpRouters {
			if r.DnsProviderID == nil {
				continue
			}
			if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Rule, *r.DnsProviderID, q); err != nil {
				return nil, err
			}
		}
	}

	// TCP
	if tcpRouters, err := q.GetTcpRouterDomains(ctx); err != nil {
		return nil, err
	} else {
		for _, r := range tcpRouters {
			if r.DnsProviderID == nil {
				continue
			}
			if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Rule, *r.DnsProviderID, q); err != nil {
				return nil, err
			}
		}
	}

	return domainMap, nil
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
