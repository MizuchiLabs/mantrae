package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type ProfileService struct {
	app *config.App
}

func NewProfileService(app *config.App) *ProfileService {
	return &ProfileService{app: app}
}

func (s *ProfileService) GetProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.GetProfileRequest],
) (*connect.Response[mantraev1.GetProfileResponse], error) {
	result, err := s.app.Conn.GetQuery().GetProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetProfileResponse{
		Profile: convert.ProfileToProto(&result),
	}), nil
}

func (s *ProfileService) CreateProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateProfileRequest],
) (*connect.Response[mantraev1.CreateProfileResponse], error) {
	params := db.CreateProfileParams{
		Name: req.Msg.Name,
	}
	if req.Msg.Description != nil {
		params.Description = req.Msg.Description
	}

	result, err := s.app.Conn.GetQuery().CreateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateProfileResponse{
		Profile: convert.ProfileToProto(&result),
	}), nil
}

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateProfileRequest],
) (*connect.Response[mantraev1.UpdateProfileResponse], error) {
	params := db.UpdateProfileParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
	}
	if req.Msg.Description != nil {
		params.Description = req.Msg.Description
	}

	result, err := s.app.Conn.GetQuery().UpdateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateProfileResponse{
		Profile: convert.ProfileToProto(&result),
	}), nil
}

func (s *ProfileService) DeleteProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteProfileRequest],
) (*connect.Response[mantraev1.DeleteProfileResponse], error) {
	err := s.app.Conn.GetQuery().DeleteProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteProfileResponse{}), nil
}

func (s *ProfileService) ListProfiles(
	ctx context.Context,
	req *connect.Request[mantraev1.ListProfilesRequest],
) (*connect.Response[mantraev1.ListProfilesResponse], error) {
	var params db.ListProfilesParams
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

	result, err := s.app.Conn.GetQuery().ListProfiles(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountProfiles(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.ListProfilesResponse{
		Profiles:   convert.ProfilesToProto(result),
		TotalCount: totalCount,
	}), nil
}
