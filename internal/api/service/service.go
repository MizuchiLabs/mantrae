// Package service provides the gRPC service implementations.
package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
)

type Service struct {
	app      *config.App
	dispatch map[mantraev1.ProtocolType]ServiceOps
}

func NewServiceService(app *config.App) *Service {
	return &Service{
		app: app,
		dispatch: map[mantraev1.ProtocolType]ServiceOps{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP: NewHTTPServiceOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP:  NewTCPServiceOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_UDP:  NewUDPServiceOps(app),
		},
	}
}

func (s *Service) GetService(
	ctx context.Context,
	req *connect.Request[mantraev1.GetServiceRequest],
) (*connect.Response[mantraev1.GetServiceResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Get(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *Service) CreateService(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateServiceRequest],
) (*connect.Response[mantraev1.CreateServiceResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Create(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *Service) UpdateService(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateServiceRequest],
) (*connect.Response[mantraev1.UpdateServiceResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Update(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *Service) DeleteService(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteServiceRequest],
) (*connect.Response[mantraev1.DeleteServiceResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Delete(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *Service) ListServices(
	ctx context.Context,
	req *connect.Request[mantraev1.ListServicesRequest],
) (*connect.Response[mantraev1.ListServicesResponse], error) {
	if req.Msg.Type != nil {
		ops, ok := s.dispatch[*req.Msg.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid service type"),
			)
		}

		result, err := ops.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return connect.NewResponse(result), nil
	} else {
		// Get HTTP services
		httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
		httpResult, err := httpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Get TCP services
		tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
		tcpResult, err := tcpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Get UDP services
		udpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_UDP]
		udpResult, err := udpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Combine results
		allServices := append(httpResult.Services, tcpResult.Services...)
		allServices = append(allServices, udpResult.Services...)
		totalCount := httpResult.TotalCount + tcpResult.TotalCount + udpResult.TotalCount

		return connect.NewResponse(&mantraev1.ListServicesResponse{
			Services:   allServices,
			TotalCount: totalCount,
		}), nil
	}
}
