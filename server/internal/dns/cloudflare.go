package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
)

type CloudflareProvider struct {
	Client     *cloudflare.Client
	ExternalIP string
	Proxied    bool
}

func NewCloudflareProvider(d *schema.DNSProviderConfig) *CloudflareProvider {
	if d == nil || d.APIKey == "" {
		slog.Error("Invalid Cloudflare provider config")
		return nil
	}
	return &CloudflareProvider{
		Client:     cloudflare.NewClient(option.WithAPIToken(d.APIKey)),
		ExternalIP: d.IP,
		Proxied:    d.Proxied,
	}
}

func (c *CloudflareProvider) UpsertRecord(subdomain string) error {
	if c.Client == nil {
		return nil
	}

	recordType, err := c.getRecordType()
	if err != nil {
		return err
	}

	// List existing records
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	// Check if records are managed by us
	if !c.isManagedByUs(records, subdomain) && len(records) > 0 {
		return fmt.Errorf("record not managed by Mantrae")
	}

	// Separate DNS records from TXT marker
	dnsRecords, txtRecord := c.separateRecords(records, subdomain)

	// No records exist - create new
	if len(dnsRecords) == 0 {
		if err := c.createRecord(subdomain, recordType); err != nil {
			return err
		}
		slog.Info("Created record", "name", subdomain, "type", recordType, "content", c.ExternalIP)
		return nil
	}

	// Records exist - check if update needed
	if c.needsUpdate(dnsRecords, recordType) {
		for _, record := range dnsRecords {
			if err := c.updateRecord(record.ID, recordType, subdomain); err != nil {
				return err
			}
			slog.Info(
				"Updated record",
				"name",
				record.Name,
				"type",
				record.Type,
				"content",
				c.ExternalIP,
			)
		}
	}

	// Ensure TXT marker exists
	if txtRecord == nil {
		if err := c.createTXTMarker(subdomain); err != nil {
			return fmt.Errorf("failed to create TXT marker: %w", err)
		}
	}

	return nil
}

func (c *CloudflareProvider) DeleteRecord(subdomain string) error {
	if c.Client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	// Check if managed by us before deleting
	if !c.isManagedByUs(records, subdomain) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return err
	}

	// Delete all records (including TXT marker)
	for _, record := range records {
		_, err := c.Client.DNS.Records.Delete(
			context.Background(),
			record.ID,
			dns.RecordDeleteParams{
				ZoneID: cloudflare.F(zoneID),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.ID, err)
		}

		slog.Info("Deleted record",
			"subdomain", subdomain,
			"type", record.Type,
			"content", record.Content,
		)
	}

	return nil
}

func (c *CloudflareProvider) createRecord(subdomain, recordType string) error {
	if c.Client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return err
	}

	// Create DNS record
	switch recordType {
	case "A":
		params := dns.RecordNewParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ExternalIP),
				Proxied: cloudflare.F(c.Proxied),
				Type:    cloudflare.F(dns.ARecordTypeA),
			},
		}
		_, err = c.Client.DNS.Records.New(context.Background(), params)
	case "AAAA":
		params := dns.RecordNewParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.AAAARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ExternalIP),
				Proxied: cloudflare.F(c.Proxied),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			},
		}
		_, err = c.Client.DNS.Records.New(context.Background(), params)
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}
	if err != nil {
		return fmt.Errorf("failed to create %s record: %w", recordType, err)
	}

	// Create TXT marker
	return c.createTXTMarker(subdomain)
}

func (c *CloudflareProvider) updateRecord(recordID, recordType, subdomain string) error {
	if c.Client == nil {
		return nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return err
	}

	switch recordType {
	case "A":
		params := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ExternalIP),
				Proxied: cloudflare.F(c.Proxied),
				Type:    cloudflare.F(dns.ARecordTypeA),
			},
		}
		_, err = c.Client.DNS.Records.Update(context.Background(), recordID, params)
	case "AAAA":
		params := dns.RecordUpdateParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.AAAARecordParam{
				Name:    cloudflare.F(subdomain),
				Content: cloudflare.F(c.ExternalIP),
				Proxied: cloudflare.F(c.Proxied),
				Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			},
		}
		_, err = c.Client.DNS.Records.Update(context.Background(), recordID, params)
	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}

	if err != nil {
		return fmt.Errorf("failed to update %s record: %w", recordType, err)
	}

	return nil
}

