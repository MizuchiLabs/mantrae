package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/gosimple/slug"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type ProfileService struct {
	app *config.App
}

func NewProfileService(app *config.App) *ProfileService {
	return &ProfileService{app: app}
}

func (s *ProfileService) GetProfile(
	ctx context.Context,
	req *mantraev1.GetProfileRequest,
) (*mantraev1.GetProfileResponse, error) {
	result, err := s.app.Conn.Q.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.GetProfileResponse{Profile: result.ToProto()}, nil
}

func (s *ProfileService) CreateProfile(
	ctx context.Context,
	req *mantraev1.CreateProfileRequest,
) (*mantraev1.CreateProfileResponse, error) {
	params := &db.CreateProfileParams{
		Name:        slug.Make(req.Name),
		Description: req.Description,
		Token:       util.GenerateToken(6),
	}

	result, err := s.app.Conn.Q.CreateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.CreateProfileResponse{Profile: result.ToProto()}, nil
}

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	req *mantraev1.UpdateProfileRequest,
) (*mantraev1.UpdateProfileResponse, error) {
	params := &db.UpdateProfileParams{
		ID:          req.Id,
		Name:        slug.Make(req.Name),
		Description: req.Description,
	}
	if req.GetRegenerateToken() {
		params.Token = util.GenerateToken(6)
	} else {
		profile, err := s.app.Conn.Q.GetProfile(ctx, params.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Token = profile.Token
	}

	result, err := s.app.Conn.Q.UpdateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.UpdateProfileResponse{Profile: result.ToProto()}, nil
}

func (s *ProfileService) DeleteProfile(
	ctx context.Context,
	req *mantraev1.DeleteProfileRequest,
) (*mantraev1.DeleteProfileResponse, error) {
	if err := s.app.Conn.Q.DeleteProfile(ctx, req.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.DeleteProfileResponse{}, nil
}

func (s *ProfileService) ListProfiles(
	ctx context.Context,
	req *mantraev1.ListProfilesRequest,
) (*mantraev1.ListProfilesResponse, error) {
	params := &db.ListProfilesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	result, err := s.app.Conn.Q.ListProfiles(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	totalCount, err := s.app.Conn.Q.CountProfiles(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	profiles := make([]*mantraev1.Profile, 0, len(result))
	for _, p := range result {
		profiles = append(profiles, p.ToProto())
	}
	return &mantraev1.ListProfilesResponse{
		Profiles:   profiles,
		TotalCount: totalCount,
	}, nil
}
