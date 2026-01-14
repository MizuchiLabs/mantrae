package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/ryanwholey/go-pihole"
)

type PiholeProvider struct {
	client *pihole.Client
	ip     string
}

func NewPiholeProvider(d *schema.DNSProviderConfig) *PiholeProvider {
	if d == nil || d.APIKey == "" || d.APIUrl == "" {
		slog.Error("Invalid Pi-hole provider config")
		return nil
	}
	client, err := pihole.New(pihole.Config{
		BaseURL:  d.APIUrl,
		Password: d.APIKey,
	})
	if err != nil {
		slog.Error("failed to create pihole client", "error", err)
		return nil
	}
	return &PiholeProvider{
		client: client,
		ip:     d.IP,
	}
}

func (p *PiholeProvider) UpsertRecord(ctx context.Context, subdomain string) error {
	if p.client == nil {
		return nil
	}

	// Delete existing records if any
	existing, err := p.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}

	for _, record := range existing {
		if err := p.client.LocalDNS.Delete(ctx, record.Name); err != nil {
			// Log but continue - record might not exist
			slog.Warn("failed to delete existing record", "domain", record.Name, "error", err)
		}
	}

	// Create new record
	_, err = p.client.LocalDNS.Create(ctx, subdomain, p.ip)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	return nil
}

func (p *PiholeProvider) DeleteRecord(ctx context.Context, subdomain string) error {
	if p.client == nil {
		return nil
	}

	err := p.client.LocalDNS.Delete(ctx, subdomain)
	if err != nil {
		return fmt.Errorf("failed to delete record %s: %w", subdomain, err)
	}

	return nil
}

func (p *PiholeProvider) ListRecords(ctx context.Context, subdomain string) ([]DNSRecord, error) {
	if p.client == nil {
		return nil, nil
	}

	allRecords, err := p.client.LocalDNS.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}

	var records []DNSRecord

	for _, record := range allRecords {
		if record.Domain == subdomain {
			recordType := "A"
			if util.IsValidIPv6(record.IP) {
				recordType = "AAAA"
			}

			records = append(records, DNSRecord{
				ID:      record.Domain,
				Name:    record.Domain,
				Type:    recordType,
				Content: record.IP,
			})
		}
	}

	return records, nil
}
