package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type DnsProviderService struct {
	app *config.App
}

func NewDnsProviderService(app *config.App) *DnsProviderService {
	return &DnsProviderService{app: app}
}

func (s *DnsProviderService) GetDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.GetDnsProviderRequest],
) (*connect.Response[mantraev1.GetDnsProviderResponse], error) {
	dnsProvider, err := s.app.Conn.GetQuery().GetDnsProvider(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetDnsProviderResponse{
		DnsProvider: &mantraev1.DnsProvider{
			Id:   dnsProvider.ID,
			Name: dnsProvider.Name,
			Type: dnsProvider.Type,
			Config: &mantraev1.DnsProviderConfig{
				ApiKey:     dnsProvider.Config.APIKey,
				ApiUrl:     dnsProvider.Config.APIUrl,
				Ip:         dnsProvider.Config.IP,
				Proxied:    dnsProvider.Config.Proxied,
				AutoUpdate: dnsProvider.Config.AutoUpdate,
				ZoneType:   dnsProvider.Config.ZoneType,
			},
			IsActive:  dnsProvider.IsActive,
			CreatedAt: SafeTimestamp(dnsProvider.CreatedAt),
			UpdatedAt: SafeTimestamp(dnsProvider.UpdatedAt),
		},
	}), nil
}

func (s *DnsProviderService) CreateDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateDnsProviderRequest],
) (*connect.Response[mantraev1.CreateDnsProviderResponse], error) {
	params := db.CreateDnsProviderParams{
		Name: req.Msg.Name,
		Type: req.Msg.Type,
		Config: &schema.DNSProviderConfig{
			APIKey:     req.Msg.Config.ApiKey,
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
			ZoneType:   req.Msg.Config.ZoneType,
		},
		IsActive: req.Msg.IsActive,
	}

	dnsProvider, err := s.app.Conn.GetQuery().CreateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateDnsProviderResponse{
		DnsProvider: &mantraev1.DnsProvider{
			Id:   dnsProvider.ID,
			Name: dnsProvider.Name,
			Type: dnsProvider.Type,
			Config: &mantraev1.DnsProviderConfig{
				ApiKey:     dnsProvider.Config.APIKey,
				ApiUrl:     dnsProvider.Config.APIUrl,
				Ip:         dnsProvider.Config.IP,
				Proxied:    dnsProvider.Config.Proxied,
				AutoUpdate: dnsProvider.Config.AutoUpdate,
				ZoneType:   dnsProvider.Config.ZoneType,
			},
			IsActive:  dnsProvider.IsActive,
			CreatedAt: SafeTimestamp(dnsProvider.CreatedAt),
			UpdatedAt: SafeTimestamp(dnsProvider.UpdatedAt),
		},
	}), nil
}

func (s *DnsProviderService) UpdateDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateDnsProviderRequest],
) (*connect.Response[mantraev1.UpdateDnsProviderResponse], error) {
	params := db.UpdateDnsProviderParams{
		Name: req.Msg.Name,
		Type: req.Msg.Type,
		Config: &schema.DNSProviderConfig{
			APIKey:     req.Msg.Config.ApiKey,
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
			ZoneType:   req.Msg.Config.ZoneType,
		},
		IsActive: req.Msg.IsActive,
	}

	dnsProvider, err := s.app.Conn.GetQuery().UpdateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateDnsProviderResponse{
		DnsProvider: &mantraev1.DnsProvider{
			Id:   dnsProvider.ID,
			Name: dnsProvider.Name,
			Type: dnsProvider.Type,
			Config: &mantraev1.DnsProviderConfig{
				ApiKey:     dnsProvider.Config.APIKey,
				ApiUrl:     dnsProvider.Config.APIUrl,
				Ip:         dnsProvider.Config.IP,
				Proxied:    dnsProvider.Config.Proxied,
				AutoUpdate: dnsProvider.Config.AutoUpdate,
				ZoneType:   dnsProvider.Config.ZoneType,
			},
			IsActive:  dnsProvider.IsActive,
			CreatedAt: SafeTimestamp(dnsProvider.CreatedAt),
			UpdatedAt: SafeTimestamp(dnsProvider.UpdatedAt),
		},
	}), nil
}

func (s *DnsProviderService) DeleteDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteDnsProviderRequest],
) (*connect.Response[mantraev1.DeleteDnsProviderResponse], error) {
	err := s.app.Conn.GetQuery().DeleteDnsProvider(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteDnsProviderResponse{}), nil
}

func (s *DnsProviderService) ListDnsProviders(
	ctx context.Context,
	req *connect.Request[mantraev1.ListDnsProvidersRequest],
) (*connect.Response[mantraev1.ListDnsProvidersResponse], error) {
	var params db.ListDnsProvidersParams
	if req.Msg.Limit == nil {
		params.Limit = 100
	} else {
		params.Limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *req.Msg.Offset
	}

	dbDNSProviders, err := s.app.Conn.GetQuery().ListDnsProviders(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountDnsProviders(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var dnsProviders []*mantraev1.DnsProvider
	for _, dnsProvider := range dbDNSProviders {
		dnsProviders = append(dnsProviders, &mantraev1.DnsProvider{
			Id:   dnsProvider.ID,
			Name: dnsProvider.Name,
			Type: dnsProvider.Type,
			Config: &mantraev1.DnsProviderConfig{
				ApiKey:     dnsProvider.Config.APIKey,
				ApiUrl:     dnsProvider.Config.APIUrl,
				Ip:         dnsProvider.Config.IP,
				Proxied:    dnsProvider.Config.Proxied,
				AutoUpdate: dnsProvider.Config.AutoUpdate,
				ZoneType:   dnsProvider.Config.ZoneType,
			},
			IsActive:  dnsProvider.IsActive,
			CreatedAt: SafeTimestamp(dnsProvider.CreatedAt),
			UpdatedAt: SafeTimestamp(dnsProvider.UpdatedAt),
		})
	}
	return connect.NewResponse(&mantraev1.ListDnsProvidersResponse{
		DnsProviders: dnsProviders,
		TotalCount:   totalCount,
	}), nil
}
