package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func DNSProviderToProto(p *db.DnsProvider) *mantraev1.DnsProvider {
	var dnsType mantraev1.DnsProviderType
	switch p.Type {
	case "cloudflare":
		dnsType = mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_CLOUDFLARE
	case "powerdns":
		dnsType = mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_POWERDNS
	case "technitium":
		dnsType = mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_TECHNITIUM
	default:
		return nil
	}

	return &mantraev1.DnsProvider{
		Id:   p.ID,
		Name: p.Name,
		Type: dnsType,
		Config: &mantraev1.DnsProviderConfig{
			ApiKey:     p.Config.APIKey,
			ApiUrl:     p.Config.APIUrl,
			Ip:         p.Config.IP,
			Proxied:    p.Config.Proxied,
			AutoUpdate: p.Config.AutoUpdate,
			ZoneType:   p.Config.ZoneType,
		},
		IsDefault: p.IsDefault,
		CreatedAt: SafeTimestamp(p.CreatedAt),
		UpdatedAt: SafeTimestamp(p.UpdatedAt),
	}
}

func DNSProvidersToProto(providers []db.DnsProvider) []*mantraev1.DnsProvider {
	var dnsProviders []*mantraev1.DnsProvider
	for _, p := range providers {
		dnsProviders = append(dnsProviders, DNSProviderToProto(&p))
	}
	return dnsProviders
}
