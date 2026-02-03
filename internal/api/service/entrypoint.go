package service

import (
	"context"
	"slices"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

type EntryPointService struct {
	app *config.App
}

func NewEntryPointService(app *config.App) *EntryPointService {
	return &EntryPointService{app: app}
}

func (s *EntryPointService) GetEntryPoint(
	ctx context.Context,
	req *mantraev1.GetEntryPointRequest,
) (*mantraev1.GetEntryPointResponse, error) {
	result, err := s.app.Conn.Q.GetEntryPoint(ctx, req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.GetEntryPointResponse{EntryPoint: result.ToProto()}, nil
}

func (s *EntryPointService) CreateEntryPoint(
	ctx context.Context,
	req *mantraev1.CreateEntryPointRequest,
) (*mantraev1.CreateEntryPointResponse, error) {
	params := &db.CreateEntryPointParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	}
	if req.IsDefault {
		if err := s.app.Conn.Q.UnsetDefaultEntryPoint(ctx, req.ProfileId); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.Q.CreateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.CreateEntryPointResponse{EntryPoint: result.ToProto()}, nil
}

func (s *EntryPointService) UpdateEntryPoint(
	ctx context.Context,
	req *mantraev1.UpdateEntryPointRequest,
) (*mantraev1.UpdateEntryPointResponse, error) {
	params := &db.UpdateEntryPointParams{
		ID:        req.Id,
		Name:      req.Name,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	}
	if req.IsDefault {
		if err := s.app.Conn.Q.UnsetDefaultEntryPoint(ctx, req.ProfileId); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Remove old EntryPoint name and replace with new one (Order is important!)
	if err := s.updateRouterEntrypoints(ctx, req.Id, req.Name); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := s.app.Conn.Q.UpdateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.UpdateEntryPointResponse{EntryPoint: result.ToProto()}, nil
}

func (s *EntryPointService) DeleteEntryPoint(
	ctx context.Context,
	req *mantraev1.DeleteEntryPointRequest,
) (*mantraev1.DeleteEntryPointResponse, error) {
	if err := s.updateRouterEntrypoints(ctx, req.Id, ""); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.app.Conn.Q.DeleteEntryPointByID(ctx, req.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.DeleteEntryPointResponse{}, nil
}

func (s *EntryPointService) ListEntryPoints(
	ctx context.Context,
	req *mantraev1.ListEntryPointsRequest,
) (*mantraev1.ListEntryPointsResponse, error) {
	params := &db.ListEntryPointsParams{
		ProfileID: req.ProfileId,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	result, err := s.app.Conn.Q.ListEntryPoints(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.Q.CountEntryPoints(ctx, req.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	entryPoints := make([]*mantraev1.EntryPoint, 0, len(result))
	for _, e := range result {
		entryPoints = append(entryPoints, e.ToProto())
	}
	return &mantraev1.ListEntryPointsResponse{
		EntryPoints: entryPoints,
		TotalCount:  totalCount,
	}, nil
}

// Helper functions
func (s *EntryPointService) updateRouterEntrypoints(
	ctx context.Context,
	id,
	newEntrypoint string,
) error {
	entrypoint, err := s.app.Conn.Q.GetEntryPoint(ctx, id)
	if err != nil {
		return err
	}
	httpRouters, err := s.app.Conn.Q.
		GetHttpRoutersUsingEntryPoint(ctx, &db.GetHttpRoutersUsingEntryPointParams{
			ProfileID: entrypoint.ProfileID,
			ID:        entrypoint.ID,
		})
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	for _, r := range httpRouters {
		if idx := slices.Index(r.Config.EntryPoints, entrypoint.Name); idx != -1 {
			r.Config.EntryPoints = slices.Delete(r.Config.EntryPoints, idx, idx+1)
		}
		if newEntrypoint != "" {
			r.Config.EntryPoints = append(r.Config.EntryPoints, newEntrypoint)
		}
		if _, err = s.app.Conn.Q.UpdateHttpRouter(ctx, &db.UpdateHttpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}
	tcpRouters, err := s.app.Conn.Q.
		GetTcpRoutersUsingEntryPoint(ctx, &db.GetTcpRoutersUsingEntryPointParams{
			ProfileID: entrypoint.ProfileID,
			ID:        entrypoint.ID,
		})
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	for _, r := range tcpRouters {
		if idx := slices.Index(r.Config.EntryPoints, entrypoint.Name); idx != -1 {
			r.Config.EntryPoints = slices.Delete(r.Config.EntryPoints, idx, idx+1)
		}
		if newEntrypoint != "" {
			r.Config.EntryPoints = append(r.Config.EntryPoints, newEntrypoint)
		}
		if _, err = s.app.Conn.Q.UpdateTcpRouter(ctx, &db.UpdateTcpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}
	udpRouters, err := s.app.Conn.Q.
		GetUdpRoutersUsingEntryPoint(ctx, &db.GetUdpRoutersUsingEntryPointParams{
			ProfileID: entrypoint.ProfileID,
			ID:        entrypoint.ID,
		})
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	for _, r := range udpRouters {
		if idx := slices.Index(r.Config.EntryPoints, entrypoint.Name); idx != -1 {
			r.Config.EntryPoints = slices.Delete(r.Config.EntryPoints, idx, idx+1)
		}
		if newEntrypoint != "" {
			r.Config.EntryPoints = append(r.Config.EntryPoints, newEntrypoint)
		}
		if _, err = s.app.Conn.Q.UpdateUdpRouter(ctx, &db.UpdateUdpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}
