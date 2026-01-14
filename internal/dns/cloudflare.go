package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type CloudflareProvider struct {
	client *cloudflare.Client
	ip     string
	proxy  bool
}

func NewCloudflareProvider(d *schema.DNSProviderConfig) *CloudflareProvider {
	if d == nil || d.APIKey == "" {
		slog.Error("Invalid Cloudflare provider config")
		return nil
	}
	return &CloudflareProvider{
		client: cloudflare.NewClient(option.WithAPIToken(d.APIKey)),
		ip:     d.IP,
		proxy:  d.Proxied,
	}
}

func (c *CloudflareProvider) UpsertRecord(ctx context.Context, subdomain string) error {
	if c.client == nil {
		return nil
	}

	rm, err := NewRecordManager(subdomain, c.ip)
	if err != nil {
		return err
	}

	records, err := c.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}

	ops := UpsertOperation{
		CreateDNSRecord: func(recordType string) error {
			return c.createRecord(ctx, subdomain, recordType)
		},
		CreateTXTMarker: func() error {
			return c.createTXTMarker(ctx, subdomain)
		},
		UpdateDNSRecord: func(recordID, recordType string) error {
			return c.updateRecord(ctx, recordID, recordType, subdomain)
		},
	}

	return rm.ExecuteUpsert(records, ops)
}

func (c *CloudflareProvider) DeleteRecord(ctx context.Context, subdomain string) error {
	if c.client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	records, err := c.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}

	rm, err := NewRecordManager(subdomain, c.ip)
	if err != nil {
		return err
	}

	if !rm.IsManagedByUs(records) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return err
	}

	for _, record := range records {
		_, err := c.client.DNS.Records.Delete(
			ctx,
			record.ID,
			dns.RecordDeleteParams{
				ZoneID: cloudflare.F(zoneID),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.ID, err)
		}
	}

	return nil
}

func (c *CloudflareProvider) createRecord(ctx context.Context, subdomain, recordType string) error {
	if c.client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return err
	}

	switch recordType {
	case "A":
		params := dns.RecordNewParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ip),
				Proxied: cloudflare.F(c.proxy),
				Type:    cloudflare.F(dns.ARecordTypeA),
			},
		}
		_, err = c.client.DNS.Records.New(ctx, params)
	case "AAAA":
		params := dns.RecordNewParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.AAAARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ip),
				Proxied: cloudflare.F(c.proxy),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			},
		}
		_, err = c.client.DNS.Records.New(ctx, params)
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}
	if err != nil {
		return fmt.Errorf("failed to create %s record: %w", recordType, err)
	}

	return nil
}

func (c *CloudflareProvider) updateRecord(
	ctx context.Context,
	recordID, recordType, subdomain string,
) error {
	if c.client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return err
	}

	switch recordType {
	case "A":
		params := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ip),
				Proxied: cloudflare.F(c.proxy),
				Type:    cloudflare.F(dns.ARecordTypeA),
			},
		}
		_, err = c.client.DNS.Records.Update(ctx, recordID, params)
	case "AAAA":
		params := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.AAAARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ip),
				Proxied: cloudflare.F(c.proxy),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			},
		}
		_, err = c.client.DNS.Records.Update(ctx, recordID, params)
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}

	if err != nil {
		return fmt.Errorf("failed to update %s record: %w", recordType, err)
	}

	return nil
}

func (c *CloudflareProvider) ListRecords(
	ctx context.Context,
	subdomain string,
) ([]DNSRecord, error) {
	if c.client == nil {
		return nil, nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("error getting zone ID for subdomain %s: %w", subdomain, err)
	}

	marker := markerName(subdomain)

	var allRecords []dns.RecordResponse

	recordsA, err := c.client.DNS.Records.List(
		ctx,
		dns.RecordListParams{
			ZoneID: cloudflare.F(zoneID),
			Type:   cloudflare.F(dns.RecordListParamsTypeA),
			Name:   cloudflare.F(dns.RecordListParamsName{Contains: cloudflare.F(subdomain)}),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing A records for subdomain %s: %w", subdomain, err)
	}
	allRecords = append(allRecords, recordsA.Result...)

	recordsAAAA, err := c.client.DNS.Records.List(
		ctx,
		dns.RecordListParams{
			ZoneID: cloudflare.F(zoneID),
			Type:   cloudflare.F(dns.RecordListParamsTypeAAAA),
			Name:   cloudflare.F(dns.RecordListParamsName{Contains: cloudflare.F(subdomain)}),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing AAAA records for subdomain %s: %w", subdomain, err)
	}
	allRecords = append(allRecords, recordsAAAA.Result...)

	recordsTXT, err := c.client.DNS.Records.List(
		ctx,
		dns.RecordListParams{
			ZoneID: cloudflare.F(zoneID),
			Type:   cloudflare.F(dns.RecordListParamsTypeTXT),
			Name:   cloudflare.F(dns.RecordListParamsName{Contains: cloudflare.F(marker)}),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing TXT records for subdomain %s: %w", subdomain, err)
	}
	allRecords = append(allRecords, recordsTXT.Result...)

	// Filter exact names to avoid "Contains" false positives.
	var out []DNSRecord
	for _, record := range allRecords {
		if record.Name != subdomain && record.Name != marker {
			continue
		}
		out = append(out, DNSRecord{
			ID:      record.ID,
			Name:    record.Name,
			Type:    string(record.Type),
			Content: record.Content,
		})
	}

	return out, nil
}

func (c *CloudflareProvider) createTXTMarker(ctx context.Context, subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(ctx, domain)
	if err != nil {
		return err
	}

	params := dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body: dns.TXTRecordParam{
			Name:    cloudflare.F(markerName(subdomain)),
			Content: cloudflare.F(managedTXT),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
		},
	}
	_, err = c.client.DNS.Records.New(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create TXT marker: %w", err)
	}

	return nil
}

// getZoneID retrieves the zone ID for a given domain
func (c *CloudflareProvider) getZoneID(ctx context.Context, domain string) (string, error) {
	zoneList, err := c.client.Zones.List(
		ctx,
		zones.ZoneListParams{
			Name: cloudflare.F(domain),
		},
	)
	if err != nil {
		return "", err
	}

	if len(zoneList.Result) == 0 {
		return "", fmt.Errorf("no zone found for domain: %s", domain)
	}

	return zoneList.Result[0].ID, nil
}
