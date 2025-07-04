package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
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
	result, err := s.app.Conn.GetQuery().GetDnsProvider(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.GetDnsProviderResponse{
		DnsProvider: convert.DNSProviderToProto(&result),
	}), nil
}

func (s *DnsProviderService) CreateDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateDnsProviderRequest],
) (*connect.Response[mantraev1.CreateDnsProviderResponse], error) {
	var dnsType string
	switch req.Msg.Type {
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_CLOUDFLARE:
		dnsType = "cloudflare"
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_POWERDNS:
		dnsType = "powerdns"
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_TECHNITIUM:
		dnsType = "technitium"
	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid dns provider type"),
		)
	}

	params := db.CreateDnsProviderParams{
		Name: req.Msg.Name,
		Type: dnsType,
		Config: &schema.DNSProviderConfig{
			APIKey:     req.Msg.Config.ApiKey,
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
			ZoneType:   req.Msg.Config.ZoneType,
		},
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.IsDefault {
		if err := s.app.Conn.GetQuery().UnsetDefaultDNSProvider(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.GetQuery().CreateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateDnsProviderResponse{
		DnsProvider: convert.DNSProviderToProto(&result),
	}), nil
}

func (s *DnsProviderService) UpdateDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateDnsProviderRequest],
) (*connect.Response[mantraev1.UpdateDnsProviderResponse], error) {
	var dnsType string
	switch req.Msg.Type {
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_CLOUDFLARE:
		dnsType = "cloudflare"
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_POWERDNS:
		dnsType = "powerdns"
	case mantraev1.DnsProviderType_DNS_PROVIDER_TYPE_TECHNITIUM:
		dnsType = "technitium"
	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid dns provider type"),
		)
	}

	params := db.UpdateDnsProviderParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
		Type: dnsType,
		Config: &schema.DNSProviderConfig{
			APIKey:     req.Msg.Config.ApiKey,
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
			ZoneType:   req.Msg.Config.ZoneType,
		},
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.IsDefault {
		if err := s.app.Conn.GetQuery().UnsetDefaultDNSProvider(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.GetQuery().UpdateDnsProvider(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateDnsProviderResponse{
		DnsProvider: convert.DNSProviderToProto(&result),
	}), nil
}

func (s *DnsProviderService) DeleteDnsProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteDnsProviderRequest],
) (*connect.Response[mantraev1.DeleteDnsProviderResponse], error) {
	if err := s.app.Conn.GetQuery().DeleteDnsProvider(ctx, req.Msg.Id); err != nil {
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

	result, err := s.app.Conn.GetQuery().ListDnsProviders(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountDnsProviders(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.ListDnsProvidersResponse{
		DnsProviders: convert.DNSProvidersToProto(result),
		TotalCount:   totalCount,
	}), nil
}
