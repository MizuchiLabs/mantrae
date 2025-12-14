package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
)

type DNSProviderService struct {
	app *config.App
}

func NewDNSProviderService(app *config.App) *DNSProviderService {
	return &DNSProviderService{app: app}
}

func (s *DNSProviderService) GetDNSProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.GetDNSProviderRequest],
) (*connect.Response[mantraev1.GetDNSProviderResponse], error) {
	result, err := s.app.Conn.GetQuery().GetDnsProvider(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	decryptedAPIKey, err := util.DecryptSecret(result.Config.APIKey, s.app.Secret)
	if err != nil {
		return nil, err
	}
	result.Config.APIKey = decryptedAPIKey
	return connect.NewResponse(&mantraev1.GetDNSProviderResponse{
		DnsProvider: result.ToProto(),
	}), nil
}

func (s *DNSProviderService) CreateDNSProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateDNSProviderRequest],
) (*connect.Response[mantraev1.CreateDNSProviderResponse], error) {
	params := &db.CreateDnsProviderParams{
		ID:   uuid.New().String(),
		Name: req.Msg.Name,
		Type: int64(req.Msg.Type),
		Config: &schema.DNSProviderConfig{
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
		},
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.Config.ApiKey != "" {
		apiKeyHash, err := util.EncryptSecret(req.Msg.Config.ApiKey, s.app.Secret)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Config.APIKey = apiKeyHash
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

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_DnsProvider{
			DnsProvider: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.CreateDNSProviderResponse{
		DnsProvider: result.ToProto(),
	}), nil
}

func (s *DNSProviderService) UpdateDNSProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateDNSProviderRequest],
) (*connect.Response[mantraev1.UpdateDNSProviderResponse], error) {
	params := &db.UpdateDnsProviderParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
		Type: int64(req.Msg.Type),
		Config: &schema.DNSProviderConfig{
			APIUrl:     req.Msg.Config.ApiUrl,
			IP:         req.Msg.Config.Ip,
			Proxied:    req.Msg.Config.Proxied,
			AutoUpdate: req.Msg.Config.AutoUpdate,
		},
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.Config.ApiKey != "" {
		apiKeyHash, err := util.EncryptSecret(req.Msg.Config.ApiKey, s.app.Secret)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Config.APIKey = apiKeyHash
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

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_DnsProvider{
			DnsProvider: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.UpdateDNSProviderResponse{
		DnsProvider: result.ToProto(),
	}), nil
}

func (s *DNSProviderService) DeleteDNSProvider(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteDNSProviderRequest],
) (*connect.Response[mantraev1.DeleteDNSProviderResponse], error) {
	dnsProvider, err := s.app.Conn.GetQuery().GetDnsProvider(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.app.Conn.GetQuery().DeleteDnsProvider(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_DnsProvider{
			DnsProvider: dnsProvider.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.DeleteDNSProviderResponse{}), nil
}

func (s *DNSProviderService) ListDNSProviders(
	ctx context.Context,
	req *connect.Request[mantraev1.ListDNSProvidersRequest],
) (*connect.Response[mantraev1.ListDNSProvidersResponse], error) {
	params := &db.ListDnsProvidersParams{
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	}

	result, err := s.app.Conn.GetQuery().ListDnsProviders(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountDnsProviders(ctx)
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
	return connect.NewResponse(&mantraev1.ListDNSProvidersResponse{
		DnsProviders: dnsProviders,
		TotalCount:   totalCount,
	}), nil
}
