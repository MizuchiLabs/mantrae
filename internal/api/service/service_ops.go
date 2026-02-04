package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type ServiceOps interface {
	Get(
		ctx context.Context,
		req *mantraev1.GetServiceRequest,
	) (*mantraev1.GetServiceResponse, error)
	Create(
		ctx context.Context,
		req *mantraev1.CreateServiceRequest,
	) (*mantraev1.CreateServiceResponse, error)
	Update(
		ctx context.Context,
		req *mantraev1.UpdateServiceRequest,
	) (*mantraev1.UpdateServiceResponse, error)
	Delete(
		ctx context.Context,
		req *mantraev1.DeleteServiceRequest,
	) (*mantraev1.DeleteServiceResponse, error)
	List(
		ctx context.Context,
		req *mantraev1.ListServicesRequest,
	) (*mantraev1.ListServicesResponse, error)
}

type HTTPServiceOps struct {
	app *config.App
}

type TCPServiceOps struct {
	app *config.App
}

type UDPServiceOps struct {
	app *config.App
}

func NewHTTPServiceOps(app *config.App) *HTTPServiceOps {
	return &HTTPServiceOps{app: app}
}

func NewTCPServiceOps(app *config.App) *TCPServiceOps {
	return &TCPServiceOps{app: app}
}

func NewUDPServiceOps(app *config.App) *UDPServiceOps {
	return &UDPServiceOps{app: app}
}

// HTTP Service Operations ----------------------------------------------------

