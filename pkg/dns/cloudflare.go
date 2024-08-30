package dns

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareProvider struct {
	Container *cloudflare.ResourceContainer
	Client    *cloudflare.API
}

func NewCloudflareDNSProvider(key, id string) *CloudflareProvider {
	client, err := cloudflare.NewWithAPIToken(key)
	if err != nil {
		log.Fatal(err)
	}

	return &CloudflareProvider{
		Container: cloudflare.AccountIdentifier(id),
		Client:    client,
	}
}

// CreateRecord creates a new DNS record for the given subdomain and IP address
// and adds a TXT record to keep track of the subdomain
func (c *CloudflareProvider) CreateRecord(subdomain, ip string) error {
	recordType, err := RecordType(ip)
	if err != nil {
		return err
	}

	recordA := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    subdomain,
		Content: ip,
		Proxied: BoolPointer(false),
	}

	recordTXT := cloudflare.CreateDNSRecordParams{
		Type:    "TXT",
		Name:    "_mantrae-" + subdomain,
		Content: "managed-by=mantrae",
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

	// If the record doesn't exist, create it, or if it's managed by us, update it
	if len(records) == 0 || managed {
		_, err := c.Client.CreateDNSRecord(context.Background(), c.Container, recordA)
		if err != nil {
			return err
		}

		_, err = c.Client.CreateDNSRecord(context.Background(), c.Container, recordTXT)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("record not managed by Mantrae")
	}

	return nil
}

func (c *CloudflareProvider) UpdateRecord(subdomain, ip string) error {
	recordType, err := RecordType(ip)
	if err != nil {
		return err
	}

	paramsA := cloudflare.UpdateDNSRecordParams{
		Type:    recordType,
		Name:    subdomain,
		Content: ip,
		Proxied: BoolPointer(false),
	}
	paramsTXT := cloudflare.UpdateDNSRecordParams{
		Type:    "TXT",
		Name:    "_mantrae-" + subdomain,
		Content: "managed-by=mantrae",
	}

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

	_, err = c.Client.UpdateDNSRecord(context.Background(), c.Container, paramsA)
	if err != nil {
		return err
	}
	_, err = c.Client.UpdateDNSRecord(context.Background(), c.Container, paramsTXT)
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudflareProvider) DeleteRecord(subdomain, ip string) error {
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

	for _, record := range records {
		if err := c.Client.DeleteDNSRecord(context.Background(), c.Container, record.ID); err != nil {
			return err
		}
	}
	return nil
}

func (c *CloudflareProvider) ListRecords(subdomain string) ([]cloudflare.DNSRecord, error) {
	params := cloudflare.ListDNSRecordsParams{Name: subdomain}

	records, _, err := c.Client.ListDNSRecords(context.Background(), c.Container, params)
	if err != nil {
		return nil, fmt.Errorf("error listing A records for subdomain %s: %w", subdomain, err)
	}

	return records, nil
}

func (c *CloudflareProvider) CheckRecord(subdomain string) (bool, error) {
	params := cloudflare.ListDNSRecordsParams{
		Type: "TXT",
		Name: "_mantrae-" + subdomain,
	}

	records, _, err := c.Client.ListDNSRecords(context.Background(), c.Container, params)
	if err != nil {
		return false, fmt.Errorf("error checking TXT record for subdomain %s: %w", subdomain, err)
	}

	for _, record := range records {
		if record.Content == "managed-by=mantrae" {
			return true, nil
		}
	}

	return false, nil
}
