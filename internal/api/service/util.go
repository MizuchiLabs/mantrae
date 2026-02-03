package service

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type UtilService struct {
	app *config.App
}

func NewUtilService(app *config.App) *UtilService {
	return &UtilService{app: app}
}

func (s *UtilService) GetVersion(
	ctx context.Context,
	req *mantraev1.GetVersionRequest,
) (*mantraev1.GetVersionResponse, error) {
	return &mantraev1.GetVersionResponse{Version: s.app.Version}, nil
}

func (s *UtilService) GetDynamicConfig(
	ctx context.Context,
	req *mantraev1.GetDynamicConfigRequest,
) (*mantraev1.GetDynamicConfigResponse, error) {
	profile, err := s.app.Conn.Q.GetProfile(ctx, req.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cfg, err := traefik.BuildDynamicConfig(ctx, s.app.Conn.Q, *profile)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.GetDynamicConfigResponse{Config: string(jsonBytes)}, nil
}

func (s *UtilService) GetPublicIP(
	ctx context.Context,
	req *mantraev1.GetPublicIPRequest,
) (*mantraev1.GetPublicIPResponse, error) {
	ips, err := util.GetPublicIPs()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.GetPublicIPResponse{
		Ipv4: ips.IPv4,
		Ipv6: ips.IPv6,
	}, nil
}

func (s *UtilService) EventStream(
	ctx context.Context,
	req *mantraev1.EventStreamRequest,
	stream *connect.ServerStream[mantraev1.EventStreamResponse],
) error {
	id, ch := s.app.Event.Subscribe(req.ProfileId)
	defer s.app.Event.Unsubscribe(id)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}
