package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type RouterService struct {
	app *config.App
}

func NewRouterService(app *config.App) *RouterService {
	return &RouterService{app: app}
}

func (s *RouterService) GetRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.GetRouterRequest],
) (*connect.Response[mantraev1.GetRouterResponse], error) {
	var router *mantraev1.Router

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		res, err := s.app.Conn.GetQuery().GetHttpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoHttpRouter(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		res, err := s.app.Conn.GetQuery().GetTcpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		router, err = buildProtoTcpRouter(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		res, err := s.app.Conn.GetQuery().GetUdpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		router, err = buildProtoUdpRouter(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	return connect.NewResponse(&mantraev1.GetRouterResponse{Router: router}), nil
}

func (s *RouterService) CreateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateRouterRequest],
) (*connect.Response[mantraev1.CreateRouterResponse], error) {
	var router *mantraev1.Router
	var err error

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		var params db.CreateHttpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[dynamic.Router](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().CreateHttpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoHttpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		var params db.CreateTcpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = UnmarshalStruct[dynamic.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().CreateTcpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoTcpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		var params db.CreateUdpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = UnmarshalStruct[dynamic.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().CreateUdpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoUdpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	return connect.NewResponse(&mantraev1.CreateRouterResponse{Router: router}), nil
}

func (s *RouterService) UpdateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateRouterRequest],
) (*connect.Response[mantraev1.UpdateRouterResponse], error) {
	var router *mantraev1.Router
	var err error

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		var params db.UpdateHttpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[dynamic.Router](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().UpdateHttpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoHttpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		var params db.UpdateTcpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[dynamic.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().UpdateTcpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoTcpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		var params db.UpdateUdpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[dynamic.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		dbRouter, err := s.app.Conn.GetQuery().UpdateUdpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		router, err = buildProtoUdpRouter(dbRouter)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	return connect.NewResponse(&mantraev1.UpdateRouterResponse{Router: router}), nil
}

func (s *RouterService) DeleteRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteRouterRequest],
) (*connect.Response[mantraev1.DeleteRouterResponse], error) {
	err := s.app.Conn.GetQuery().DeleteHttpRouter(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteRouterResponse{}), nil
}

func (s *RouterService) ListRouters(
	ctx context.Context,
	req *connect.Request[mantraev1.ListRoutersRequest],
) (*connect.Response[mantraev1.ListRoutersResponse], error) {
	var limit int64
	var offset int64
	if req.Msg.Limit == nil {
		limit = 100
	} else {
		limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		offset = 0
	} else {
		offset = *req.Msg.Offset
	}

	var routers []*mantraev1.Router
	var totalCount int64
	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		params := db.ListHttpRoutersParams{
			Limit:  limit,
			Offset: offset,
		}
		dbRouters, err := s.app.Conn.GetQuery().ListHttpRouters(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountHttpRouters(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		for _, dbRouter := range dbRouters {
			router, err := buildProtoHttpRouter(dbRouter)
			if err != nil {
				slog.Error("Failed to build proto router", "error", err)
				continue
			}
			routers = append(routers, router)
		}

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		params := db.ListTcpRoutersParams{
			Limit:  limit,
			Offset: offset,
		}
		dbRouters, err := s.app.Conn.GetQuery().ListTcpRouters(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountTcpRouters(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		for _, dbRouter := range dbRouters {
			router, err := buildProtoTcpRouter(dbRouter)
			if err != nil {
				slog.Error("Failed to build proto router", "error", err)
				continue
			}
			routers = append(routers, router)
		}

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		params := db.ListUdpRoutersParams{
			Limit:  limit,
			Offset: offset,
		}
		dbRouters, err := s.app.Conn.GetQuery().ListUdpRouters(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountUdpRouters(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		for _, dbRouter := range dbRouters {
			router, err := buildProtoUdpRouter(dbRouter)
			if err != nil {
				slog.Error("Failed to build proto router", "error", err)
				continue
			}
			routers = append(routers, router)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	return connect.NewResponse(&mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}), nil
}

func buildProtoHttpRouter(r db.HttpRouter) (*mantraev1.Router, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HTTP config: %w", err)
	}
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Enabled:   r.Enabled,
		Type:      mantraev1.RouterType_ROUTER_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoTcpRouter(r db.TcpRouter) (*mantraev1.Router, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP config: %w", err)
	}
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Enabled:   r.Enabled,
		Type:      mantraev1.RouterType_ROUTER_TYPE_TCP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoUdpRouter(r db.UdpRouter) (*mantraev1.Router, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UDP config: %w", err)
	}
	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Enabled:   r.Enabled,
		Type:      mantraev1.RouterType_ROUTER_TYPE_UDP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}
