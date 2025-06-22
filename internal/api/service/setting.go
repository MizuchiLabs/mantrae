package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type SettingService struct {
	app *config.App
}

func NewSettingService(app *config.App) *SettingService {
	return &SettingService{app: app}
}

func (s *SettingService) GetSetting(
	ctx context.Context,
	req *connect.Request[mantraev1.GetSettingRequest],
) (*connect.Response[mantraev1.GetSettingResponse], error) {
	value, ok := s.app.SM.Get(req.Msg.Key)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("setting not found"))
	}

	return connect.NewResponse(&mantraev1.GetSettingResponse{Value: value}), nil
}

func (s *SettingService) UpdateSetting(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateSettingRequest],
) (*connect.Response[mantraev1.UpdateSettingResponse], error) {
	params := db.UpsertSettingParams{
		Key:   req.Msg.Key,
		Value: req.Msg.Value,
	}
	if err := s.app.SM.Set(ctx, params.Key, params.Value); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateSettingResponse{
		Setting: &mantraev1.Setting{Key: params.Key, Value: params.Value},
	}), nil
}

func (s *SettingService) ListSettings(
	ctx context.Context,
	req *connect.Request[mantraev1.ListSettingsRequest],
) (*connect.Response[mantraev1.ListSettingsResponse], error) {
	var settings []*mantraev1.Setting
	for key, val := range s.app.SM.GetAll() {
		settings = append(settings, &mantraev1.Setting{Key: key, Value: val})
	}
	return connect.NewResponse(&mantraev1.ListSettingsResponse{Settings: settings}), nil
}
