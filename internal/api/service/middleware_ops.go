package service

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
)

type MiddlewareOps interface {
	Get(
		ctx context.Context,
		req *mantraev1.GetMiddlewareRequest,
	) (*mantraev1.GetMiddlewareResponse, error)
	Create(
		ctx context.Context,
		req *mantraev1.CreateMiddlewareRequest,
	) (*mantraev1.CreateMiddlewareResponse, error)
	Update(
		ctx context.Context,
		req *mantraev1.UpdateMiddlewareRequest,
	) (*mantraev1.UpdateMiddlewareResponse, error)
	Delete(
		ctx context.Context,
		req *mantraev1.DeleteMiddlewareRequest,
	) (*mantraev1.DeleteMiddlewareResponse, error)
	List(
		ctx context.Context,
		req *mantraev1.ListMiddlewaresRequest,
	) (*mantraev1.ListMiddlewaresResponse, error)
}

type HTTPMiddlewareOps struct {
	app *config.App
}

type TCPMiddlewareOps struct {
	app *config.App
}

func NewHTTPMiddlewareOps(app *config.App) *HTTPMiddlewareOps {
	return &HTTPMiddlewareOps{app: app}
}

func NewTCPMiddlewareOps(app *config.App) *TCPMiddlewareOps {
	return &TCPMiddlewareOps{app: app}
}

// HTTP Middleware Operations -------------------------------------------------

