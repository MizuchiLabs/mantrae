package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
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
	if req.Msg.ProfileId == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("profile id is required"),
		)
	}
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name is required"))
	}

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

		params.Config, err = UnmarshalStruct[schema.Router](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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

		params.Config, err = UnmarshalStruct[schema.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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

		params.Config, err = UnmarshalStruct[schema.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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
	if req.Msg.Id == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name is required"))
	}

	var router *mantraev1.Router
	var err error

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		var params db.UpdateHttpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.Router](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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
		params.Config, err = UnmarshalStruct[schema.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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
		params.Config, err = UnmarshalStruct[schema.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

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
	if req.Msg.ProfileId == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("profile id is required"),
		)
	}

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

	if req.Msg.Type == nil {
		httpRouters, totalHttp, err := listRouters[db.HttpRouter, mantraev1.Router, db.ListHttpRoutersParams](
			ctx,
			s.app.Conn.GetQuery().ListHttpRouters,
			s.app.Conn.GetQuery().CountHttpRouters,
			buildProtoHttpRouter,
			db.ListHttpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		tcpRouters, totalTcp, err := listRouters[db.TcpRouter, mantraev1.Router, db.ListTcpRoutersParams](
			ctx,
			s.app.Conn.GetQuery().ListTcpRouters,
			s.app.Conn.GetQuery().CountTcpRouters,
			buildProtoTcpRouter,
			db.ListTcpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		udpRouters, totalUdp, err := listRouters[db.UdpRouter, mantraev1.Router, db.ListUdpRoutersParams](
			ctx,
			s.app.Conn.GetQuery().ListUdpRouters,
			s.app.Conn.GetQuery().CountUdpRouters,
			buildProtoUdpRouter,
			db.ListUdpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		routers = append(routers, httpRouters...)
		routers = append(routers, tcpRouters...)
		routers = append(routers, udpRouters...)
		totalCount = totalHttp + totalTcp + totalUdp
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.RouterType_ROUTER_TYPE_HTTP:
			routers, totalCount, err = listRouters[db.HttpRouter, mantraev1.Router, db.ListHttpRoutersParams](
				ctx,
				s.app.Conn.GetQuery().ListHttpRouters,
				s.app.Conn.GetQuery().CountHttpRouters,
				buildProtoHttpRouter,
				db.ListHttpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
			)

		case mantraev1.RouterType_ROUTER_TYPE_TCP:
			routers, totalCount, err = listRouters[db.TcpRouter, mantraev1.Router, db.ListTcpRoutersParams](
				ctx,
				s.app.Conn.GetQuery().ListTcpRouters,
				s.app.Conn.GetQuery().CountTcpRouters,
				buildProtoTcpRouter,
				db.ListTcpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
			)

		case mantraev1.RouterType_ROUTER_TYPE_UDP:
			routers, totalCount, err = listRouters[db.UdpRouter, mantraev1.Router, db.ListUdpRoutersParams](
				ctx,
				s.app.Conn.GetQuery().ListUdpRouters,
				s.app.Conn.GetQuery().CountUdpRouters,
				buildProtoUdpRouter,
				db.ListUdpRoutersParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
			)

		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}), nil
}

// Helpers
func listRouters[
	DBType any,
	ProtoType any,
	ParamsType any,
](
	ctx context.Context,
	listFn func(context.Context, ParamsType) ([]DBType, error),
	countFn func(context.Context) (int64, error),
	buildFn func(DBType) (*mantraev1.Router, error),
	params ParamsType,
) ([]*mantraev1.Router, int64, error) {
	dbRouters, err := listFn(ctx, params)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	totalCount, err := countFn(ctx)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	var routers []*mantraev1.Router
	for _, dbRouter := range dbRouters {
		router, err := buildFn(dbRouter)
		if err != nil {
			slog.Error("Failed to build proto router", "error", err)
			continue
		}
		routers = append(routers, router)
	}

	return routers, totalCount, nil
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
