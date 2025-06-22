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
	"github.com/mizuchilabs/mantrae/internal/store/schema"
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
	var err error

	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		var params db.CreateHttpMiddlewareParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[schema.Middleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[schema.TCPMiddleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
	var err error

	switch req.Msg.Type {
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
		var params db.UpdateHttpMiddlewareParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.Middleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

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
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.TCPMiddleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

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

	if req.Msg.Type == nil {
		var err error
		if req.Msg.AgentId == nil {
			middlewares, totalCount, err = listMiddlewares[db.ListMiddlewaresByProfileRow, mantraev1.Middleware, db.ListMiddlewaresByProfileParams, db.CountMiddlewaresByProfileParams](
				ctx,
				s.app.Conn.GetQuery().ListMiddlewaresByProfile,
				s.app.Conn.GetQuery().CountMiddlewaresByProfile,
				buildMiddlewaresByProfile,
				db.ListMiddlewaresByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
					Limit:       limit,
					Offset:      offset,
				},
				db.CountMiddlewaresByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
				},
			)
		} else {
			middlewares, totalCount, err = listMiddlewares[db.ListMiddlewaresByAgentRow, mantraev1.Middleware, db.ListMiddlewaresByAgentParams, db.CountMiddlewaresByAgentParams](
				ctx,
				s.app.Conn.GetQuery().ListMiddlewaresByAgent,
				s.app.Conn.GetQuery().CountMiddlewaresByAgent,
				buildMiddlewaresByAgent,
				db.ListMiddlewaresByAgentParams{
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
					Limit:     limit,
					Offset:    offset,
				},
				db.CountMiddlewaresByAgentParams{
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
				},
			)
		}
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
			if req.Msg.AgentId == nil {
				middlewares, totalCount, err = listMiddlewares[db.HttpMiddleware, mantraev1.Middleware, db.ListHttpMiddlewaresParams, int64](
					ctx,
					s.app.Conn.GetQuery().ListHttpMiddlewares,
					s.app.Conn.GetQuery().CountHttpMiddlewaresByProfile,
					buildProtoHttpMiddleware,
					db.ListHttpMiddlewaresParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
					req.Msg.ProfileId,
				)
			} else {
				middlewares, totalCount, err = listMiddlewares[db.HttpMiddleware, mantraev1.Middleware, db.ListHttpMiddlewaresByAgentParams, *string](
					ctx,
					s.app.Conn.GetQuery().ListHttpMiddlewaresByAgent,
					s.app.Conn.GetQuery().CountHttpMiddlewaresByAgent,
					buildProtoHttpMiddleware,
					db.ListHttpMiddlewaresByAgentParams{AgentID: req.Msg.AgentId, Limit: limit, Offset: offset},
					req.Msg.AgentId,
				)
			}

		case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
			if req.Msg.AgentId == nil {
				middlewares, totalCount, err = listMiddlewares[db.TcpMiddleware, mantraev1.Middleware, db.ListTcpMiddlewaresParams, int64](
					ctx,
					s.app.Conn.GetQuery().ListTcpMiddlewares,
					s.app.Conn.GetQuery().CountTcpMiddlewaresByProfile,
					buildProtoTcpMiddleware,
					db.ListTcpMiddlewaresParams{ProfileID: req.Msg.ProfileId, Limit: limit, Offset: offset},
					req.Msg.ProfileId,
				)
			} else {
				middlewares, totalCount, err = listMiddlewares[db.TcpMiddleware, mantraev1.Middleware, db.ListTcpMiddlewaresByAgentParams, *string](
					ctx,
					s.app.Conn.GetQuery().ListTcpMiddlewaresByAgent,
					s.app.Conn.GetQuery().CountTcpMiddlewaresByAgent,
					buildProtoTcpMiddleware,
					db.ListTcpMiddlewaresByAgentParams{AgentID: req.Msg.AgentId, Limit: limit, Offset: offset},
					req.Msg.AgentId,
				)
			}

		default:
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid middleware type"),
			)
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.ListMiddlewaresResponse{
		Middlewares: middlewares,
		TotalCount:  totalCount,
	}), nil
}