func (s *HTTPMiddlewareOps) Get(
	ctx context.Context,
	req *mantraev1.GetMiddlewareRequest,
) (*mantraev1.GetMiddlewareResponse, error) {
	result, err := s.app.Conn.Q.GetHttpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *HTTPMiddlewareOps) Create(
	ctx context.Context,
	req *mantraev1.CreateMiddlewareRequest,
) (*mantraev1.CreateMiddlewareResponse, error) {
	params := &db.CreateHttpMiddlewareParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
		Name:      req.Name,
		IsDefault: req.IsDefault,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPMiddleware](req.Config)
	if err != nil {
		return nil, err
	}
	if err = params.Config.Verify(); err != nil {
		return nil, err
	}

	if req.IsDefault {
		if err = s.app.Conn.Q.UnsetDefaultHttpMiddleware(ctx, req.ProfileId); err != nil {
			return nil, err
		}
	}

	result, err := s.app.Conn.Q.CreateHttpMiddleware(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *HTTPMiddlewareOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateMiddlewareRequest,
) (*mantraev1.UpdateMiddlewareResponse, error) {
	params := &db.UpdateHttpMiddlewareParams{
		ID:        req.Id,
		Name:      req.Name,
		Enabled:   req.Enabled,
		IsDefault: req.IsDefault,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPMiddleware](req.Config)
	if err != nil {
		return nil, err
	}
	if err = params.Config.Verify(); err != nil {
		return nil, err
	}

	if req.IsDefault {
		if err = s.app.Conn.Q.UnsetDefaultHttpMiddleware(ctx, req.ProfileId); err != nil {
			return nil, err
		}
	}

	// Get old middleware for router update
	middleware, err := s.app.Conn.Q.GetHttpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Make sure routers using this middleware use the new name
	routers, err := s.app.Conn.Q.
		GetHttpRoutersUsingMiddleware(ctx, &db.GetHttpRoutersUsingMiddlewareParams{
			ProfileID: middleware.ProfileID,
			ID:        middleware.ID,
		})
	if err != nil {
		return nil, err
	}
	for _, r := range routers {
		if idx := slices.Index(r.Config.Middlewares, middleware.Name); idx != -1 {
			r.Config.Middlewares = slices.Delete(r.Config.Middlewares, idx, idx+1)
		}
		r.Config.Middlewares = append(r.Config.Middlewares, req.Name)
		if _, err = s.app.Conn.Q.UpdateHttpRouter(ctx, &db.UpdateHttpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return nil, err
		}
	}

	result, err := s.app.Conn.Q.UpdateHttpMiddleware(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *HTTPMiddlewareOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteMiddlewareRequest,
) (*mantraev1.DeleteMiddlewareResponse, error) {
	middleware, err := s.app.Conn.Q.GetHttpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Make sure to delete the middleware from related routers
	routers, err := s.app.Conn.Q.
		GetHttpRoutersUsingMiddleware(ctx, &db.GetHttpRoutersUsingMiddlewareParams{
			ProfileID: middleware.ProfileID,
			ID:        middleware.ID,
		})
	if err != nil {
		return nil, err
	}
	for _, r := range routers {
		if idx := slices.Index(r.Config.Middlewares, middleware.Name); idx != -1 {
			r.Config.Middlewares = slices.Delete(r.Config.Middlewares, idx, idx+1)
		}
		if _, err := s.app.Conn.Q.UpdateHttpRouter(ctx, &db.UpdateHttpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return nil, err
		}
	}

	if err := s.app.Conn.Q.DeleteHttpMiddleware(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteMiddlewareResponse{}, nil
}

func (s *HTTPMiddlewareOps) List(
	ctx context.Context,
	req *mantraev1.ListMiddlewaresRequest,
) (*mantraev1.ListMiddlewaresResponse, error) {
	result, err := s.app.Conn.Q.
		ListHttpMiddlewares(ctx, &db.ListHttpMiddlewaresParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.
		CountHttpMiddlewares(ctx, &db.CountHttpMiddlewaresParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
		})
	if err != nil {
		return nil, err
	}

	middlewares := make([]*mantraev1.Middleware, 0, len(result))
	for _, m := range result {
		middlewares = append(middlewares, m.ToProto())
	}
	return &mantraev1.ListMiddlewaresResponse{
		Middlewares: middlewares,
		TotalCount:  totalCount,
	}, nil
}

// TCP Middleware Operations --------------------------------------------------

func (s *TCPMiddlewareOps) Get(
	ctx context.Context,
	req *mantraev1.GetMiddlewareRequest,
) (*mantraev1.GetMiddlewareResponse, error) {
	result, err := s.app.Conn.Q.GetTcpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *TCPMiddlewareOps) Create(
	ctx context.Context,
	req *mantraev1.CreateMiddlewareRequest,
) (*mantraev1.CreateMiddlewareResponse, error) {
	params := &db.CreateTcpMiddlewareParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
		Name:      req.Name,
		IsDefault: req.IsDefault,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPMiddleware](req.Config)
	if err != nil {
		return nil, err
	}

	if req.IsDefault {
		if err = s.app.Conn.Q.UnsetDefaultTcpMiddleware(ctx, req.ProfileId); err != nil {
			return nil, err
		}
	}

	result, err := s.app.Conn.Q.CreateTcpMiddleware(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *TCPMiddlewareOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateMiddlewareRequest,
) (*mantraev1.UpdateMiddlewareResponse, error) {
	params := &db.UpdateTcpMiddlewareParams{
		Name:      req.Name,
		Enabled:   req.Enabled,
		IsDefault: req.IsDefault,
		ID:        req.Id,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPMiddleware](req.Config)
	if err != nil {
		return nil, err
	}

	if req.IsDefault {
		if err = s.app.Conn.Q.UnsetDefaultTcpMiddleware(ctx, req.ProfileId); err != nil {
			return nil, err
		}
	}

	// Get old middleware for router update
	middleware, err := s.app.Conn.Q.GetTcpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Make sure routers using this middleware use the new name
	routers, err := s.app.Conn.Q.
		GetTcpRoutersUsingMiddleware(ctx, &db.GetTcpRoutersUsingMiddlewareParams{
			ProfileID: middleware.ProfileID,
			ID:        middleware.ID,
		})
	if err != nil {
		return nil, err
	}
	for _, r := range routers {
		if idx := slices.Index(r.Config.Middlewares, middleware.Name); idx != -1 {
			r.Config.Middlewares = slices.Delete(r.Config.Middlewares, idx, idx+1)
		}
		r.Config.Middlewares = append(r.Config.Middlewares, req.Name)
		if _, err = s.app.Conn.Q.UpdateTcpRouter(ctx, &db.UpdateTcpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return nil, err
		}
	}

	result, err := s.app.Conn.Q.UpdateTcpMiddleware(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateMiddlewareResponse{
		Middleware: result.ToProto(),
	}, nil
}

func (s *TCPMiddlewareOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteMiddlewareRequest,
) (*mantraev1.DeleteMiddlewareResponse, error) {
	middleware, err := s.app.Conn.Q.GetTcpMiddleware(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Make sure to delete the middleware from related routers
	routers, err := s.app.Conn.Q.
		GetTcpRoutersUsingMiddleware(ctx, &db.GetTcpRoutersUsingMiddlewareParams{
			ProfileID: middleware.ProfileID,
			ID:        middleware.ID,
		})
	if err != nil {
		return nil, err
	}
	for _, r := range routers {
		if idx := slices.Index(r.Config.Middlewares, middleware.Name); idx != -1 {
			r.Config.Middlewares = slices.Delete(r.Config.Middlewares, idx, idx+1)
		}
		if _, err := s.app.Conn.Q.UpdateTcpRouter(ctx, &db.UpdateTcpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return nil, err
		}
	}

	if err := s.app.Conn.Q.DeleteTcpMiddleware(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteMiddlewareResponse{}, nil
}

func (s *TCPMiddlewareOps) List(
	ctx context.Context,
	req *mantraev1.ListMiddlewaresRequest,
) (*mantraev1.ListMiddlewaresResponse, error) {
	result, err := s.app.Conn.Q.
		ListTcpMiddlewares(ctx, &db.ListTcpMiddlewaresParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.CountTcpMiddlewares(ctx, &db.CountTcpMiddlewaresParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	middlewares := make([]*mantraev1.Middleware, 0, len(result))
	for _, m := range result {
		middlewares = append(middlewares, m.ToProto())
	}
	return &mantraev1.ListMiddlewaresResponse{
		Middlewares: middlewares,
		TotalCount:  totalCount,
	}, nil
}
