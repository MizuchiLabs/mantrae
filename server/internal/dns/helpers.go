package dns

import (
	"fmt"
	"strings"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

const managedTXT = "managed-by=mantrae"

// RecordManager handles common DNS operations
type RecordManager struct {
	subdomain  string
	externalIP string
	recordType string
}

// UpsertOperation defines the required operations for upserting
type UpsertOperation struct {
	UpdateDNSRecord func(recordID, recordType string) error
	CreateDNSRecord func(recordType string) error
	CreateTXTMarker func() error
}

func NewRecordManager(subdomain, externalIP string) (*RecordManager, error) {
	recordType, err := determineRecordType(externalIP)
	if err != nil {
		return nil, err
	}
	return &RecordManager{
		subdomain:  subdomain,
		externalIP: externalIP,
		recordType: recordType,
	}, nil
}

func (rm *RecordManager) MarkerName() string {
	return markerName(rm.subdomain)
}

func (rm *RecordManager) RecordType() string {
	return rm.recordType
}

func (rm *RecordManager) IP() string {
	return rm.externalIP
}

// IsManagedByUs checks if TXT marker exists in records
func (rm *RecordManager) IsManagedByUs(records []DNSRecord) bool {
	marker := rm.MarkerName()
	for _, record := range records {
		if record.Name == marker && record.Type == "TXT" &&
			normalizeTXT(record.Content) == managedTXT {
			return true
		}
	}
	return false
}

// NeedsUpdate checks if records need updating
func (rm *RecordManager) NeedsUpdate(records []DNSRecord) bool {
	hasCorrectRecord := false
	for _, record := range records {
		if record.Type == "TXT" {
			continue
		}
		if record.Type != rm.recordType {
			return true
		}
		if record.Content == rm.externalIP {
			hasCorrectRecord = true
		} else {
			return true
		}
	}
	return !hasCorrectRecord
}

// SeparateRecords splits DNS records from TXT marker
func (rm *RecordManager) SeparateRecords(
	records []DNSRecord,
) (dnsRecords []DNSRecord, hasTXT bool) {
	marker := rm.MarkerName()
	for _, record := range records {
		if record.Type == "TXT" && record.Name == marker {
			hasTXT = true
		} else if record.Type == "A" || record.Type == "AAAA" {
			dnsRecords = append(dnsRecords, record)
		}
	}
	return dnsRecords, hasTXT
}

// ExecuteUpsert contains the common upsert logic
func (rm *RecordManager) ExecuteUpsert(records []DNSRecord, ops UpsertOperation) error {
	if len(records) > 0 && !rm.IsManagedByUs(records) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	dnsRecords, hasTXT := rm.SeparateRecords(records)

	if len(dnsRecords) == 0 {
		if err := ops.CreateDNSRecord(rm.recordType); err != nil {
			return err
		}
		if err := ops.CreateTXTMarker(); err != nil {
			return fmt.Errorf("failed to create TXT marker: %w", err)
		}
		return nil
	}

	if rm.NeedsUpdate(dnsRecords) {
		for _, record := range dnsRecords {
			if err := ops.UpdateDNSRecord(record.ID, rm.recordType); err != nil {
				return err
			}
		}
	}

	if !hasTXT {
		if err := ops.CreateTXTMarker(); err != nil {
			return fmt.Errorf("failed to create TXT marker: %w", err)
		}
	}

	return nil
}

func determineRecordType(ip string) (string, error) {
	if util.IsValidIPv4(ip) {
		return "A", nil
	}
	if util.IsValidIPv6(ip) {
		return "AAAA", nil
	}
	return "", fmt.Errorf("invalid IP address: %s", ip)
}

func markerName(subdomain string) string {
	return "_mantrae." + subdomain
}

func normalizeTXT(s string) string {
	return strings.Trim(s, "\"")
}

// Some DNS backends expect TXT in presentation format (quoted).
func quoteTXT(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s
	}
	return `"` + s + `"`
}
