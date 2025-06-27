package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"

	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/pkg/util"
)

type TechnitiumProvider struct {
	BaseURL    string
	APIKey     string
	ExternalIP string
	ZoneType   string // primary, forwarder
}

func NewTechnitiumProvider(d *schema.DNSProviderConfig) *TechnitiumProvider {
	if !slices.Contains(ZoneTypes, d.ZoneType) {
		slog.Error("Invalid zone type", "type", d.ZoneType)
	}
	return &TechnitiumProvider{
		BaseURL:    d.APIUrl,
		APIKey:     d.APIKey,
		ExternalIP: d.IP,
		ZoneType:   d.ZoneType,
	}
}

// Generic HTTP request helper
func (t *TechnitiumProvider) doRequest(
	ctx context.Context,
	method, endpoint string,
	body any,
) (*http.Response, error) {
	url := t.BaseURL + endpoint

	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (t *TechnitiumProvider) UpsertRecord(subdomain string) error {
	var recordType string
	if util.IsValidIPv4(t.ExternalIP) {
		recordType = "A"
	} else if util.IsValidIPv6(t.ExternalIP) {
		recordType = "AAAA"
	} else {
		return fmt.Errorf("invalid IP address: %s", t.ExternalIP)
	}

	// Check if the record exists
	records, err := t.ListRecords(subdomain)
	if err != nil {
		return err
	}

	if err := t.checkRecord(subdomain); err != nil {
		return err
	}

	shouldUpdate := verifyRecords(records, subdomain, t.ExternalIP)
	if len(records) <= 1 { // At least 2 records must exist TXT+A/AAAA
		if err := t.createRecord(subdomain, recordType); err != nil {
			return err
		}
		slog.Info("Created record", "name", subdomain, "type", recordType, "content", t.ExternalIP)
	} else if shouldUpdate {
		for _, record := range records {
			if record.Type != "TXT" {
				if err := t.updateRecord(subdomain, recordType); err != nil {
					return err
				}
				slog.Info("Updated record", "name", record.Name, "type", record.Type, "content", record.Content)
			}
		}
	}

	return nil
}

func (t *TechnitiumProvider) DeleteRecord(subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	if err = t.checkRecord(subdomain); err != nil {
		return err
	}

	records, err := t.ListRecords(subdomain)
	if err != nil {
		return err
	}

	for _, record := range records {
		endpoint := fmt.Sprintf(
			"/api/zones/records/delete?token=%s&zone=%s&type=%s",
			t.APIKey,
			domain,
			record.Type,
		)

		if t.ZoneType == "forwarder" {
			endpoint = endpoint + "&forwarder=true"
		}
		if record.Type == "A" || record.Type == "AAAA" {
			endpoint = endpoint + "&domain=" + subdomain + "&ipAddress=" + record.Content
		}
		if record.Type == "TXT" {
			endpoint = endpoint + "&domain=_mantrae-" + subdomain + "&text=" + url.QueryEscape(
				record.Content,
			)
		}

		resp, err := t.doRequest(
			context.Background(),
			http.MethodPost,
			endpoint,
			nil,
		)
		if err != nil {
			return err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Error("failed to close response body", "error", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete record: %s", string(bodyBytes))
		}

		if record.Type != "TXT" {
			slog.Info(
				"Deleted record",
				"name",
				record.Name,
				"type",
				record.Type,
				"content",
				record.Content,
			)
		}
	}

	return nil
}

func (t *TechnitiumProvider) createRecord(subdomain, recordType string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/add?token=%s&zone=%s",
		t.APIKey,
		domain,
	)

	if t.ZoneType == "forwarder" {
		endpoint = endpoint + "&forwarder=this-server"
	}

	// Create the A/AAAA record
	resp, err := t.doRequest(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf(
			"%s&type=%s&domain=%s&ipAddress=%s",
			endpoint,
			recordType,
			subdomain,
			t.ExternalIP,
		),
		nil,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record: %s", string(bodyBytes))
	}

	// Create the TXT record
	resp, err = t.doRequest(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf(
			"%s&type=TXT&domain=_mantrae-%s&text=%s",
			endpoint,
			subdomain,
			url.QueryEscape(managedTXT),
		),
		nil,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record: %s", string(bodyBytes))
	}
	return nil
}

func (t *TechnitiumProvider) updateRecord(subdomain, recordType string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/update?token=%s&zone=%s&type=%s&ipAddress=%s",
		t.APIKey,
		domain,
		recordType,
		t.ExternalIP,
	)
	if t.ZoneType == "forwarder" {
		endpoint = endpoint + "&forwarder=this-server"
	}

	resp, err := t.doRequest(
		context.Background(),
		http.MethodPatch,
		endpoint,
		nil,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update record: %s", string(bodyBytes))
	}
	return nil
}

func (t *TechnitiumProvider) ListRecords(subdomain string) ([]DNSRecord, error) {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/get?token=%s&domain=%s&zone=%s&listZone=true",
		t.APIKey,
		subdomain,
		domain,
	)

	resp, err := t.doRequest(
		context.Background(),
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	var tRecords struct {
		Status       string `json:"status"`
		ErrorMessage string `json:"errorMessage"`
		Response     struct {
			Records []struct {
				Name  string `json:"name"`
				Type  string `json:"type"`
				RData struct {
					IP   string `json:"ipAddress"`
					Text string `json:"text"`
				}
			}
		} `json:"response"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&tRecords); err != nil {
		return nil, err
	}

	if tRecords.Status == "error" {
		return nil, fmt.Errorf("%s", tRecords.ErrorMessage)
	}

	var records []DNSRecord
	for _, record := range tRecords.Response.Records {
		if record.Name == "_mantrae-"+subdomain && record.Type == "TXT" &&
			record.RData.Text == managedTXT {
			records = append(records, DNSRecord{
				Name:    record.Name,
				Type:    record.Type,
				Content: record.RData.Text,
			})
		}
		if record.Name == subdomain && record.Type == "A" {
			records = append(records, DNSRecord{
				Name:    record.Name,
				Type:    record.Type,
				Content: record.RData.IP,
			})
		}
		if record.Name == subdomain && record.Type == "AAAA" {
			records = append(records, DNSRecord{
				Name:    record.Name,
				Type:    record.Type,
				Content: record.RData.IP,
			})
		}
	}

	return records, nil
}

func (t *TechnitiumProvider) checkRecord(subdomain string) error {
	records, err := t.ListRecords(subdomain)
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
