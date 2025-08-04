// Package dns provides functionality for managing DNS records.
package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
)

type DNSManager struct {
	app *config.App
}

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

var managedTXT = "\"managed-by=mantrae\""

func NewManager(app *config.App) *DNSManager {
	return &DNSManager{app: app}
}

// UpdateDNS updates the DNS records for all locally managed domains
func (d *DNSManager) UpdateDNS(ctx context.Context) (err error) {
	domainMap, err := d.getDomainConfig(ctx)
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
func (d *DNSManager) DeleteDNS(ctx context.Context, proto, routerID string) error {
	switch proto {
	case "http":
		router, err := d.app.Conn.GetQuery().GetHttpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get traefik config: %w", err)
		}
		domains, err := util.ExtractDomainFromRule(router.Config.Rule)
		if err != nil {
			return fmt.Errorf("failed to extract domains: %w", err)
		}

		providers, err := d.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get DNS provider: %w", err)
		}

		for _, p := range providers {
			provider, err := d.getProvider(p.ID)
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
		router, err := d.app.Conn.GetQuery().GetTcpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get traefik config: %w", err)
		}
		domains, err := util.ExtractDomainFromRule(router.Config.Rule)
		if err != nil {
			return fmt.Errorf("failed to extract domains: %w", err)
		}

		providers, err := d.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, routerID)
		if err != nil {
			return fmt.Errorf("failed to get DNS provider: %w", err)
		}

		for _, p := range providers {
			provider, err := d.getProvider(p.ID)
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

func (d *DNSManager) getProvider(id string) (DNSProvider, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid provider id")
	}

	provider, err := d.app.Conn.GetQuery().GetDnsProvider(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if provider.Config.APIKey == "" {
		return nil, fmt.Errorf("invalid provider config")
	}
	decryptedAPIKey, err := util.DecryptSecret(provider.Config.APIKey, d.app.Secret)
	if err != nil {
		return nil, err
	}
	provider.Config.APIKey = decryptedAPIKey

	if provider.Config.AutoUpdate {
		machineIPs, err := util.GetPublicIPs()
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
	case int64(mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_CLOUDFLARE):
		dnsProvider = NewCloudflareProvider(provider.Config)
	case int64(mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_POWERDNS):
		dnsProvider = NewPowerDNSProvider(provider.Config)
	case int64(mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_TECHNITIUM):
		dnsProvider = NewTechnitiumProvider(provider.Config)
	default:
		return nil, fmt.Errorf("invalid provider type")
	}

	if dnsProvider == nil {
		return nil, fmt.Errorf("failed to initialize provider")
	}

	return dnsProvider, nil
}

// Result: map from domain → slice of provider info
func (d *DNSManager) getDomainConfig(ctx context.Context) (map[string][]DNSRouterInfo, error) {
	domainMap := make(map[string][]DNSRouterInfo)

	process := func(
		routerName string,
		profileName string,
		rule string,
		providerID string,
	) error {
		provider, err := d.getProvider(providerID)
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
	if httpRouters, err := d.app.Conn.GetQuery().GetHttpRouterDomains(ctx); err != nil {
		return nil, err
	} else {
		for _, r := range httpRouters {
			if r.DnsProviderID == nil {
				continue
			}
			if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Rule, *r.DnsProviderID); err != nil {
				return nil, err
			}
		}
	}

	// TCP
	if tcpRouters, err := d.app.Conn.GetQuery().GetTcpRouterDomains(ctx); err != nil {
		return nil, err
	} else {
		for _, r := range tcpRouters {
			if r.DnsProviderID == nil {
				continue
			}
			if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Rule, *r.DnsProviderID); err != nil {
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
