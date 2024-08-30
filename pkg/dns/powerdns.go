package dns

import (
	"context"
	"fmt"
	"net/url"

	"github.com/joeig/go-powerdns/v3"
)

type PowerDNSProvider struct {
	Client     *powerdns.Client
	ManagedTXT string
}

func NewPowerDNSProvider(url, key string) *PowerDNSProvider {
	client := powerdns.NewClient(url, "localhost", map[string]string{"X-API-Key": key}, nil)

	return &PowerDNSProvider{
		Client:     client,
		ManagedTXT: "managed-by=mantrae",
	}
}

func (p *PowerDNSProvider) CreateRecord(subdomain, ip string) error {
	recordType, err := RecordType(ip)
	if err != nil {
		return err
	}
	pdnsRecordType := powerdns.RRTypeA
	if recordType == "AAAA" {
		pdnsRecordType = powerdns.RRTypeAAAA
	}

	// Check if the record is managed by us
	managed, err := p.CheckRecord(subdomain)
	if err != nil {
		return err
	}

	// Fetch existing records
	existingRecords, err := p.ListRecords(subdomain, ip)
	if err != nil {
		return err
	}

	// If not managed by us and records exist, return an error
	if len(existingRecords) > 0 && !managed {
		return fmt.Errorf("record not managed by Mantrae")
	}

	u, err := url.Parse(subdomain)
	if err != nil {
		return err
	}

	err = p.Client.Records.Add(
		context.Background(),
		u.Host,
		subdomain,
		pdnsRecordType,
		60,
		[]string{ip},
	)
	if err != nil {
		return err
	}

	err = p.Client.Records.Add(
		context.Background(),
		u.Host,
		"_mantrae-"+subdomain,
		powerdns.RRTypeTXT,
		60,
		[]string{p.ManagedTXT},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PowerDNSProvider) DeleteRecord(subdomain, ip string) error {
	managed, err := p.CheckRecord(subdomain)
	if err != nil {
		return err
	}

	if !managed {
		return fmt.Errorf("record not managed by Mantrae")
	}

	records, err := p.ListRecords(subdomain, ip)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found for %s", subdomain)
	}

	u, err := url.Parse(subdomain)
	if err != nil {
		return err
	}

	for _, rrset := range records {
		err := p.Client.Records.
			Delete(context.Background(), u.Host, *rrset.Name, *rrset.Type)
		if err != nil {
			return fmt.Errorf("failed to delete record %s: %w", *rrset.Name, err)
		}
	}

	return nil
}

func (p *PowerDNSProvider) ListRecords(subdomain, ip string) ([]powerdns.RRset, error) {
	u, err := url.Parse(subdomain)
	if err != nil {
		return nil, err
	}

	recordType, err := RecordType(ip)
	if err != nil {
		return nil, err
	}

	pdnsRecordType := powerdns.RRTypeA
	if recordType == "AAAA" {
		pdnsRecordType = powerdns.RRTypeAAAA
	}

	records, err := p.Client.Records.Get(
		context.Background(),
		u.Host,
		subdomain,
		powerdns.RRTypePtr(pdnsRecordType),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve records for %s: %w", subdomain, err)
	}

	return records, nil
}

func (p *PowerDNSProvider) CheckRecord(subdomain string) (bool, error) {
	u, err := url.Parse(subdomain)
	if err != nil {
		return false, err
	}
	records, err := p.Client.Records.Get(
		context.Background(),
		u.Host,
		"_mantrae-"+subdomain,
		powerdns.RRTypePtr(powerdns.RRTypeTXT),
	)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve records for %s: %w", subdomain, err)
	}

	name := "_mantrae-" + subdomain
	for _, rrset := range records {
		if rrset.Name == &name && rrset.Type == powerdns.RRTypePtr(powerdns.RRTypeTXT) {
			for _, record := range rrset.Records {
				if record.Content == &p.ManagedTXT {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
