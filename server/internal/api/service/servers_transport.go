package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type ServersTransportService struct {
	app      *config.App
	dispatch map[mantraev1.ProtocolType]ServersTransportOps
}

func NewServersTransportService(app *config.App) *ServersTransportService {
	return &ServersTransportService{
		app: app,
		dispatch: map[mantraev1.ProtocolType]ServersTransportOps{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP: NewHTTPServersTransportOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP:  NewTCPServersTransportOps(app),
		},
	}
}

func (s *ServersTransportService) GetServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.GetServersTransportRequest],
) (*connect.Response[mantraev1.GetServersTransportResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid servers transport type"),
		)
	}

	result, err := ops.Get(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(result), nil
}

func (s *ServersTransportService) CreateServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateServersTransportRequest],
) (*connect.Response[mantraev1.CreateServersTransportResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid servers transport type"),
		)
	}

	result, err := ops.Create(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(result), nil
}

func (s *ServersTransportService) UpdateServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateServersTransportRequest],
) (*connect.Response[mantraev1.UpdateServersTransportResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid servers transport type"),
		)
	}

	result, err := ops.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(result), nil
}

func (s *ServersTransportService) DeleteServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteServersTransportRequest],
) (*connect.Response[mantraev1.DeleteServersTransportResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid servers transport type"),
		)
	}

	result, err := ops.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(result), nil
}

func (s *ServersTransportService) ListServersTransports(
	ctx context.Context,
	req *connect.Request[mantraev1.ListServersTransportsRequest],
) (*connect.Response[mantraev1.ListServersTransportsResponse], error) {
	if req.Msg.Type != nil {
		ops, ok := s.dispatch[*req.Msg.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid servers transport type"),
			)
		}

		result, err := ops.List(ctx, req.Msg)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(result), nil
	} else {
		// Get HTTP servers transports
		httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
		httpResult, err := httpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, err
		}

		// Get TCP servers transports
		tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
		tcpResult, err := tcpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, err
		}

		// Combine results
		allServersTransports := append(httpResult.ServersTransports, tcpResult.ServersTransports...)
		totalCount := httpResult.TotalCount + tcpResult.TotalCount

		return connect.NewResponse(&mantraev1.ListServersTransportsResponse{
			ServersTransports: allServersTransports,
			TotalCount:        totalCount,
		}), nil
	}
}
