package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

type SettingService struct {
	app *config.App
}

func NewSettingService(app *config.App) *SettingService {
	return &SettingService{app: app}
}

func (s *SettingService) GetSetting(
	ctx context.Context,
	req *mantraev1.GetSettingRequest,
) (*mantraev1.GetSettingResponse, error) {
	value, ok := s.app.SM.Get(ctx, req.Key)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("setting not found"))
	}
	return &mantraev1.GetSettingResponse{Value: value}, nil
}

func (s *SettingService) UpdateSetting(
	ctx context.Context,
	req *mantraev1.UpdateSettingRequest,
) (*mantraev1.UpdateSettingResponse, error) {
	params := &db.UpsertSettingParams{
		Key:   req.Key,
		Value: req.Value,
	}
	if err := s.app.SM.Set(ctx, params.Key, params.Value); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.UpdateSettingResponse{
		Setting: &mantraev1.Setting{Key: params.Key, Value: params.Value},
	}, nil
}

func (s *SettingService) ListSettings(
	ctx context.Context,
	req *mantraev1.ListSettingsRequest,
) (*mantraev1.ListSettingsResponse, error) {
	var settings []*mantraev1.Setting
	for key, val := range s.app.SM.GetAll(ctx) {
		settings = append(settings, &mantraev1.Setting{Key: key, Value: val})
	}
	return &mantraev1.ListSettingsResponse{Settings: settings}, nil
}
