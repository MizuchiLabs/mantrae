package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type Service struct {
	app *config.App
}

func NewServiceService(app *config.App) *Service {
	return &Service{app: app}
}

func (s *Service) GetService(
	ctx context.Context,
	req *connect.Request[mantraev1.GetServiceRequest],
) (*connect.Response[mantraev1.GetServiceResponse], error) {
	var service *mantraev1.Service

	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		res, err := s.app.Conn.GetQuery().GetHttpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoHttpService(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		res, err := s.app.Conn.GetQuery().GetTcpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service, err = buildProtoTcpService(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		res, err := s.app.Conn.GetQuery().GetUdpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service, err = buildProtoUdpService(res)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	return connect.NewResponse(&mantraev1.GetServiceResponse{Service: service}), nil
}

func (s *Service) CreateService(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateServiceRequest],
) (*connect.Response[mantraev1.CreateServiceResponse], error) {
	var service *mantraev1.Service

	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		var params db.CreateHttpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		dbService, err := s.app.Conn.GetQuery().CreateHttpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoHttpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		var params db.CreateTcpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		dbService, err := s.app.Conn.GetQuery().CreateTcpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoTcpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		var params db.CreateUdpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		dbService, err := s.app.Conn.GetQuery().CreateUdpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoUdpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return connect.NewResponse(&mantraev1.CreateServiceResponse{Service: service}), nil
}

func (s *Service) UpdateService(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateServiceRequest],
) (*connect.Response[mantraev1.UpdateServiceResponse], error) {
	var service *mantraev1.Service

	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		var params db.UpdateHttpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name

		dbService, err := s.app.Conn.GetQuery().UpdateHttpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoHttpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		var params db.UpdateTcpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name

		dbService, err := s.app.Conn.GetQuery().UpdateTcpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoTcpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		var params db.UpdateUdpServiceParams
		if err := json.Unmarshal([]byte(req.Msg.Config), &params.Config); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name

		dbService, err := s.app.Conn.GetQuery().UpdateUdpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		service, err = buildProtoUdpService(dbService)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return connect.NewResponse(&mantraev1.UpdateServiceResponse{Service: service}), nil
}

func (s *Service) DeleteService(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteServiceRequest],
) (*connect.Response[mantraev1.DeleteServiceResponse], error) {
	err := s.app.Conn.GetQuery().DeleteHttpService(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteServiceResponse{}), nil
}

func (s *Service) ListServices(
	ctx context.Context,
	req *connect.Request[mantraev1.ListServicesRequest],
) (*connect.Response[mantraev1.ListServicesResponse], error) {
	var limit int64
	var offset int64
	if req.Msg.Limit == nil {
		limit = 100
	} else {
		limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		offset = 0
	} else {
		offset = *req.Msg.Offset
	}

	var services []*mantraev1.Service
	var totalCount int64
	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		params := db.ListHttpServicesParams{Limit: limit, Offset: offset}
		dbServices, err := s.app.Conn.GetQuery().ListHttpServices(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountHttpServices(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		for _, dbService := range dbServices {
			service, err := buildProtoHttpService(dbService)
			if err != nil {
				slog.Error("Failed to build proto service", "error", err)
				continue
			}
			services = append(services, service)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		params := db.ListTcpServicesParams{Limit: limit, Offset: offset}
		dbServices, err := s.app.Conn.GetQuery().ListTcpServices(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountTcpServices(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		for _, dbService := range dbServices {
			service, err := buildProtoTcpService(dbService)
			if err != nil {
				slog.Error("Failed to build proto service", "error", err)
				continue
			}
			services = append(services, service)
		}

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		params := db.ListUdpServicesParams{Limit: limit, Offset: offset}
		dbServices, err := s.app.Conn.GetQuery().ListUdpServices(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		totalCount, err = s.app.Conn.GetQuery().CountUdpServices(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		for _, dbService := range dbServices {
			service, err := buildProtoUdpService(dbService)
			if err != nil {
				slog.Error("Failed to build proto service", "error", err)
				continue
			}
			services = append(services, service)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	return connect.NewResponse(&mantraev1.ListServicesResponse{
		Services:   services,
		TotalCount: totalCount,
	}), nil
}

func buildProtoHttpService(r db.HttpService) (*mantraev1.Service, error) {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HTTP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    string(configBytes),
		Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoTcpService(r db.TcpService) (*mantraev1.Service, error) {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    string(configBytes),
		Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoUdpService(r db.UdpService) (*mantraev1.Service, error) {
	configBytes, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UDP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    string(configBytes),
		Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}
