package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type DNSProviderService struct {
	app *config.App
}

func NewDNSProviderService(app *config.App) *DNSProviderService {
	return &DNSProviderService{app: app}
}

func (s *DNSProviderService) GetDNSProvider(
	ctx context.Context,
	req *mantraev1.GetDNSProviderRequest,
) (*mantraev1.GetDNSProviderResponse, error) {
	result, err := s.app.Conn.Q.GetDnsProvider(ctx, req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	decryptedAPIKey, err := util.DecryptSecret(result.Config.APIKey, s.app.Secret)
	if err != nil {
		return nil, err
	}
	result.Config.APIKey = decryptedAPIKey
	return &mantraev1.GetDNSProviderResponse{DnsProvider: result.ToProto()}, nil
}

func (s *DNSProviderService) CreateDNSProvider(
	ctx context.Context,
	req *mantraev1.CreateDNSProviderRequest,
) (*mantraev1.CreateDNSProviderResponse, error) {
	params := &db.CreateDnsProviderParams{
		ID:   uuid.New().String(),
		Name: req.Name,
		Type: int64(req.Type),
		Config: &schema.DNSProviderConfig{
			APIUrl:     req.Config.ApiUrl,
			IP:         req.Config.Ip,
			Proxied:    req.Config.Proxied,
			AutoUpdate: req.Config.AutoUpdate,
		},
		IsDefault: req.IsDefault,
	}
	if req.Config.ApiKey != "" {
		apiKeyHash, err := util.EncryptSecret(req.Config.ApiKey, s.app.Secret)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Config.APIKey = apiKeyHash
	}
	if req.IsDefault {
		if err := s.app.Conn.Q.UnsetDefaultDNSProvider(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.Q.CreateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &mantraev1.CreateDNSProviderResponse{DnsProvider: result.ToProto()}, nil
}

func (s *DNSProviderService) UpdateDNSProvider(
	ctx context.Context,
	req *mantraev1.UpdateDNSProviderRequest,
) (*mantraev1.UpdateDNSProviderResponse, error) {
	params := &db.UpdateDnsProviderParams{
		ID:   req.Id,
		Name: req.Name,
		Type: int64(req.Type),
		Config: &schema.DNSProviderConfig{
			APIUrl:     req.Config.ApiUrl,
			IP:         req.Config.Ip,
			Proxied:    req.Config.Proxied,
			AutoUpdate: req.Config.AutoUpdate,
		},
		IsDefault: req.IsDefault,
	}
	if req.Config.ApiKey != "" {
		apiKeyHash, err := util.EncryptSecret(req.Config.ApiKey, s.app.Secret)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Config.APIKey = apiKeyHash
	}
	if req.IsDefault {
		if err := s.app.Conn.Q.UnsetDefaultDNSProvider(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.Q.UpdateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &mantraev1.UpdateDNSProviderResponse{DnsProvider: result.ToProto()}, nil
}

func (s *DNSProviderService) DeleteDNSProvider(
	ctx context.Context,
	req *mantraev1.DeleteDNSProviderRequest,
) (*mantraev1.DeleteDNSProviderResponse, error) {
	if err := s.app.Conn.Q.DeleteDnsProvider(ctx, req.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.DeleteDNSProviderResponse{}, nil
}

func (s *DNSProviderService) ListDNSProviders(
	ctx context.Context,
	req *mantraev1.ListDNSProvidersRequest,
) (*mantraev1.ListDNSProvidersResponse, error) {
	params := &db.ListDnsProvidersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	result, err := s.app.Conn.Q.ListDnsProviders(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.Q.CountDnsProviders(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	dnsProviders := make([]*mantraev1.DNSProvider, 0, len(result))
	for _, p := range result {
		decryptedAPIKey, err := util.DecryptSecret(p.Config.APIKey, s.app.Secret)
		if err != nil {
			return nil, err
		}
		p.Config.APIKey = decryptedAPIKey
		dnsProviders = append(dnsProviders, p.ToProto())
	}
	return &mantraev1.ListDNSProvidersResponse{
		DnsProviders: dnsProviders,
		TotalCount:   totalCount,
	}, nil
}
