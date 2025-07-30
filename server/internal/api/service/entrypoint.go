package service

import (
	"context"
	"slices"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type EntryPointService struct {
	app *config.App
}

func NewEntryPointService(app *config.App) *EntryPointService {
	return &EntryPointService{app: app}
}

func (s *EntryPointService) GetEntryPoint(
	ctx context.Context,
	req *connect.Request[mantraev1.GetEntryPointRequest],
) (*connect.Response[mantraev1.GetEntryPointResponse], error) {
	result, err := s.app.Conn.GetQuery().GetEntryPoint(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetEntryPointResponse{
		EntryPoint: result.ToProto(),
	}), nil
}

func (s *EntryPointService) CreateEntryPoint(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateEntryPointRequest],
) (*connect.Response[mantraev1.CreateEntryPointResponse], error) {
	params := db.CreateEntryPointParams{
		ProfileID: req.Msg.ProfileId,
		Name:      req.Msg.Name,
		Address:   req.Msg.Address,
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.IsDefault {
		if err := s.app.Conn.GetQuery().UnsetDefaultEntryPoint(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	result, err := s.app.Conn.GetQuery().CreateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateEntryPointResponse{
		EntryPoint: result.ToProto(),
	}), nil
}

func (s *EntryPointService) UpdateEntryPoint(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateEntryPointRequest],
) (*connect.Response[mantraev1.UpdateEntryPointResponse], error) {
	params := db.UpdateEntryPointParams{
		ID:        req.Msg.Id,
		Name:      req.Msg.Name,
		Address:   req.Msg.Address,
		IsDefault: req.Msg.IsDefault,
	}
	if req.Msg.IsDefault {
		if err := s.app.Conn.GetQuery().UnsetDefaultEntryPoint(ctx); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Remove old EntryPoint name and replace with new one (Order is important!)
	if err := s.updateRouterEntrypoints(ctx, req.Msg.Id, req.Msg.Name); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := s.app.Conn.GetQuery().UpdateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateEntryPointResponse{
		EntryPoint: result.ToProto(),
	}), nil
}

func (s *EntryPointService) DeleteEntryPoint(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteEntryPointRequest],
) (*connect.Response[mantraev1.DeleteEntryPointResponse], error) {
	if err := s.updateRouterEntrypoints(ctx, req.Msg.Id, ""); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.app.Conn.GetQuery().DeleteEntryPointByID(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteEntryPointResponse{}), nil
}

func (s *EntryPointService) ListEntryPoints(
	ctx context.Context,
	req *connect.Request[mantraev1.ListEntryPointsRequest],
) (*connect.Response[mantraev1.ListEntryPointsResponse], error) {
	params := db.ListEntryPointsParams{
		ProfileID: req.Msg.ProfileId,
		Limit:     req.Msg.Limit,
		Offset:    req.Msg.Offset,
	}

	result, err := s.app.Conn.GetQuery().ListEntryPoints(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountEntryPoints(ctx, req.Msg.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	entryPoints := make([]*mantraev1.EntryPoint, 0, len(result))
	for _, e := range result {
		entryPoints = append(entryPoints, e.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListEntryPointsResponse{
		EntryPoints: entryPoints,
		TotalCount:  totalCount,
	}), nil
}

// Helper functions
func (s *EntryPointService) updateRouterEntrypoints(
	ctx context.Context,
	id int64,
	newEntrypoint string,
) error {
	entrypoint, err := s.app.Conn.GetQuery().GetEntryPoint(ctx, id)
	if err != nil {
		return err
	}
	httpRouters, err := s.app.Conn.GetQuery().
		GetHttpRoutersUsingEntryPoint(ctx, db.GetHttpRoutersUsingEntryPointParams{
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
		if _, err = s.app.Conn.GetQuery().UpdateHttpRouter(ctx, db.UpdateHttpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}
	tcpRouters, err := s.app.Conn.GetQuery().
		GetTcpRoutersUsingEntryPoint(ctx, db.GetTcpRoutersUsingEntryPointParams{
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
		if _, err = s.app.Conn.GetQuery().UpdateTcpRouter(ctx, db.UpdateTcpRouterParams{
			ID:      r.ID,
			Enabled: r.Enabled,
			Config:  r.Config,
			Name:    r.Name,
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}
	udpRouters, err := s.app.Conn.GetQuery().
		GetUdpRoutersUsingEntryPoint(ctx, db.GetUdpRoutersUsingEntryPointParams{
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
		if _, err = s.app.Conn.GetQuery().UpdateUdpRouter(ctx, db.UpdateUdpRouterParams{
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
