package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type MiddlewareService struct {
	app *config.App
}

func NewMiddlewareService(app *config.App) *MiddlewareService {
	return &MiddlewareService{app: app}
}

func (s *MiddlewareService) GetMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.GetMiddlewareRequest],
) (*connect.Response[mantraev1.GetMiddlewareResponse], error) {
	var middleware *mantraev1.Middleware

	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		res, err := s.app.Conn.GetQuery().GetHttpMiddleware(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		middleware, err = buildProtoHttpMiddleware(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		res, err := s.app.Conn.GetQuery().GetTcpMiddleware(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware, err = buildProtoTcpMiddleware(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	return connect.NewResponse(&mantraev1.GetMiddlewareResponse{Middleware: middleware}), nil
}

func (s *MiddlewareService) CreateMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateMiddlewareRequest],
) (*connect.Response[mantraev1.CreateMiddlewareResponse], error) {
	var middleware *mantraev1.Middleware

	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		var params db.CreateHttpMiddlewareParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		dbMiddleware, err := s.app.Conn.GetQuery().CreateHttpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		middleware, err = buildProtoHttpMiddleware(dbMiddleware)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		var params db.CreateTcpMiddlewareParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		dbMiddleware, err := s.app.Conn.GetQuery().CreateTcpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		middleware, err = buildProtoTcpMiddleware(dbMiddleware)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	return connect.NewResponse(&mantraev1.CreateMiddlewareResponse{Middleware: middleware}), nil
}

func (s *MiddlewareService) UpdateMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateMiddlewareRequest],
) (*connect.Response[mantraev1.UpdateMiddlewareResponse], error) {
	var middleware *mantraev1.Middleware

	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		var params db.UpdateHttpMiddlewareParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name

		dbMiddleware, err := s.app.Conn.GetQuery().UpdateHttpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		middleware, err = buildProtoHttpMiddleware(dbMiddleware)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		var params db.UpdateTcpMiddlewareParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name

		dbMiddleware, err := s.app.Conn.GetQuery().UpdateTcpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		middleware, err = buildProtoTcpMiddleware(dbMiddleware)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	return connect.NewResponse(&mantraev1.UpdateMiddlewareResponse{Middleware: middleware}), nil
}

func (s *MiddlewareService) DeleteMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteMiddlewareRequest],
) (*connect.Response[mantraev1.DeleteMiddlewareResponse], error) {
	err := s.app.Conn.GetQuery().DeleteHttpMiddleware(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteMiddlewareResponse{}), nil
}

func (s *MiddlewareService) ListMiddlewares(
	ctx context.Context,
	req *connect.Request[mantraev1.ListMiddlewaresRequest],
) (*connect.Response[mantraev1.ListMiddlewaresResponse], error) {
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

	var middlewares []*mantraev1.Middleware
	var totalCount int64
	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		params := db.ListHttpMiddlewaresParams{Limit: limit, Offset: offset}
		dbMiddlewares, err := s.app.Conn.GetQuery().ListHttpMiddlewares(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountHttpMiddlewares(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		for _, dbMiddleware := range dbMiddlewares {
			middleware, err := buildProtoHttpMiddleware(dbMiddleware)
			if err != nil {
				slog.Error("Failed to build proto middleware", "error", err)
				continue
			}
			middlewares = append(middlewares, middleware)
		}

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		params := db.ListTcpMiddlewaresParams{Limit: limit, Offset: offset}
		dbMiddlewares, err := s.app.Conn.GetQuery().ListTcpMiddlewares(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountTcpMiddlewares(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		for _, dbMiddleware := range dbMiddlewares {
			middleware, err := buildProtoTcpMiddleware(dbMiddleware)
			if err != nil {
				slog.Error("Failed to build proto middleware", "error", err)
				continue
			}
			middlewares = append(middlewares, middleware)
		}

	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	return connect.NewResponse(&mantraev1.ListMiddlewaresResponse{
		Middlewares: middlewares,
		TotalCount:  totalCount,
	}), nil
}

func buildProtoHttpMiddleware(r db.HttpMiddleware) (*mantraev1.Middleware, error) {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HTTP config: %w", err)
	}
	return &mantraev1.Middleware{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    string(configBytes),
		Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoTcpMiddleware(r db.TcpMiddleware) (*mantraev1.Middleware, error) {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP config: %w", err)
	}
	return &mantraev1.Middleware{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    string(configBytes),
		Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func (s *MiddlewareService) GetMiddlewarePlugins(
	ctx context.Context,
	req *connect.Request[mantraev1.GetMiddlewarePluginsRequest],
) (*connect.Response[mantraev1.GetMiddlewarePluginsResponse], error) {
	resp, err := http.Get("https://plugins.traefik.io/api/services/plugins")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer resp.Body.Close()

	var allPlugins []*mantraev1.Plugin
	if err := json.NewDecoder(resp.Body).Decode(&allPlugins); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	plugins := slices.DeleteFunc(allPlugins, func(p *mantraev1.Plugin) bool {
		return p.Type != "middleware"
	})

	return connect.NewResponse(&mantraev1.GetMiddlewarePluginsResponse{Plugins: plugins}), nil
}
