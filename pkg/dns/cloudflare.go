package dns

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareProvider struct {
	Client     *cloudflare.API
	ManagedTXT string
	ExternalIP string
	Proxied    *bool
}

func NewCloudflareProvider(key, ip string, proxied bool) *CloudflareProvider {
	client, err := cloudflare.NewWithAPIToken(key)
	if err != nil {
		log.Fatal(err)
	}

	return &CloudflareProvider{
		Client:     client,
		ManagedTXT: "managed-by=mantrae",
		ExternalIP: ip,
		Proxied:    &proxied,
	}
}

// CreateRecord creates a new DNS record for the given subdomain and IP address
// and adds a TXT record to keep track of the subdomain
func (c *CloudflareProvider) UpsertRecord(subdomain string) error {
	recordType, err := cfRecordType(c.ExternalIP)
	if err != nil {
		return err
	}

	// Check if the record already exists
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}
	// Check if the record is managed by us
	managed, err := c.CheckRecord(subdomain)
	if err != nil {
		return err
	}
	if len(records) > 0 && !managed {
		return fmt.Errorf("record not managed by Mantrae")
	}

	zoneID, err := c.Client.ZoneIDByName(getBaseDomain(subdomain))
	if err != nil {
		return err
	}

	// Create the record if it doesn't exist
	if len(records) == 0 {
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
				Content: c.ManagedTXT,
			},
		)
		if err != nil {
			return err
		}

		slog.Info(
			"Created record",
			"subdomain",
			subdomain,
			"type",
			recordType,
			"content",
			c.ExternalIP,
		)
	}

	shouldUpdate := verifyRecords(records, subdomain, c.ExternalIP)

	// If the record doesn't exist, create it, or if it's managed by us, update it
	if len(records) > 0 && shouldUpdate {
		for _, record := range records {
			if record.Type == "A" {
				_, err = c.Client.UpdateDNSRecord(
					context.Background(),
					cloudflare.ZoneIdentifier(zoneID),
					cloudflare.UpdateDNSRecordParams{
						ID:      record.ID,
						Type:    recordType,
						Name:    subdomain,
						Content: c.ExternalIP,
						Proxied: c.Proxied,
					},
				)
				if err != nil {
					return err
				}
			}

			if record.Type == "TXT" {
				_, err = c.Client.UpdateDNSRecord(
					context.Background(),
					cloudflare.ZoneIdentifier(zoneID),
					cloudflare.UpdateDNSRecordParams{
						ID:      record.ID,
						Type:    "TXT",
						Name:    "_mantrae-" + subdomain,
						Content: c.ManagedTXT,
					},
				)
				if err != nil {
					return err
				}
			}
		}

		slog.Info(
			"Updated record",
			"subdomain",
			subdomain,
			"type",
			recordType,
			"content",
			c.ExternalIP,
		)
	}

	return nil
}

func (c *CloudflareProvider) DeleteRecord(subdomain string) error {
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	managed, err := c.CheckRecord(subdomain)
	if err != nil {
		return err
	}
	if !managed {
		return fmt.Errorf("record not managed by mantrae")
	}

	zoneID, err := c.Client.ZoneIDByName(getBaseDomain(subdomain))
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

func (c *CloudflareProvider) ListRecords(subdomain string) ([]DNSRecord, error) {
	zoneID, err := c.Client.ZoneIDByName(getBaseDomain(subdomain))
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

func (c *CloudflareProvider) CheckRecord(subdomain string) (bool, error) {
	zoneID, err := c.Client.ZoneIDByName(getBaseDomain(subdomain))
	if err != nil {
		return false, fmt.Errorf("error getting zone ID for subdomain %s: %w", subdomain, err)
	}

	records, _, err := c.Client.ListDNSRecords(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{
			Type: "TXT",
			Name: "_mantrae-" + subdomain,
		},
	)
	if err != nil {
		return false, fmt.Errorf("error checking TXT record for subdomain %s: %w", subdomain, err)
	}

	for _, record := range records {
		if record.Content == c.ManagedTXT {
			return true, nil
		}
	}

	return false, nil
}

func boolPointer(b bool) *bool {
	return &b
}

func cfRecordType(ip string) (string, error) {
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if net.ParseIP(ip).To4() != nil {
		return "A", nil
	}
	return "AAAA", nil
}
