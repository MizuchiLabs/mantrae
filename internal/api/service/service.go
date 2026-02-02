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
	req *mantraev1.GetServiceRequest,
) (*mantraev1.GetServiceResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Get(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *Service) CreateService(
	ctx context.Context,
	req *mantraev1.CreateServiceRequest,
) (*mantraev1.CreateServiceResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Create(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *Service) UpdateService(
	ctx context.Context,
	req *mantraev1.UpdateServiceRequest,
) (*mantraev1.UpdateServiceResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Update(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *Service) DeleteService(
	ctx context.Context,
	req *mantraev1.DeleteServiceRequest,
) (*mantraev1.DeleteServiceResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid service type"),
		)
	}

	result, err := ops.Delete(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *Service) ListServices(
	ctx context.Context,
	req *mantraev1.ListServicesRequest,
) (*mantraev1.ListServicesResponse, error) {
	if req.Type != nil {
		ops, ok := s.dispatch[*req.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid service type"),
			)
		}

		result, err := ops.List(ctx, req)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return result, nil
	}

	// Get HTTP services
	httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
	httpResult, err := httpOps.List(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Get TCP services
	tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
	tcpResult, err := tcpOps.List(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Get UDP services
	udpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_UDP]
	udpResult, err := udpOps.List(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Combine results
	allServices := append(httpResult.Services, tcpResult.Services...)
	allServices = append(allServices, udpResult.Services...)
	totalCount := httpResult.TotalCount + tcpResult.TotalCount + udpResult.TotalCount

	return &mantraev1.ListServicesResponse{
		Services:   allServices,
		TotalCount: totalCount,
	}, nil
}
