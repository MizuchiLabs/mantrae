package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type MiddlewareService struct {
	app      *config.App
	dispatch map[mantraev1.ProtocolType]MiddlewareOps
}

func NewMiddlewareService(app *config.App) *MiddlewareService {
	return &MiddlewareService{
		app: app,
		dispatch: map[mantraev1.ProtocolType]MiddlewareOps{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP: NewHTTPMiddlewareOps(app),
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP:  NewTCPMiddlewareOps(app),
		},
	}
}

func (s *MiddlewareService) GetMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.GetMiddlewareRequest],
) (*connect.Response[mantraev1.GetMiddlewareResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Get(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *MiddlewareService) CreateMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateMiddlewareRequest],
) (*connect.Response[mantraev1.CreateMiddlewareResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Create(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *MiddlewareService) UpdateMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateMiddlewareRequest],
) (*connect.Response[mantraev1.UpdateMiddlewareResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Update(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *MiddlewareService) DeleteMiddleware(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteMiddlewareRequest],
) (*connect.Response[mantraev1.DeleteMiddlewareResponse], error) {
	ops, ok := s.dispatch[req.Msg.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Delete(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(result), nil
}

func (s *MiddlewareService) ListMiddlewares(
	ctx context.Context,
	req *connect.Request[mantraev1.ListMiddlewaresRequest],
) (*connect.Response[mantraev1.ListMiddlewaresResponse], error) {
	if req.Msg.Type != nil {
		ops, ok := s.dispatch[*req.Msg.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid middleware type"),
			)
		}

		result, err := ops.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return connect.NewResponse(result), nil
	} else {
		// Get HTTP middlewares
		httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
		httpResult, err := httpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Get TCP middlewares
		tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
		tcpResult, err := tcpOps.List(ctx, req.Msg)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		// Combine results
		allMiddlewares := append(httpResult.Middlewares, tcpResult.Middlewares...)
		totalCount := httpResult.TotalCount + tcpResult.TotalCount

		return connect.NewResponse(&mantraev1.ListMiddlewaresResponse{
			Middlewares: allMiddlewares,
			TotalCount:  totalCount,
		}), nil
	}
}

func (s *MiddlewareService) GetMiddlewarePlugins(
	ctx context.Context,
	req *connect.Request[mantraev1.GetMiddlewarePluginsRequest],
) (*connect.Response[mantraev1.GetMiddlewarePluginsResponse], error) {
	resp, err := http.Get("https://plugins.traefik.io/api/services/plugins")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	var allPlugins []schema.Plugin
	if err := json.NewDecoder(resp.Body).Decode(&allPlugins); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	filtered := slices.DeleteFunc(allPlugins, func(p schema.Plugin) bool {
		return p.Type != "middleware"
	})
	var plugins []*mantraev1.Plugin
	for _, p := range filtered {
		plugins = append(plugins, &mantraev1.Plugin{
			Id:            p.ID,
			Name:          p.Name,
			DisplayName:   p.DisplayName,
			Author:        p.Author,
			Type:          p.Type,
			Import:        p.Import,
			Summary:       p.Summary,
			IconUrl:       p.IconURL,
			BannerUrl:     p.BannerURL,
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
