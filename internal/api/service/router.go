package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
)

type RouterService struct {
	app      *config.App
	dispatch map[mantraev1.ProtocolType]RouterOps
}

func NewRouterService(app *config.App) *RouterService {
	return &RouterService{
		app: app,
		dispatch: map[mantraev1.ProtocolType]RouterOps{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP: NewHTTPRouterOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP:  NewTCPRouterOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_UDP:  NewUDPRouterOps(app),
		},
	}
}

func (s *RouterService) GetRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.GetRouterRequest],
) (*connect.Response[mantraev1.GetRouterResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid router type"),
		)
	}

	result, err := ops.Get(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *RouterService) CreateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateRouterRequest],
) (*connect.Response[mantraev1.CreateRouterResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid router type"),
		)
	}

	result, err := ops.Create(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *RouterService) UpdateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateRouterRequest],
) (*connect.Response[mantraev1.UpdateRouterResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid router type"),
		)
	}

	result, err := ops.Update(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *RouterService) DeleteRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteRouterRequest],
) (*connect.Response[mantraev1.DeleteRouterResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid router type"),
		)
	}

	result, err := ops.Delete(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *RouterService) ListRouters(
	ctx context.Context,
	req *connect.Request[mantraev1.ListRoutersRequest],
) (*connect.Response[mantraev1.ListRoutersResponse], error) {
	if req.Msg.Type != nil {
		ops, ok := s.dispatch[*req.Msg.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid router type"),
			)
		}

		result, err := ops.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return connect.NewResponse(result), nil
	} else {
		// Get HTTP routers
		httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
		httpResult, err := httpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Get TCP routers
		tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
		tcpResult, err := tcpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Get UDP routers
		udpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_UDP]
		udpResult, err := udpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Combine results
		allRouters := append(httpResult.Routers, tcpResult.Routers...)
		allRouters = append(allRouters, udpResult.Routers...)
		totalCount := httpResult.TotalCount + tcpResult.TotalCount + udpResult.TotalCount

		return connect.NewResponse(&mantraev1.ListRoutersResponse{
			Routers:    allRouters,
			TotalCount: totalCount,
		}), nil
	}
}
