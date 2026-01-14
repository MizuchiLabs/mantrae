package dns

import (
	"context"
	"fmt"

	"github.com/joeig/go-powerdns/v3"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type PowerDNSProvider struct {
	client *powerdns.Client
	ip     string
}

func NewPowerDNSProvider(d *schema.DNSProviderConfig) *PowerDNSProvider {
	client := powerdns.New(d.APIUrl, "", powerdns.WithAPIKey(d.APIKey))
	return &PowerDNSProvider{
		client: client,
		ip:     d.IP,
	}
}

func (p *PowerDNSProvider) UpsertRecord(ctx context.Context, subdomain string) error {
	rm, err := NewRecordManager(subdomain, p.ip)
	if err != nil {
		return err
	}

	records, err := p.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}

	ops := UpsertOperation{
		CreateDNSRecord: func(recordType string) error {
			return p.createRecord(ctx, subdomain, recordType)
		},
		CreateTXTMarker: func() error {
			return p.createTXTMarker(ctx, subdomain)
		},
		UpdateDNSRecord: func(_ string, recordType string) error {
			return p.updateRecord(ctx, subdomain, recordType)
		},
	}

	return rm.ExecuteUpsert(records, ops)
}

func (p *PowerDNSProvider) createRecord(
	ctx context.Context,
	subdomain string,
	recordType string,
) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	return p.client.Records.Add(
		ctx,
		domain,
		subdomain,
		powerdns.RRType(recordType),
		60,
		[]string{p.ip},
	)
}

func (p *PowerDNSProvider) updateRecord(
	ctx context.Context,
	subdomain string,
	recordType string,
) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	return p.client.Records.Change(
		ctx,
		domain,
		subdomain,
		powerdns.RRType(recordType),
		60,
		[]string{p.ip},
	)
}

func (p *PowerDNSProvider) DeleteRecord(ctx context.Context, subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	records, err := p.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("no records found for %s", subdomain)
	}

	rm, err := NewRecordManager(subdomain, p.ip)
	if err != nil {
		return err
	}
	if !rm.IsManagedByUs(records) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	dnsRecords, hasTXT := rm.SeparateRecords(records)

	for _, record := range dnsRecords {
		if err := p.client.Records.Delete(
			ctx,
			domain,
			record.Name,
			powerdns.RRType(record.Type),
		); err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.Name, err)
		}
	}

	if hasTXT {
		if err := p.client.Records.Delete(
			ctx,
			domain,
			rm.MarkerName(),
			powerdns.RRTypeTXT,
		); err != nil {
			return fmt.Errorf("failed to delete TXT marker %s: %w", rm.MarkerName(), err)
		}
	}

	return nil
}

func (p *PowerDNSProvider) ListRecords(ctx context.Context, subdomain string) ([]DNSRecord, error) {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	records, err := p.client.Records.Get(ctx, domain, subdomain, nil)
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

	txtName := markerName(subdomain)
	txtRecords, err := p.client.Records.Get(ctx, domain, txtName, nil)
	if err == nil {
		for _, record := range txtRecords {
			dnsRecords = append(dnsRecords, DNSRecord{
				Name:    *record.Name,
				Type:    string(*record.Type),
				Content: *record.Records[0].Content,
			})
		}
	}

	return dnsRecords, nil
}

func (p *PowerDNSProvider) createTXTMarker(ctx context.Context, subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	// PowerDNS generally expects TXT content in quoted presentation format.
	return p.client.Records.Add(
		ctx,
		domain,
		markerName(subdomain),
		powerdns.RRTypeTXT,
		60,
		[]string{quoteTXT(managedTXT)},
	)
}
