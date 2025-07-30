package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
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
