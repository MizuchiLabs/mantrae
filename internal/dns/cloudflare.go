package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type CloudflareProvider struct {
	Client     *cloudflare.API
	ExternalIP string
	Proxied    *bool
}

func NewCloudflareProvider(d *schema.DNSProviderConfig) *CloudflareProvider {
	if d == nil || d.APIKey == "" {
		slog.Error("Invalid Cloudflare provider config")
		return nil
	}
	client, err := cloudflare.NewWithAPIToken(d.APIKey)
	if err != nil {
		slog.Error("Failed to create Cloudflare client", "error", err)
		return nil
	}

	return &CloudflareProvider{
		Client:     client,
		ExternalIP: d.TraefikIP,
		Proxied:    &d.Proxied,
	}
}

func (c *CloudflareProvider) UpsertRecord(subdomain string) error {
	if c.Client == nil {
		return nil
	}
	var recordType string
	if util.IsValidIPv4(c.ExternalIP) {
		recordType = "A"
	} else if util.IsValidIPv6(c.ExternalIP) {
		recordType = "AAAA"
	} else {
		return fmt.Errorf("invalid IP address: %s", c.ExternalIP)
	}

	// Check if the record exists
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	// Check if the record is managed by us
	if err := c.checkRecord(subdomain); err != nil {
		return err
	}

	shouldUpdate := verifyRecords(records, subdomain, c.ExternalIP)
	if len(records) <= 1 {
		if err := c.createRecord(subdomain, recordType); err != nil {
			return err
		}
		slog.Info("Created record", "name", subdomain, "type", recordType, "content", c.ExternalIP)
	} else if shouldUpdate {
		for _, record := range records {
			if record.Type != "TXT" {
				if err := c.updateRecord(record.ID, recordType, subdomain); err != nil {
					return err
				}
				slog.Info("Updated record", "name", record.Name, "type", record.Type, "content", record.Content)
			}
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

	if err = c.checkRecord(subdomain); err != nil {
		return err
	}

	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.Client.ZoneIDByName(domain)
	if err != nil {
		return err
	}

	for _, record := range records {
		if err := c.Client.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), record.ID); err != nil {
			return err
		}

		slog.Info(
			"Deleted record",
			"subdomain",
			subdomain,
			"type",
			record.Type,
			"content",
			record.Content,
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

	zoneID, err := c.Client.ZoneIDByName(domain)
	if err != nil {
		return err
	}

	// Create the A/AAAA record
	_, err = c.Client.CreateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.CreateDNSRecordParams{
			Type:    recordType,
			Name:    subdomain,
			Content: c.ExternalIP,
			Proxied: c.Proxied,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.Client.CreateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.CreateDNSRecordParams{
			Type:    "TXT",
			Name:    "_mantrae-" + subdomain,
			Content: managedTXT,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudflareProvider) updateRecord(recordID, recordType, subdomain string) error {
	if c.Client == nil {
		return nil
	}
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	zoneID, err := c.Client.ZoneIDByName(domain)
	if err != nil {
		return err
	}

	_, err = c.Client.UpdateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.UpdateDNSRecordParams{
			ID:      recordID,
			Type:    recordType,
			Name:    subdomain,
			Content: c.ExternalIP,
			Proxied: c.Proxied,
		},
	)
	if err != nil {
		return err
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

	zoneID, err := c.Client.ZoneIDByName(domain)
	if err != nil {
		return nil, fmt.Errorf("error getting zone ID for subdomain %s: %w", subdomain, err)
	}

	recordsA, _, err := c.Client.ListDNSRecords(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{Name: subdomain},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing A records for subdomain %s: %w", subdomain, err)
	}

	recordsTXT, _, err := c.Client.ListDNSRecords(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{
			Type: "TXT",
			Name: "_mantrae-" + subdomain,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing TXT records for subdomain %s: %w", subdomain, err)
	}

	var dnsRecords []DNSRecord
	for _, record := range append(recordsA, recordsTXT...) {
		dnsRecords = append(dnsRecords, DNSRecord{
			ID:      record.ID,
			Name:    record.Name,
			Type:    record.Type,
			Content: record.Content,
		})
	}

	return dnsRecords, nil
}

// checkRecord verifies if the TXT record for verification exists and is managed by us.
func (c *CloudflareProvider) checkRecord(subdomain string) error {
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		if record.Name == "_mantrae-"+subdomain && record.Type == "TXT" &&
			record.Content == managedTXT {
			return nil
		}
	}

	return fmt.Errorf("record not managed by Mantrae")
}