// Helpers
func listMiddlewares[
	DBType any,
	ProtoType any,
	ListParams any,
	CountParams any,
](
	ctx context.Context,
	listFn func(context.Context, ListParams) ([]DBType, error),
	countFn func(context.Context, CountParams) (int64, error),
	buildFn func(DBType) (*mantraev1.Middleware, error),
	listParams ListParams,
	countParams CountParams,
) ([]*mantraev1.Middleware, int64, error) {
	dbMiddlewares, err := listFn(ctx, listParams)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	totalCount, err := countFn(ctx, countParams)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	var middlewares []*mantraev1.Middleware
	for _, dbMiddleware := range dbMiddlewares {
		middleware, err := buildFn(dbMiddleware)
		if err != nil {
			slog.Error("Failed to build proto middleware", "error", err)
			continue
		}
		middlewares = append(middlewares, middleware)
	}

	return middlewares, totalCount, nil
}

func buildMiddlewaresByProfile(r db.ListMiddlewaresByProfileRow) (*mantraev1.Middleware, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware config: %w", err)
	}
	var middleware mantraev1.Middleware
	switch r.Type {
	case "http":
		middleware = mantraev1.Middleware{
			Id:        r.ID,
			ProfileId: r.ProfileID,
			AgentId:   SafeString(r.AgentID),
			Name:      r.Name,
			Config:    config,
			Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
			CreatedAt: SafeTimestamp(r.CreatedAt),
			UpdatedAt: SafeTimestamp(r.UpdatedAt),
		}
	case "tcp":
		middleware = mantraev1.Middleware{
			Id:        r.ID,
			ProfileId: r.ProfileID,
			AgentId:   SafeString(r.AgentID),
			Name:      r.Name,
			Config:    config,
			Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
			CreatedAt: SafeTimestamp(r.CreatedAt),
			UpdatedAt: SafeTimestamp(r.UpdatedAt),
		}
	default:
		return nil, fmt.Errorf("invalid middleware type: %s", r.Type)
	}

	return &middleware, nil
}

func buildMiddlewaresByAgent(r db.ListMiddlewaresByAgentRow) (*mantraev1.Middleware, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal middleware config: %w", err)
	}
	var middleware mantraev1.Middleware
	switch r.Type {
	case "http":
		middleware = mantraev1.Middleware{
			Id:        r.ID,
			ProfileId: r.ProfileID,
			AgentId:   SafeString(r.AgentID),
			Name:      r.Name,
			Config:    config,
			Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
			CreatedAt: SafeTimestamp(r.CreatedAt),
			UpdatedAt: SafeTimestamp(r.UpdatedAt),
		}
	case "tcp":
		middleware = mantraev1.Middleware{
			Id:        r.ID,
			ProfileId: r.ProfileID,
			AgentId:   SafeString(r.AgentID),
			Name:      r.Name,
			Config:    config,
			Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
			CreatedAt: SafeTimestamp(r.CreatedAt),
			UpdatedAt: SafeTimestamp(r.UpdatedAt),
		}
	default:
		return nil, fmt.Errorf("invalid middleware type: %s", r.Type)
	}

	return &middleware, nil
}

func buildProtoHttpMiddleware(r db.HttpMiddleware) (*mantraev1.Middleware, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HTTP config: %w", err)
	}
	return &mantraev1.Middleware{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    config,
		Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoTcpMiddleware(r db.TcpMiddleware) (*mantraev1.Middleware, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP config: %w", err)
	}
	return &mantraev1.Middleware{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    config,
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

	var allPlugins []schema.Plugin
	if err := json.NewDecoder(resp.Body).Decode(&allPlugins); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	cleanPlugins := slices.DeleteFunc(allPlugins, func(p schema.Plugin) bool {
		return p.Type != "middleware"
	})
	var plugins []*mantraev1.Plugin
	for _, p := range cleanPlugins {
		plugins = append(plugins, &mantraev1.Plugin{
			Id:            p.ID,
			Name:          p.Name,
			DisplayName:   p.DisplayName,
			Author:        p.Author,
			Type:          p.Type,
			Import:        p.Import,
			Summary:       p.Summary,
			IconUrl:       p.IconUrl,
			BannerUrl:     p.BannerUrl,
			Readme:        p.Readme,
			LatestVersion: p.LatestVersion,
			Versions:      p.Versions,
			Stars:         p.Stars,
			Snippet: &mantraev1.PluginSnippet{
				K8S:  p.Snippet.K8S,
				Yaml: p.Snippet.Yaml,
				Toml: p.Snippet.Toml,
			},
			CreatedAt: p.CreatedAt,
		})
	}

	return connect.NewResponse(&mantraev1.GetMiddlewarePluginsResponse{Plugins: plugins}), nil
}
