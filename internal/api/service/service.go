package service

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
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
	var err error

	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		var params db.CreateHttpServiceParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[schema.Service](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[schema.TCPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = UnmarshalStruct[schema.UDPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
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
	var err error

	switch req.Msg.Type {
	case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
		var params db.UpdateHttpServiceParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.Service](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

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
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.TCPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

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
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = UnmarshalStruct[schema.UDPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

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

	if req.Msg.Type == nil {
		httpParams := db.ListHttpServicesParams{Limit: limit, Offset: offset}
		tcpParams := db.ListTcpServicesParams{Limit: limit, Offset: offset}
		udpParams := db.ListUdpServicesParams{Limit: limit, Offset: offset}

		httpServices, totalHttp, err := listServices[db.HttpService, mantraev1.Service, db.ListHttpServicesParams](
			ctx,
			s.app.Conn.GetQuery().ListHttpServices,
			s.app.Conn.GetQuery().CountHttpServices,
			buildProtoHttpService,
			httpParams,
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		tcpServices, totalTcp, err := listServices[db.TcpService, mantraev1.Service, db.ListTcpServicesParams](
			ctx,
			s.app.Conn.GetQuery().ListTcpServices,
			s.app.Conn.GetQuery().CountTcpServices,
			buildProtoTcpService,
			tcpParams,
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		udpServices, totalUdp, err := listServices[db.UdpService, mantraev1.Service, db.ListUdpServicesParams](
			ctx,
			s.app.Conn.GetQuery().ListUdpServices,
			s.app.Conn.GetQuery().CountUdpServices,
			buildProtoUdpService,
			udpParams,
		)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		services = append(services, httpServices...)
		services = append(services, tcpServices...)
		services = append(services, udpServices...)
		totalCount = totalHttp + totalTcp + totalUdp
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
			params := db.ListHttpServicesParams{Limit: limit, Offset: offset}
			services, totalCount, err = listServices[db.HttpService, mantraev1.Service, db.ListHttpServicesParams](
				ctx,
				s.app.Conn.GetQuery().ListHttpServices,
				s.app.Conn.GetQuery().CountHttpServices,
				buildProtoHttpService,
				params,
			)

		case mantraev1.ServiceType_SERVICE_TYPE_TCP:
			params := db.ListTcpServicesParams{Limit: limit, Offset: offset}
			services, totalCount, err = listServices[db.TcpService, mantraev1.Service, db.ListTcpServicesParams](
				ctx,
				s.app.Conn.GetQuery().ListTcpServices,
				s.app.Conn.GetQuery().CountTcpServices,
				buildProtoTcpService,
				params,
			)

		case mantraev1.ServiceType_SERVICE_TYPE_UDP:
			params := db.ListUdpServicesParams{Limit: limit, Offset: offset}
			services, totalCount, err = listServices[db.UdpService, mantraev1.Service, db.ListUdpServicesParams](
				ctx,
				s.app.Conn.GetQuery().ListUdpServices,
				s.app.Conn.GetQuery().CountUdpServices,
				buildProtoUdpService,
				params,
			)

		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, nil)
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.ListServicesResponse{
		Services:   services,
		TotalCount: totalCount,
	}), nil
}

// Helpers
func listServices[
	DBType any,
	ProtoType any,
	ParamsType any,
](
	ctx context.Context,
	listFn func(context.Context, ParamsType) ([]DBType, error),
	countFn func(context.Context) (int64, error),
	buildFn func(DBType) (*mantraev1.Service, error),
	params ParamsType,
) ([]*mantraev1.Service, int64, error) {
	dbServices, err := listFn(ctx, params)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	totalCount, err := countFn(ctx)
	if err != nil {
		return nil, 0, connect.NewError(connect.CodeInternal, err)
	}

	var services []*mantraev1.Service
	for _, dbService := range dbServices {
		service, err := buildFn(dbService)
		if err != nil {
			slog.Error("Failed to build proto service", "error", err)
			continue
		}
		services = append(services, service)
	}

	return services, totalCount, nil
}

func buildProtoHttpService(r db.HttpService) (*mantraev1.Service, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HTTP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoTcpService(r db.TcpService) (*mantraev1.Service, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}

func buildProtoUdpService(r db.UdpService) (*mantraev1.Service, error) {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UDP config: %w", err)
	}
	return &mantraev1.Service{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		Name:      r.Name,
		Config:    config,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}, nil
}