func (c *CloudflareProvider) ListRecords(subdomain string) ([]DNSRecord, error) {
	if c.Client == nil {
		return nil, nil
	}

	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return nil, fmt.Errorf("error getting zone ID for subdomain %s: %w", subdomain, err)
	}

	var allRecords []dns.RecordResponse

	// Query A records
	recordsA, err := c.Client.DNS.Records.List(
		context.Background(),
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

	// Query AAAA records
	recordsAAAA, err := c.Client.DNS.Records.List(
		context.Background(),
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

	// Query TXT marker
	recordsTXT, err := c.Client.DNS.Records.List(
		context.Background(),
		dns.RecordListParams{
			ZoneID: cloudflare.F(zoneID),
			Type:   cloudflare.F(dns.RecordListParamsTypeTXT),
			Name: cloudflare.F(
				dns.RecordListParamsName{Contains: cloudflare.F("_mantrae-" + subdomain)},
			),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing TXT records for subdomain %s: %w", subdomain, err)
	}
	allRecords = append(allRecords, recordsTXT.Result...)

	// Convert to DNSRecord
	var dnsRecords []DNSRecord
	for _, record := range allRecords {
		dnsRecords = append(dnsRecords, DNSRecord{
			ID:      record.ID,
			Name:    record.Name,
			Type:    string(record.Type),
			Content: record.Content,
		})
	}

	return dnsRecords, nil
}

// separateRecords splits DNS records from the TXT marker
func (c *CloudflareProvider) separateRecords(
	records []DNSRecord,
	subdomain string,
) ([]DNSRecord, *DNSRecord) {
	var dnsRecords []DNSRecord
	var txtRecord *DNSRecord
	markerName := "_mantrae-" + subdomain

	for i := range records {
		if records[i].Type == "TXT" && records[i].Name == markerName {
			txtRecord = &records[i]
		} else if records[i].Type == "A" || records[i].Type == "AAAA" {
			dnsRecords = append(dnsRecords, records[i])
		}
	}

	return dnsRecords, txtRecord
}

// isManagedByUs checks if the TXT marker exists
func (c *CloudflareProvider) isManagedByUs(records []DNSRecord, subdomain string) bool {
	markerName := "_mantrae-" + subdomain
	for _, record := range records {
		if record.Name == markerName && record.Type == "TXT" && record.Content == managedTXT {
			return true
		}
	}
	return false
}

// getRecordType determines A or AAAA based on IP
func (c *CloudflareProvider) getRecordType() (string, error) {
	if util.IsValidIPv4(c.ExternalIP) {
		return "A", nil
	}
	if util.IsValidIPv6(c.ExternalIP) {
		return "AAAA", nil
	}
	return "", fmt.Errorf("invalid IP address: %s", c.ExternalIP)
}

func (c *CloudflareProvider) createTXTMarker(subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.getZoneID(domain)
	if err != nil {
		return err
	}

	params := dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body: dns.TXTRecordParam{
			Name:    cloudflare.F("_mantrae-" + subdomain),
			Content: cloudflare.F(managedTXT),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
		},
	}
	_, err = c.Client.DNS.Records.New(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create TXT marker: %w", err)
	}

	return nil
}

// needsUpdate checks if any DNS record needs updating
func (c *CloudflareProvider) needsUpdate(records []DNSRecord, expectedType string) bool {
	for _, record := range records {
		// Wrong type or wrong content
		if record.Type != expectedType || record.Content != c.ExternalIP {
			return true
		}
	}
	return false
}

// getZoneID retrieves the zone ID for a given domain
func (c *CloudflareProvider) getZoneID(domain string) (string, error) {
	zoneList, err := c.Client.Zones.List(
		context.Background(),
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
