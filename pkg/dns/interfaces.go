package dns

import (
	"fmt"
	"net"
)

type DNSProvider interface {
	CreateRecord(subdomain, ip string) error
	UpdateRecord(subdomain, ip string) error
	DeleteRecord(subdomain string) error
	ListRecords(subdomain string) ([]DNSRecord, error)
	CheckRecord(subdomain string) (bool, error)
}

type DNSRecord struct {
	Name  string
	Type  string
	Value string
}

func BoolPointer(b bool) *bool {
	return &b
}

func RecordType(ip string) (string, error) {
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if net.ParseIP(ip).To4() != nil {
		return "A", nil
	}
	return "AAAA", nil
}
