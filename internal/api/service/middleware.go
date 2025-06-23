package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
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
		result, err := s.app.Conn.GetQuery().GetHttpMiddleware(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.HTTPMiddlewareToProto(&result)
	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		result, err := s.app.Conn.GetQuery().GetTcpMiddleware(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.TCPMiddlewareToProto(&result)

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
		params.Config, err = convert.UnmarshalStruct[schema.Middleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().CreateHttpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.HTTPMiddlewareToProto(&result)

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		var params db.CreateTcpMiddlewareParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = convert.UnmarshalStruct[schema.TCPMiddleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().CreateTcpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.TCPMiddlewareToProto(&result)

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
		params.Config, err = convert.UnmarshalStruct[schema.Middleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().UpdateHttpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.HTTPMiddlewareToProto(&result)

	case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
		var params db.UpdateTcpMiddlewareParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = convert.UnmarshalStruct[schema.TCPMiddleware](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().UpdateTcpMiddleware(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		middleware = convert.TCPMiddlewareToProto(&result)
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
	if err := s.app.Conn.GetQuery().DeleteHttpMiddleware(ctx, req.Msg.Id); err != nil {
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
			result, err := s.app.Conn.GetQuery().
				ListMiddlewaresByProfile(ctx, db.ListMiddlewaresByProfileParams{
					Limit:       limit,
					Offset:      offset,
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			totalCount, err = s.app.Conn.GetQuery().
				CountMiddlewaresByProfile(ctx, db.CountMiddlewaresByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			middlewares = convert.MiddlewaresByProfileToProto(result)
		} else {
			result, err := s.app.Conn.GetQuery().
				ListMiddlewaresByAgent(ctx, db.ListMiddlewaresByAgentParams{
					Limit:     limit,
					Offset:    offset,
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			totalCount, err = s.app.Conn.GetQuery().
				CountMiddlewaresByAgent(ctx, db.CountMiddlewaresByAgentParams{
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			middlewares = convert.MiddlewaresByAgentToProto(result)
		}
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountHttpMiddlewaresByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpMiddlewares(ctx, db.ListHttpMiddlewaresParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				middlewares = convert.HTTPMiddlewaresToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountHttpMiddlewaresByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpMiddlewaresByAgent(ctx, db.ListHttpMiddlewaresByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				middlewares = convert.HTTPMiddlewaresToProto(result)
			}

		case mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountTcpMiddlewaresByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpMiddlewares(ctx, db.ListTcpMiddlewaresParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				middlewares = convert.TCPMiddlewaresToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountTcpMiddlewaresByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpMiddlewaresByAgent(ctx, db.ListTcpMiddlewaresByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				middlewares = convert.TCPMiddlewaresToProto(result)
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
