package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	profilev1 "github.com/mizuchilabs/mantrae/proto/gen/server/v1"
)

type ProfileService struct {
	app *config.App
}

func NewProfileService(app *config.App) *ProfileService {
	return &ProfileService{app: app}
}

func (s *ProfileService) GetProfile(
	ctx context.Context,
	req *connect.Request[profilev1.GetProfileRequest],
) (*connect.Response[profilev1.GetProfileResponse], error) {
	profile, err := s.app.Conn.GetQuery().GetProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&profilev1.GetProfileResponse{
		Profile: &profilev1.Profile{
			Id:          profile.ID,
			Name:        profile.Name,
			Description: SafeString(profile.Description),
			CreatedAt:   SafeTimestamp(profile.CreatedAt),
			UpdatedAt:   SafeTimestamp(profile.UpdatedAt),
		},
	}), nil
}

func (s *ProfileService) CreateProfile(
	ctx context.Context,
	req *connect.Request[profilev1.CreateProfileRequest],
) (*connect.Response[profilev1.CreateProfileResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	params := db.CreateProfileParams{
		ID:   id.String(),
		Name: req.Msg.Name,
	}
	if req.Msg.Description != nil {
		params.Description = req.Msg.Description
	}

	profile, err := s.app.Conn.GetQuery().CreateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&profilev1.CreateProfileResponse{
		Profile: &profilev1.Profile{
			Id:          profile.ID,
			Name:        profile.Name,
			Description: SafeString(profile.Description),
			CreatedAt:   SafeTimestamp(profile.CreatedAt),
			UpdatedAt:   SafeTimestamp(profile.UpdatedAt),
		},
	}), nil
}

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	req *connect.Request[profilev1.UpdateProfileRequest],
) (*connect.Response[profilev1.UpdateProfileResponse], error) {
	params := db.UpdateProfileParams{
		ID:   req.Msg.Id,
		Name: req.Msg.Name,
	}
	if req.Msg.Description != nil {
		params.Description = req.Msg.Description
	}

	profile, err := s.app.Conn.GetQuery().UpdateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&profilev1.UpdateProfileResponse{
		Profile: &profilev1.Profile{
			Id:          profile.ID,
			Name:        profile.Name,
			Description: SafeString(profile.Description),
			CreatedAt:   SafeTimestamp(profile.CreatedAt),
			UpdatedAt:   SafeTimestamp(profile.UpdatedAt),
		},
	}), nil
}

func (s *ProfileService) DeleteProfile(
	ctx context.Context,
	req *connect.Request[profilev1.DeleteProfileRequest],
) (*connect.Response[profilev1.DeleteProfileResponse], error) {
	err := s.app.Conn.GetQuery().DeleteProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&profilev1.DeleteProfileResponse{}), nil
}

func (s *ProfileService) ListProfiles(
	ctx context.Context,
	req *connect.Request[profilev1.ListProfilesRequest],
) (*connect.Response[profilev1.ListProfilesResponse], error) {
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

	dbProfiles, err := s.app.Conn.GetQuery().ListProfiles(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountProfiles(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var profiles []*profilev1.Profile
	for _, profile := range dbProfiles {
		profiles = append(profiles, &profilev1.Profile{
			Id:          profile.ID,
			Name:        profile.Name,
			Description: SafeString(profile.Description),
			CreatedAt:   SafeTimestamp(profile.CreatedAt),
			UpdatedAt:   SafeTimestamp(profile.UpdatedAt),
		})
	}
	return connect.NewResponse(&profilev1.ListProfilesResponse{
		Profiles:   profiles,
		TotalCount: totalCount,
	}), nil
}
