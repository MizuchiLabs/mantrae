// Package dns provides functionality for managing DNS records.
package dns

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type DNSManager struct {
	conn   *store.Connection
	secret string
}

type DNSProvider interface {
	UpsertRecord(ctx context.Context, subdomain string) error
	DeleteRecord(ctx context.Context, subdomain string) error
	ListRecords(ctx context.Context, subdomain string) ([]DNSRecord, error)
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

func NewManager(conn *store.Connection, secret string) *DNSManager {
	return &DNSManager{conn: conn, secret: secret}
}

// UpdateDNS updates the DNS records for all locally managed domains
func (d *DNSManager) UpdateDNS() {
	subdomains := d.getSubdomains()

	for sub, entries := range subdomains {
		for _, entry := range entries {
			if entry.Provider == nil {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			slog.Info("Updating DNS record", "domain", sub)
			if err := entry.Provider.UpsertRecord(ctx, sub); err != nil {
				slog.Error("Failed to update DNS record", "domain", sub, "error", err)
			}
			cancel()
		}
	}
}

// DeleteDNS deletes the DNS record for a router if it's managed by us
func (d *DNSManager) DeleteDNS(providerID, rule string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	domains, err := util.ExtractDomainFromRule(rule)
	if err != nil {
		slog.Error("Failed to extract domains", "error", err)
		return
	}

	provider, err := d.getProvider(providerID)
	if err != nil {
		slog.Error("Failed to get provider", "error", err)
		return
	}
	for _, domain := range domains {
		if err := provider.DeleteRecord(ctx, domain); err != nil {
			slog.Error("Failed to delete DNS record", "domain", domain, "error", err)
		}
	}
}

func (d *DNSManager) getProvider(id string) (DNSProvider, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid provider id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	provider, err := d.conn.Q.GetDnsProvider(ctx, id)
	if err != nil {
		return nil, err
	}
	if provider.Config.Data.ApiKey == "" {
		return nil, fmt.Errorf("invalid provider config")
	}
	decryptedAPIKey, err := util.DecryptSecret(provider.Config.Data.ApiKey, d.secret)
	if err != nil {
		return nil, err
	}
	provider.Config.Data.ApiKey = decryptedAPIKey

	if provider.Config.Data.AutoUpdate {
		machineIPs, err := util.GetPublicIPs()
		if err != nil {
			return nil, err
		}
		if machineIPs.IPv4 != "" {
			provider.Config.Data.Ip = machineIPs.IPv4
		} else if machineIPs.IPv6 != "" {
			provider.Config.Data.Ip = machineIPs.IPv6
		}
	}

	var dnsProvider DNSProvider
	switch provider.Type {
	case int64(mantraev1.DNSProviderType_DNS_PROVIDER_TYPE_CLOUDFLARE):
		dnsProvider = NewCloudflareProvider(provider.Config.Data)
	case int64(mantraev1.DNSProviderType_DNS_PROVIDER_TYPE_POWERDNS):
		dnsProvider = NewPowerDNSProvider(provider.Config.Data)
	case int64(mantraev1.DNSProviderType_DNS_PROVIDER_TYPE_TECHNITIUM):
		dnsProvider = NewTechnitiumProvider(provider.Config.Data)
	case int64(mantraev1.DNSProviderType_DNS_PROVIDER_TYPE_PIHOLE):
		dnsProvider = NewPiholeProvider(provider.Config.Data)
	default:
		return nil, fmt.Errorf("invalid provider type")
	}

	if dnsProvider == nil {
		return nil, fmt.Errorf("failed to initialize provider")
	}

	return dnsProvider, nil
}

// Result: map from subdomain → slice of provider info
func (d *DNSManager) getSubdomains() map[string][]DNSRouterInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	domainMap := make(map[string][]DNSRouterInfo)
	process := func(routerName, profileName, rule, providerID string) error {
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
		}
		return nil
	}

	// HTTP
	httpRouters, err := d.conn.Q.GetHttpRouterDomains(ctx)
	if err != nil {
		slog.Error("Failed to get HTTP routers", "error", err)
		return nil
	}
	for _, r := range httpRouters {
		if r.DnsProviderID == nil {
			continue
		}
		if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Data.Rule, *r.DnsProviderID); err != nil {
			slog.Error("Failed to process HTTP router", "router", r.RouterName, "error", err)
			return nil
		}
	}

	// TCP
	tcpRouters, err := d.conn.Q.GetTcpRouterDomains(ctx)
	if err != nil {
		slog.Error("Failed to get TCP routers", "error", err)
		return nil
	}
	for _, r := range tcpRouters {
		if r.DnsProviderID == nil {
			continue
		}
		if err := process(r.RouterName, r.ProfileName, r.ConfigJson.Data.Rule, *r.DnsProviderID); err != nil {
			slog.Error("Failed to process TCP router", "router", r.RouterName, "error", err)
			return nil
		}
	}

	return domainMap
}
