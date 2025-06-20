package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
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
	entryPoint, err := s.app.Conn.GetQuery().GetEntryPoint(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetEntryPointResponse{
		EntryPoint: &mantraev1.EntryPoint{
			Id:        entryPoint.ID,
			Name:      entryPoint.Name,
			Address:   entryPoint.Address,
			IsDefault: entryPoint.IsDefault,
			CreatedAt: SafeTimestamp(entryPoint.CreatedAt),
			UpdatedAt: SafeTimestamp(entryPoint.UpdatedAt),
		},
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
	entryPoint, err := s.app.Conn.GetQuery().CreateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateEntryPointResponse{
		EntryPoint: &mantraev1.EntryPoint{
			Id:        entryPoint.ID,
			ProfileId: entryPoint.ProfileID,
			Name:      entryPoint.Name,
			Address:   entryPoint.Address,
			IsDefault: entryPoint.IsDefault,
			CreatedAt: SafeTimestamp(entryPoint.CreatedAt),
			UpdatedAt: SafeTimestamp(entryPoint.UpdatedAt),
		},
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
	entryPoint, err := s.app.Conn.GetQuery().UpdateEntryPoint(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateEntryPointResponse{
		EntryPoint: &mantraev1.EntryPoint{
			Id:        entryPoint.ID,
			ProfileId: entryPoint.ProfileID,
			Name:      entryPoint.Name,
			Address:   entryPoint.Address,
			IsDefault: entryPoint.IsDefault,
			CreatedAt: SafeTimestamp(entryPoint.CreatedAt),
			UpdatedAt: SafeTimestamp(entryPoint.UpdatedAt),
		},
	}), nil
}

func (s *EntryPointService) DeleteEntryPoint(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteEntryPointRequest],
) (*connect.Response[mantraev1.DeleteEntryPointResponse], error) {
	err := s.app.Conn.GetQuery().DeleteEntryPoint(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteEntryPointResponse{}), nil
}

func (s *EntryPointService) ListEntryPoints(
	ctx context.Context,
	req *connect.Request[mantraev1.ListEntryPointsRequest],
) (*connect.Response[mantraev1.ListEntryPointsResponse], error) {
	if req.Msg.ProfileId == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("profile id is required"),
		)
	}

	var params db.ListEntryPointsParams
	params.ProfileID = req.Msg.ProfileId
	if req.Msg.Limit == nil {
		params.Limit = 100
	} else {
		params.Limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *req.Msg.Offset
	}

	dbEntryPoints, err := s.app.Conn.GetQuery().ListEntryPoints(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountEntryPoints(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var entryPoints []*mantraev1.EntryPoint
	for _, entryPoint := range dbEntryPoints {
		entryPoints = append(entryPoints, &mantraev1.EntryPoint{
			Id:        entryPoint.ID,
			ProfileId: entryPoint.ProfileID,
			Name:      entryPoint.Name,
			Address:   entryPoint.Address,
			IsDefault: entryPoint.IsDefault,
			CreatedAt: SafeTimestamp(entryPoint.CreatedAt),
			UpdatedAt: SafeTimestamp(entryPoint.UpdatedAt),
		})
	}
	return connect.NewResponse(&mantraev1.ListEntryPointsResponse{
		EntryPoints: entryPoints,
		TotalCount:  totalCount,
	}), nil
}
