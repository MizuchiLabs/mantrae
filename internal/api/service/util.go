package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/meta"
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
	req *connect.Request[mantraev1.GetVersionRequest],
) (*connect.Response[mantraev1.GetVersionResponse], error) {
	return connect.NewResponse(&mantraev1.GetVersionResponse{
		Version: meta.Version,
	}), nil
}

func (s *UtilService) GetPublicIP(
	ctx context.Context,
	req *connect.Request[mantraev1.GetPublicIPRequest],
) (*connect.Response[mantraev1.GetPublicIPResponse], error) {
	ips, err := util.GetPublicIPs()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetPublicIPResponse{
		Ipv4: ips.IPv4,
		Ipv6: ips.IPv6,
	}), nil
}

func (s *UtilService) EventStream(
	ctx context.Context,
	req *connect.Request[mantraev1.EventStreamRequest],
	stream *connect.ServerStream[mantraev1.EventStreamResponse],
) error {
	id, ch := s.app.Event.Subscribe(req.Msg.ProfileId)
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
