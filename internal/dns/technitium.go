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

	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type TechnitiumProvider struct {
	baseURL string
	apiKey  string
	ip      string
	client  *http.Client
}

func NewTechnitiumProvider(d *mantraev1.DNSProviderConfig) *TechnitiumProvider {
	return &TechnitiumProvider{
		baseURL: d.ApiUrl,
		apiKey:  d.ApiKey,
		ip:      d.Ip,
		client:  http.DefaultClient,
	}
}

func (t *TechnitiumProvider) doRequest(
	ctx context.Context,
	method, endpoint string,
	body any,
) (*http.Response, error) {
	fullURL := t.baseURL + endpoint

	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return t.client.Do(req)
}

func (t *TechnitiumProvider) UpsertRecord(ctx context.Context, subdomain string) error {
	rm, err := NewRecordManager(subdomain, t.ip)
	if err != nil {
		return err
	}

	records, err := t.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}

	ops := UpsertOperation{
		CreateDNSRecord: func(recordType string) error {
			return t.createRecord(ctx, subdomain, recordType)
		},
		CreateTXTMarker: func() error {
			return t.createTXTMarker(ctx, subdomain)
		},
		UpdateDNSRecord: func(_ string, recordType string) error {
			return t.updateRecord(ctx, subdomain, recordType)
		},
	}

	return rm.ExecuteUpsert(records, ops)
}

func (t *TechnitiumProvider) DeleteRecord(ctx context.Context, subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	records, err := t.ListRecords(ctx, subdomain)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	rm, err := NewRecordManager(subdomain, t.ip)
	if err != nil {
		return err
	}
	if !rm.IsManagedByUs(records) {
		return fmt.Errorf("record not managed by Mantrae")
	}

	for _, record := range records {
		endpoint := fmt.Sprintf(
			"/api/zones/records/delete?token=%s&zone=%s&type=%s",
			t.apiKey,
			domain,
			record.Type,
		)

		switch record.Type {
		case "A", "AAAA":
			endpoint += "&domain=" + subdomain + "&ipAddress=" + record.Content
		case "TXT":
			endpoint += "&domain=" + markerName(
				subdomain,
			) + "&text=" + url.QueryEscape(
				record.Content,
			)
		}

		resp, err := t.doRequest(ctx, http.MethodPost, endpoint, nil)
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
	}

	return nil
}

func (t *TechnitiumProvider) createRecord(ctx context.Context, subdomain, recordType string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/api/zones/records/add?token=%s&zone=%s", t.apiKey, domain)

	resp, err := t.doRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"%s&type=%s&domain=%s&ipAddress=%s",
			endpoint,
			recordType,
			subdomain,
			t.ip,
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

	return nil
}

func (t *TechnitiumProvider) updateRecord(ctx context.Context, subdomain, recordType string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/update?token=%s&zone=%s&type=%s&ipAddress=%s",
		t.apiKey,
		domain,
		recordType,
		t.ip,
	)

	resp, err := t.doRequest(ctx, http.MethodPatch, endpoint, nil)
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

func (t *TechnitiumProvider) ListRecords(
	ctx context.Context,
	subdomain string,
) ([]DNSRecord, error) {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/get?token=%s&domain=%s&zone=%s&listZone=true",
		t.apiKey,
		subdomain,
		domain,
	)

	resp, err := t.doRequest(ctx, http.MethodGet, endpoint, nil)
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

	marker := markerName(subdomain)

	var records []DNSRecord
	for _, record := range tRecords.Response.Records {
		if record.Name == marker && record.Type == "TXT" &&
			normalizeTXT(record.RData.Text) == managedTXT {
			records = append(records, DNSRecord{
				Name:    record.Name,
				Type:    record.Type,
				Content: record.RData.Text,
			})
			continue
		}

		if record.Name == subdomain && (record.Type == "A" || record.Type == "AAAA") {
			records = append(records, DNSRecord{
				Name:    record.Name,
				Type:    record.Type,
				Content: record.RData.IP,
			})
		}
	}

	return records, nil
}

func (t *TechnitiumProvider) createTXTMarker(ctx context.Context, subdomain string) error {
	domain, err := util.ExtractBaseDomain(subdomain)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(
		"/api/zones/records/add?token=%s&zone=%s&type=TXT&domain=%s&text=%s",
		t.apiKey,
		domain,
		markerName(subdomain),
		url.QueryEscape(managedTXT),
	)

	resp, err := t.doRequest(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			slog.Error("failed to close response body", "error", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create TXT marker: %s", string(bodyBytes))
	}
	return nil
}
