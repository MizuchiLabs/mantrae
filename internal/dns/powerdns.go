package dns

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/joeig/go-powerdns/v3"
)

type PowerDNSProvider struct {
	Client     *powerdns.Client
	ExternalIP string
}

func NewPowerDNSProvider(d *db.DNSProviderConfig) *PowerDNSProvider {
	client := powerdns.New(d.APIUrl, "", powerdns.WithAPIKey(d.APIKey))

	return &PowerDNSProvider{
		Client:     client,
		ExternalIP: d.TraefikIP,
	}
}

func (p *PowerDNSProvider) UpsertRecord(subdomain string) error {
	var recordType powerdns.RRType
	if util.IsValidIPv4(p.ExternalIP) {
		recordType = powerdns.RRTypeA
	} else if util.IsValidIPv6(p.ExternalIP) {
		recordType = powerdns.RRTypeAAAA
	} else {
		return fmt.Errorf("invalid IP address: %s", p.ExternalIP)
	}

	// Check if the record is managed by us
	if err := p.checkRecord(subdomain); err != nil {
		return err
	}

	// Fetch existing records
	records, err := p.ListRecords(subdomain)
	if err != nil {
		return err
	}

	shouldUpdate := verifyRecords(records, subdomain, p.ExternalIP)
	if len(records) <= 1 {
		if err := p.createRecord(subdomain, recordType); err != nil {
			return err
		}
		slog.Info("Created record", "name", subdomain, "type", recordType, "content", p.ExternalIP)
	} else if shouldUpdate {
		for _, record := range records {
			if record.Type != "TXT" {
				if err := p.updateRecord(record.ID, recordType, subdomain); err != nil {
					return err
				}
				slog.Info("Updated record", "name", record.Name, "type", record.Type, "content", record.Content)
			}
		}
	}

	return nil
}

func (p *PowerDNSProvider) createRecord(subdomain string, recordType powerdns.RRType) error {
	domain, err := getBaseDomain(subdomain)
	if err != nil {
		return err
	}

	// Create the A/AAAA record
	err = p.Client.Records.Add(
		context.Background(),
		domain,
		subdomain,
		recordType,
		60,
		[]string{p.ExternalIP},
	)
	if err != nil {
		return err
	}

	// Create the TXT record
	err = p.Client.Records.Add(
		context.Background(),
		domain,
		"_mantrae-"+subdomain,
		powerdns.RRTypeTXT,
		60,
		[]string{managedTXT},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PowerDNSProvider) updateRecord(
	recordID string,
	recordType powerdns.RRType,
	subdomain string,
) error {
	domain, err := getBaseDomain(subdomain)
	if err != nil {
		return err
	}

	err = p.Client.Records.Change(
		context.Background(),
		domain,
		subdomain,
		recordType,
		60,
		[]string{p.ExternalIP},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PowerDNSProvider) DeleteRecord(subdomain string) error {
	domain, err := getBaseDomain(subdomain)
	if err != nil {
		return err
	}

	if err = p.checkRecord(subdomain); err != nil {
		return err
	}

	records, err := p.ListRecords(subdomain)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("no records found for %s", subdomain)
	}

	for _, record := range records {
		err := p.Client.Records.Delete(
			context.Background(),
			domain,
			record.Name,
			powerdns.RRType(record.Type),
		)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.Name, err)
		}

		err = p.Client.Records.Delete(
			context.Background(),
			domain,
			"_mantrae-"+subdomain,
			powerdns.RRTypeTXT,
		)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", "_mantrae-"+subdomain, err)
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

func (p *PowerDNSProvider) ListRecords(subdomain string) ([]DNSRecord, error) {
	domain, err := getBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	records, err := p.Client.Records.Get(
		context.Background(),
		domain,
		subdomain,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve records for %s: %w", subdomain, err)
	}

	var dnsRecords []DNSRecord
	for _, record := range records {
		dnsRecords = append(dnsRecords, DNSRecord{
			Name:    *record.Name,
			Type:    string(*record.Type),
			Content: *record.Records[0].Content,
		})
	}

	return dnsRecords, nil
}

func (p *PowerDNSProvider) checkRecord(subdomain string) error {
	domain, err := getBaseDomain(subdomain)
	if err != nil {
		return err
	}

	records, err := p.Client.Records.Get(
		context.Background(),
		domain,
		"_mantrae-"+subdomain,
		powerdns.RRTypePtr(powerdns.RRTypeTXT),
	)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	name := "_mantrae-" + subdomain + "."
	for _, rrset := range records {
		if *rrset.Name == name {
			for _, record := range rrset.Records {
				if *record.Content == managedTXT {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("record not managed by Mantrae")
}
