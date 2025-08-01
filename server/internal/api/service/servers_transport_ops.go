package service

import (
	"context"

	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
)

type ServersTransportOps interface {
	Get(
		ctx context.Context,
		req *mantraev1.GetServersTransportRequest,
	) (*mantraev1.GetServersTransportResponse, error)
	Create(
		ctx context.Context,
		req *mantraev1.CreateServersTransportRequest,
	) (*mantraev1.CreateServersTransportResponse, error)
	Update(
		ctx context.Context,
		req *mantraev1.UpdateServersTransportRequest,
	) (*mantraev1.UpdateServersTransportResponse, error)
	Delete(
		ctx context.Context,
		req *mantraev1.DeleteServersTransportRequest,
	) (*mantraev1.DeleteServersTransportResponse, error)
	List(
		ctx context.Context,
		req *mantraev1.ListServersTransportsRequest,
	) (*mantraev1.ListServersTransportsResponse, error)
}

type HTTPServersTransportOps struct {
	app *config.App
}

type TCPServersTransportOps struct {
	app *config.App
}

func NewHTTPServersTransportOps(app *config.App) *HTTPServersTransportOps {
	return &HTTPServersTransportOps{app: app}
}

func NewTCPServersTransportOps(app *config.App) *TCPServersTransportOps {
	return &TCPServersTransportOps{app: app}
}

// HTTP Servers Transport Operations ------------------------------------------

func (s *HTTPServersTransportOps) Get(
	ctx context.Context,
	req *mantraev1.GetServersTransportRequest,
) (*mantraev1.GetServersTransportResponse, error) {
	result, err := s.app.Conn.GetQuery().GetHttpServersTransport(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *HTTPServersTransportOps) Create(
	ctx context.Context,
	req *mantraev1.CreateServersTransportRequest,
) (*mantraev1.CreateServersTransportResponse, error) {
	params := db.CreateHttpServersTransportParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
		Name:      req.Name,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPServersTransport](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.GetQuery().CreateHttpServersTransport(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *HTTPServersTransportOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateServersTransportRequest,
) (*mantraev1.UpdateServersTransportResponse, error) {
	params := db.UpdateHttpServersTransportParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPServersTransport](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.GetQuery().UpdateHttpServersTransport(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *HTTPServersTransportOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteServersTransportRequest,
) (*mantraev1.DeleteServersTransportResponse, error) {
	if err := s.app.Conn.GetQuery().DeleteHttpServersTransport(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteServersTransportResponse{}, nil
}

func (s *HTTPServersTransportOps) List(
	ctx context.Context,
	req *mantraev1.ListServersTransportsRequest,
) (*mantraev1.ListServersTransportsResponse, error) {
	result, err := s.app.Conn.GetQuery().
		ListHttpServersTransports(ctx, db.ListHttpServersTransportsParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.GetQuery().
		CountHttpServersTransports(ctx, db.CountHttpServersTransportsParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
		})
	if err != nil {
		return nil, err
	}

	serversTransports := make([]*mantraev1.ServersTransport, 0, len(result))
	for _, s := range result {
		serversTransports = append(serversTransports, s.ToProto())
	}
	return &mantraev1.ListServersTransportsResponse{
		ServersTransports: serversTransports,
		TotalCount:        totalCount,
	}, nil
}

// TCP Servers Transport Operations -------------------------------------------

func (s *TCPServersTransportOps) Get(
	ctx context.Context,
	req *mantraev1.GetServersTransportRequest,
) (*mantraev1.GetServersTransportResponse, error) {
	result, err := s.app.Conn.GetQuery().GetTcpServersTransport(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *TCPServersTransportOps) Create(
	ctx context.Context,
	req *mantraev1.CreateServersTransportRequest,
) (*mantraev1.CreateServersTransportResponse, error) {
	params := db.CreateTcpServersTransportParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
		Name:      req.Name,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPServersTransport](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.GetQuery().CreateTcpServersTransport(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *TCPServersTransportOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateServersTransportRequest,
) (*mantraev1.UpdateServersTransportResponse, error) {
	params := db.UpdateTcpServersTransportParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPServersTransport](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.GetQuery().UpdateTcpServersTransport(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateServersTransportResponse{
		ServersTransport: result.ToProto(),
	}, nil
}

func (s *TCPServersTransportOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteServersTransportRequest,
) (*mantraev1.DeleteServersTransportResponse, error) {
	if err := s.app.Conn.GetQuery().DeleteTcpServersTransport(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteServersTransportResponse{}, nil
}

func (s *TCPServersTransportOps) List(
	ctx context.Context,
	req *mantraev1.ListServersTransportsRequest,
) (*mantraev1.ListServersTransportsResponse, error) {
	result, err := s.app.Conn.GetQuery().
		ListTcpServersTransports(ctx, db.ListTcpServersTransportsParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.GetQuery().
		CountTcpServersTransports(ctx, db.CountTcpServersTransportsParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
		})
	if err != nil {
		return nil, err
	}

	serversTransports := make([]*mantraev1.ServersTransport, 0, len(result))
	for _, s := range result {
		serversTransports = append(serversTransports, s.ToProto())
	}
	return &mantraev1.ListServersTransportsResponse{
		ServersTransports: serversTransports,
		TotalCount:        totalCount,
	}, nil
}
