package dns

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/joeig/go-powerdns/v3"
)

type PowerDNSProvider struct {
	Client     *powerdns.Client
	ManagedTXT string
	ExternalIP string
}

func NewPowerDNSProvider(URL, key, ip string) *PowerDNSProvider {
	client := powerdns.NewClient(URL, "", map[string]string{"X-API-Key": key}, nil)

	return &PowerDNSProvider{
		Client:     client,
		ManagedTXT: "\"managed-by=mantrae\"",
		ExternalIP: ip,
	}
}

func (p *PowerDNSProvider) UpsertRecord(subdomain string) error {
	recordType, err := pdnsRecordType(p.ExternalIP)
	if err != nil {
		return err
	}

	// Check if the record is managed by us
	managed, err := p.CheckRecord(subdomain)
	if err != nil {
		return err
	}

	// Fetch existing records
	existingRecords, err := p.ListRecords(subdomain)
	if err != nil {
		return err
	}

	// If not managed by us and records exist, return an error
	if len(existingRecords) > 0 && !managed {
		return fmt.Errorf("record not managed by Mantrae")
	}

	// Create the record if it doesn't exist
	if len(existingRecords) == 0 {
		err = p.Client.Records.Add(
			context.Background(),
			getBaseDomain(subdomain),
			subdomain,
			recordType,
			60,
			[]string{p.ExternalIP},
		)
		if err != nil {
			return err
		}

		err = p.Client.Records.Add(
			context.Background(),
			getBaseDomain(subdomain),
			"_mantrae-"+subdomain,
			powerdns.RRTypeTXT,
			60,
			[]string{p.ManagedTXT},
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
			p.ExternalIP,
		)
	}

	// Update the record if it exists
	if len(existingRecords) > 0 {
		err = p.Client.Records.Change(
			context.Background(),
			getBaseDomain(subdomain),
			subdomain,
			recordType,
			60,
			[]string{p.ExternalIP},
		)
		if err != nil {
			return err
		}

		err = p.Client.Records.Change(
			context.Background(),
			getBaseDomain(subdomain),
			"_mantrae-"+subdomain,
			powerdns.RRTypeTXT,
			60,
			[]string{p.ManagedTXT},
		)
		if err != nil {
			return err
		}

		slog.Info(
			"Updated record",
			"subdomain",
			subdomain,
			"type",
			recordType,
			"content",
			p.ExternalIP,
		)
	}

	return nil
}

func (p *PowerDNSProvider) DeleteRecord(subdomain string) error {
	managed, err := p.CheckRecord(subdomain)
	if err != nil {
		return err
	}
	if !managed {
		return fmt.Errorf("record not managed by Mantrae")
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
			getBaseDomain(subdomain),
			record.Name,
			powerdns.RRType(record.Type),
		)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.Name, err)
		}

		err = p.Client.Records.Delete(
			context.Background(),
			getBaseDomain(subdomain),
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
	records, err := p.Client.Records.Get(
		context.Background(),
		getBaseDomain(subdomain),
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

func (p *PowerDNSProvider) CheckRecord(subdomain string) (bool, error) {
	records, err := p.Client.Records.Get(
		context.Background(),
		getBaseDomain(subdomain),
		"_mantrae-"+subdomain,
		powerdns.RRTypePtr(powerdns.RRTypeTXT),
	)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve records for %s: %w", subdomain, err)
	}

	name := "_mantrae-" + subdomain + "."
	for _, rrset := range records {
		if *rrset.Name == name {
			for _, record := range rrset.Records {
				if *record.Content == p.ManagedTXT {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func pdnsRecordType(ip string) (powerdns.RRType, error) {
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if net.ParseIP(ip).To4() != nil {
		return powerdns.RRTypeA, nil
	}
	return powerdns.RRTypeAAAA, nil
}
