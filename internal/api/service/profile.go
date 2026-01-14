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
	req *connect.Request[mantraev1.GetProfileRequest],
) (*connect.Response[mantraev1.GetProfileResponse], error) {
	result, err := s.app.Conn.Query.GetProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetProfileResponse{
		Profile: result.ToProto(),
	}), nil
}

func (s *ProfileService) CreateProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateProfileRequest],
) (*connect.Response[mantraev1.CreateProfileResponse], error) {
	params := &db.CreateProfileParams{
		Name:        slug.Make(req.Msg.Name),
		Description: req.Msg.Description,
		Token:       util.GenerateToken(6),
	}

	result, err := s.app.Conn.Query.CreateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Profile{
			Profile: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.CreateProfileResponse{
		Profile: result.ToProto(),
	}), nil
}

func (s *ProfileService) UpdateProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateProfileRequest],
) (*connect.Response[mantraev1.UpdateProfileResponse], error) {
	params := &db.UpdateProfileParams{
		ID:          req.Msg.Id,
		Name:        slug.Make(req.Msg.Name),
		Description: req.Msg.Description,
	}
	if req.Msg.GetRegenerateToken() {
		params.Token = util.GenerateToken(6)
	} else {
		profile, err := s.app.Conn.Query.GetProfile(ctx, params.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		params.Token = profile.Token
	}

	result, err := s.app.Conn.Query.UpdateProfile(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Profile{
			Profile: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.UpdateProfileResponse{
		Profile: result.ToProto(),
	}), nil
}

func (s *ProfileService) DeleteProfile(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteProfileRequest],
) (*connect.Response[mantraev1.DeleteProfileResponse], error) {
	profile, err := s.app.Conn.Query.GetProfile(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.app.Conn.Query.DeleteProfile(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Profile{
			Profile: profile.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.DeleteProfileResponse{}), nil
}

func (s *ProfileService) ListProfiles(
	ctx context.Context,
	req *connect.Request[mantraev1.ListProfilesRequest],
) (*connect.Response[mantraev1.ListProfilesResponse], error) {
	params := &db.ListProfilesParams{
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	}

	result, err := s.app.Conn.Query.ListProfiles(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	totalCount, err := s.app.Conn.Query.CountProfiles(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	profiles := make([]*mantraev1.Profile, 0, len(result))
	for _, p := range result {
		profiles = append(profiles, p.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListProfilesResponse{
		Profiles:   profiles,
		TotalCount: totalCount,
	}), nil
}
