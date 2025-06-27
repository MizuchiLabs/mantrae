// Package service provides the gRPC service implementations.
package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
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
		result, err := s.app.Conn.GetQuery().GetHttpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.HTTPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		result, err := s.app.Conn.GetQuery().GetTcpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.TCPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		result, err := s.app.Conn.GetQuery().GetUdpService(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.UDPServiceToProto(&result)

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
		params.Config, err = convert.UnmarshalStruct[schema.Service](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().CreateHttpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.HTTPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		var params db.CreateTcpServiceParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = convert.UnmarshalStruct[schema.TCPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().CreateTcpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.TCPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		var params db.CreateUdpServiceParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}
		params.Config, err = convert.UnmarshalStruct[schema.UDPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().CreateUdpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.UDPServiceToProto(&result)

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
		params.Config, err = convert.UnmarshalStruct[schema.Service](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().UpdateHttpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.HTTPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_TCP:
		var params db.UpdateTcpServiceParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = convert.UnmarshalStruct[schema.TCPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().UpdateTcpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.TCPServiceToProto(&result)

	case mantraev1.ServiceType_SERVICE_TYPE_UDP:
		var params db.UpdateUdpServiceParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Config, err = convert.UnmarshalStruct[schema.UDPService](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		result, err := s.app.Conn.GetQuery().UpdateUdpService(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.UDPServiceToProto(&result)

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return connect.NewResponse(&mantraev1.UpdateServiceResponse{Service: service}), nil
}

func (s *Service) DeleteService(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteServiceRequest],
) (*connect.Response[mantraev1.DeleteServiceResponse], error) {
	if err := s.app.Conn.GetQuery().DeleteHttpService(ctx, req.Msg.Id); err != nil {
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
		if req.Msg.AgentId == nil {
			result, err := s.app.Conn.GetQuery().
				ListServicesByProfile(ctx, db.ListServicesByProfileParams{
					Limit:       limit,
					Offset:      offset,
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
					ProfileID_3: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			totalCount, err = s.app.Conn.GetQuery().
				CountServicesByProfile(ctx, db.CountServicesByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
					ProfileID_3: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			services = convert.ServicesByProfileToProto(result)
		} else {
			result, err := s.app.Conn.GetQuery().
				ListServicesByAgent(ctx, db.ListServicesByAgentParams{
					Limit:     limit,
					Offset:    offset,
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
					AgentID_3: req.Msg.AgentId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			totalCount, err = s.app.Conn.GetQuery().
				CountServicesByAgent(ctx, db.CountServicesByAgentParams{
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
					AgentID_3: req.Msg.AgentId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			services = convert.ServicesByAgentToProto(result)
		}
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.ServiceType_SERVICE_TYPE_HTTP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountHttpServicesByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpServices(ctx, db.ListHttpServicesParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.HTTPServicesToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountHttpServicesByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpServicesByAgent(ctx, db.ListHttpServicesByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.HTTPServicesToProto(result)
			}

		case mantraev1.ServiceType_SERVICE_TYPE_TCP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountTcpServicesByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpServices(ctx, db.ListTcpServicesParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.TCPServicesToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountTcpServicesByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpServicesByAgent(ctx, db.ListTcpServicesByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.TCPServicesToProto(result)
			}

		case mantraev1.ServiceType_SERVICE_TYPE_UDP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountUdpServicesByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListUdpServices(ctx, db.ListUdpServicesParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.UDPServicesToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountUdpServicesByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListUdpServicesByAgent(ctx, db.ListUdpServicesByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				services = convert.UDPServicesToProto(result)
			}

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

func (s *Service) GetServiceByRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.GetServiceByRouterRequest],
) (*connect.Response[mantraev1.GetServiceByRouterResponse], error) {
	var service *mantraev1.Service

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		result, err := s.app.Conn.GetQuery().GetHttpServiceByName(ctx, req.Msg.Name)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.HTTPServiceToProto(&result)
	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		result, err := s.app.Conn.GetQuery().GetTcpServiceByName(ctx, req.Msg.Name)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.TCPServiceToProto(&result)
	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		result, err := s.app.Conn.GetQuery().GetUdpServiceByName(ctx, req.Msg.Name)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		service = convert.UDPServiceToProto(&result)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	return connect.NewResponse(&mantraev1.GetServiceByRouterResponse{Service: service}), nil
}
