package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
)

type CloudflareProvider struct {
	APIKey     string
	ExternalIP string
	Proxied    bool
	BaseURL    string
}

func NewCloudflareProvider(key, ip string, proxied bool) *CloudflareProvider {
	return &CloudflareProvider{
		APIKey:     key,
		ExternalIP: ip,
		Proxied:    proxied,
		BaseURL:    "https://api.cloudflare.com/client/v4",
	}
}

// Generic HTTP request helper
func (c *CloudflareProvider) doRequest(
	ctx context.Context,
	method, endpoint string,
	body interface{},
) (*http.Response, error) {
	url := c.BaseURL + endpoint

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

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (c *CloudflareProvider) UpsertRecord(subdomain string) error {
	recordType := "A"
	if net.ParseIP(c.ExternalIP).To4() == nil {
		recordType = "AAAA"
	}

	zoneID, err := c.getZoneID(subdomain)
	if err != nil {
		return err
	}

	// Check if the record exists
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	shouldUpdate := verifyRecords(records, subdomain, c.ExternalIP)
	if len(records) <= 1 {
		if err := c.createRecord(zoneID, subdomain, recordType); err != nil {
			return err
		}
		slog.Info("Created record", "name", subdomain, "type", recordType, "content", c.ExternalIP)
	} else if shouldUpdate {
		for _, record := range records {
			if record.Type != "TXT" {
				if err := c.updateRecord(zoneID, record.ID, recordType, subdomain); err != nil {
					return err
				}
				slog.Info("Updated record", "name", record.Name, "type", record.Type, "content", record.Content)
			}
		}
	}

	return nil
}

func (c *CloudflareProvider) createRecord(zoneID, subdomain, recordType string) error {
	if !c.CheckRecord(subdomain) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	// Create the A/AAAA record
	body := map[string]interface{}{
		"type":    recordType,
		"name":    subdomain,
		"content": c.ExternalIP,
		"proxied": c.Proxied,
	}
	resp, err := c.doRequest(
		context.Background(),
		http.MethodPost,
		"/zones/"+zoneID+"/dns_records",
		body,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record: %s", string(bodyBytes))
	}

	// Create the TXT record
	body = map[string]interface{}{
		"type":    "TXT",
		"name":    "_mantrae-" + subdomain,
		"content": ManagedTXT,
	}
	resp, err = c.doRequest(
		context.Background(),
		http.MethodPost,
		"/zones/"+zoneID+"/dns_records",
		body,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create record: %s", string(bodyBytes))
	}
	return nil
}

func (c *CloudflareProvider) updateRecord(zoneID, recordID, recordType, subdomain string) error {
	if !c.CheckRecord(subdomain) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	body := map[string]interface{}{
		"type":    recordType,
		"name":    subdomain,
		"content": c.ExternalIP,
		"proxied": c.Proxied,
	}
	resp, err := c.doRequest(
		context.Background(),
		http.MethodPatch,
		"/zones/"+zoneID+"/dns_records/"+recordID,
		body,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update record: %s", string(bodyBytes))
	}
	return nil
}

func (c *CloudflareProvider) ListRecords(subdomain string) ([]DNSRecord, error) {
	zoneID, err := c.getZoneID(subdomain)
	if err != nil {
		return nil, err
	}

	// List A/AAAA/TXT records
	var records []DNSRecord
	searchRecords := []string{subdomain, "_mantrae-" + subdomain}
	for _, record := range searchRecords {
		resp, err := c.doRequest(
			context.Background(),
			http.MethodGet,
			"/zones/"+zoneID+"/dns_records?name="+record,
			nil,
		)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Decode response into records...
		var cfRecords struct {
			Records []DNSRecord `json:"result"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&cfRecords); err != nil {
			return nil, err
		}
		records = append(records, cfRecords.Records...)
	}
	return records, nil
}

func (c *CloudflareProvider) getZoneID(subdomain string) (string, error) {
	baseDomain := getBaseDomain(subdomain)
	resp, err := c.doRequest(context.Background(), http.MethodGet, "/zones?name="+baseDomain, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil || len(result.Result) == 0 {
		return "", fmt.Errorf("could not find zone for domain %s", baseDomain)
	}

	return result.Result[0].ID, nil
}

func (c *CloudflareProvider) DeleteRecord(subdomain string) error {
	if !c.CheckRecord(subdomain) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	zoneID, err := c.getZoneID(subdomain)
	if err != nil {
		return err
	}

	records, err := c.ListRecords(subdomain)
	if err != nil {
		return err
	}

	for _, record := range records {
		resp, err := c.doRequest(
			context.Background(),
			http.MethodDelete,
			"/zones/"+zoneID+"/dns_records/"+record.ID,
			nil,
		)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete record: %s", string(bodyBytes))
		}

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

	return nil
}

func (c *CloudflareProvider) CheckRecord(subdomain string) bool {
	records, err := c.ListRecords(subdomain)
	if err != nil {
		return false
	}

	if len(records) == 0 {
		return true
	}

	for _, record := range records {
		if record.Name == "_mantrae-"+subdomain && record.Type == "TXT" &&
			record.Content == ManagedTXT {
			return true
		}
	}

	return false
}