func (s *HTTPServiceOps) Get(
	ctx context.Context,
	req *mantraev1.GetServiceRequest,
) (*mantraev1.GetServiceResponse, error) {
	var result *db.HttpService
	var err error

	switch id := req.GetIdentifier().(type) {
	case *mantraev1.GetServiceRequest_Id:
		result, err = s.app.Conn.Q.GetHttpService(ctx, id.Id)
		if err != nil {
			return nil, err
		}
	case *mantraev1.GetServiceRequest_Name:
		result, err = s.app.Conn.Q.GetHttpServiceByName(ctx, &db.GetHttpServiceByNameParams{
			ProfileID: req.ProfileId,
			Name:      id.Name,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid service identifier")
	}

	return &mantraev1.GetServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *HTTPServiceOps) Create(
	ctx context.Context,
	req *mantraev1.CreateServiceRequest,
) (*mantraev1.CreateServiceResponse, error) {
	params := &db.CreateHttpServiceParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.Service](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.CreateHttpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *HTTPServiceOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateServiceRequest,
) (*mantraev1.UpdateServiceResponse, error) {
	params := &db.UpdateHttpServiceParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.Service](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.UpdateHttpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *HTTPServiceOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteServiceRequest,
) (*mantraev1.DeleteServiceResponse, error) {
	if err := s.app.Conn.Q.DeleteHttpService(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteServiceResponse{}, nil
}

func (s *HTTPServiceOps) List(
	ctx context.Context,
	req *mantraev1.ListServicesRequest,
) (*mantraev1.ListServicesResponse, error) {
	result, err := s.app.Conn.Q.
		ListHttpServices(ctx, &db.ListHttpServicesParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.CountHttpServices(ctx, &db.CountHttpServicesParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	services := make([]*mantraev1.Service, 0, len(result))
	for _, s := range result {
		services = append(services, s.ToProto())
	}
	return &mantraev1.ListServicesResponse{
		Services:   services,
		TotalCount: totalCount,
	}, nil
}

// TCP Service Operations -----------------------------------------------------

func (s *TCPServiceOps) Get(
	ctx context.Context,
	req *mantraev1.GetServiceRequest,
) (*mantraev1.GetServiceResponse, error) {
	var result *db.TcpService
	var err error

	switch id := req.GetIdentifier().(type) {
	case *mantraev1.GetServiceRequest_Id:
		result, err = s.app.Conn.Q.GetTcpService(ctx, id.Id)
		if err != nil {
			return nil, err
		}
	case *mantraev1.GetServiceRequest_Name:
		result, err = s.app.Conn.Q.GetTcpServiceByName(ctx, &db.GetTcpServiceByNameParams{
			ProfileID: req.ProfileId,
			Name:      id.Name,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid service identifier")
	}

	return &mantraev1.GetServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *TCPServiceOps) Create(
	ctx context.Context,
	req *mantraev1.CreateServiceRequest,
) (*mantraev1.CreateServiceResponse, error) {
	params := &db.CreateTcpServiceParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.TCPService](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.CreateTcpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *TCPServiceOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateServiceRequest,
) (*mantraev1.UpdateServiceResponse, error) {
	params := &db.UpdateTcpServiceParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.TCPService](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.UpdateTcpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *TCPServiceOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteServiceRequest,
) (*mantraev1.DeleteServiceResponse, error) {
	if err := s.app.Conn.Q.DeleteTcpService(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteServiceResponse{}, nil
}

func (s *TCPServiceOps) List(
	ctx context.Context,
	req *mantraev1.ListServicesRequest,
) (*mantraev1.ListServicesResponse, error) {
	result, err := s.app.Conn.Q.
		ListTcpServices(ctx, &db.ListTcpServicesParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.CountTcpServices(ctx, &db.CountTcpServicesParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	services := make([]*mantraev1.Service, 0, len(result))
	for _, s := range result {
		services = append(services, s.ToProto())
	}
	return &mantraev1.ListServicesResponse{
		Services:   services,
		TotalCount: totalCount,
	}, nil
}

// UDP Service Operations -----------------------------------------------------

func (s *UDPServiceOps) Get(
	ctx context.Context,
	req *mantraev1.GetServiceRequest,
) (*mantraev1.GetServiceResponse, error) {
	var result *db.UdpService
	var err error

	switch id := req.GetIdentifier().(type) {
	case *mantraev1.GetServiceRequest_Id:
		result, err = s.app.Conn.Q.GetUdpService(ctx, id.Id)
		if err != nil {
			return nil, err
		}
	case *mantraev1.GetServiceRequest_Name:
		result, err = s.app.Conn.Q.GetUdpServiceByName(ctx, &db.GetUdpServiceByNameParams{
			ProfileID: req.ProfileId,
			Name:      id.Name,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid service identifier")
	}

	return &mantraev1.GetServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *UDPServiceOps) Create(
	ctx context.Context,
	req *mantraev1.CreateServiceRequest,
) (*mantraev1.CreateServiceResponse, error) {
	params := &db.CreateUdpServiceParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.UDPService](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.CreateUdpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.CreateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *UDPServiceOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateServiceRequest,
) (*mantraev1.UpdateServiceResponse, error) {
	params := &db.UpdateUdpServiceParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[dynamic.UDPService](req.Config)
	if err != nil {
		return nil, err
	}

	result, err := s.app.Conn.Q.UpdateUdpService(ctx, params)
	if err != nil {
		return nil, err
	}
	return &mantraev1.UpdateServiceResponse{
		Service: result.ToProto(),
	}, nil
}

func (s *UDPServiceOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteServiceRequest,
) (*mantraev1.DeleteServiceResponse, error) {
	if err := s.app.Conn.Q.DeleteUdpService(ctx, req.Id); err != nil {
		return nil, err
	}
	return &mantraev1.DeleteServiceResponse{}, nil
}

func (s *UDPServiceOps) List(
	ctx context.Context,
	req *mantraev1.ListServicesRequest,
) (*mantraev1.ListServicesResponse, error) {
	result, err := s.app.Conn.Q.
		ListUdpServices(ctx, &db.ListUdpServicesParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.CountUdpServices(ctx, &db.CountUdpServicesParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	services := make([]*mantraev1.Service, 0, len(result))
	for _, s := range result {
		services = append(services, s.ToProto())
	}
	return &mantraev1.ListServicesResponse{
		Services:   services,
		TotalCount: totalCount,
	}, nil
}
