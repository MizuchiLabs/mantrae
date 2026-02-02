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
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
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
	req *mantraev1.GetMiddlewareRequest,
) (*mantraev1.GetMiddlewareResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Get(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *MiddlewareService) CreateMiddleware(
	ctx context.Context,
	req *mantraev1.CreateMiddlewareRequest,
) (*mantraev1.CreateMiddlewareResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Create(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *MiddlewareService) UpdateMiddleware(
	ctx context.Context,
	req *mantraev1.UpdateMiddlewareRequest,
) (*mantraev1.UpdateMiddlewareResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Update(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *MiddlewareService) DeleteMiddleware(
	ctx context.Context,
	req *mantraev1.DeleteMiddlewareRequest,
) (*mantraev1.DeleteMiddlewareResponse, error) {
	ops, ok := s.dispatch[req.Type]
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid middleware type"),
		)
	}

	result, err := ops.Delete(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return result, nil
}

func (s *MiddlewareService) ListMiddlewares(
	ctx context.Context,
	req *mantraev1.ListMiddlewaresRequest,
) (*mantraev1.ListMiddlewaresResponse, error) {
	if req.Type != nil {
		ops, ok := s.dispatch[*req.Type]
		if !ok {
			return nil, connect.NewError(
				connect.CodeInvalidArgument,
				errors.New("invalid middleware type"),
			)
		}

		result, err := ops.List(ctx, req)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		return result, nil
	}

	// Get HTTP middlewares
	httpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP]
	httpResult, err := httpOps.List(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Get TCP middlewares
	tcpOps := s.dispatch[mantraev1.ProtocolType_PROTOCOL_TYPE_TCP]
	tcpResult, err := tcpOps.List(ctx, req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Combine results
	allMiddlewares := append(httpResult.Middlewares, tcpResult.Middlewares...)
	totalCount := httpResult.TotalCount + tcpResult.TotalCount

	return &mantraev1.ListMiddlewaresResponse{
		Middlewares: allMiddlewares,
		TotalCount:  totalCount,
	}, nil
}

func (s *MiddlewareService) GetMiddlewarePlugins(
	ctx context.Context,
	req *mantraev1.GetMiddlewarePluginsRequest,
) (*mantraev1.GetMiddlewarePluginsResponse, error) {
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

	return &mantraev1.GetMiddlewarePluginsResponse{Plugins: plugins}, nil
}
