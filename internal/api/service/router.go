package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type RouterService struct {
	app *config.App
}

func NewRouterService(app *config.App) *RouterService {
	return &RouterService{app: app}
}

func (s *RouterService) GetRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.GetRouterRequest],
) (*connect.Response[mantraev1.GetRouterResponse], error) {
	var router *mantraev1.Router

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		resRouter, err := s.app.Conn.GetQuery().GetHttpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resDNS)
		router = convert.HTTPRouterToProto(&resRouter, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		resRouter, err := s.app.Conn.GetQuery().GetTcpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resDNS)
		router = convert.TCPRouterToProto(&resRouter, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		result, err := s.app.Conn.GetQuery().GetUdpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		router = convert.UDPRouterToProto(&result)

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	return connect.NewResponse(&mantraev1.GetRouterResponse{Router: router}), nil
}

func (s *RouterService) CreateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateRouterRequest],
) (*connect.Response[mantraev1.CreateRouterResponse], error) {
	var router *mantraev1.Router
	var err error

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		var params db.CreateHttpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = convert.UnmarshalStruct[schema.HTTPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		result, err := s.app.Conn.GetQuery().CreateHttpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resultDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, result.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resultDNS)
		router = convert.HTTPRouterToProto(&result, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		var params db.CreateTcpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = convert.UnmarshalStruct[schema.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		result, err := s.app.Conn.GetQuery().CreateTcpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resultDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, result.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resultDNS)
		router = convert.TCPRouterToProto(&result, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		var params db.CreateUdpRouterParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = convert.UnmarshalStruct[schema.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		result, err := s.app.Conn.GetQuery().CreateUdpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		router = convert.UDPRouterToProto(&result)

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	// Broadcast event
	// s.app.Event.BroadcastProfileEvent(&mantraev1.ProfileEvent{
	// 	EventType:    mantraev1.EventType_EVENT_TYPE_CREATED,
	// 	ResourceType: mantraev1.ResourceType_RESOURCE_TYPE_ROUTER,
	// 	ProfileId:    req.Msg.ProfileId,
	// 	Timestamp:    timestamppb.Now(),
	// 	Resource: &mantraev1.ProfileEvent_Router{
	// 		Router: router,
	// 	},
	// })

	return connect.NewResponse(&mantraev1.CreateRouterResponse{Router: router}), nil
}

func (s *RouterService) UpdateRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateRouterRequest],
) (*connect.Response[mantraev1.UpdateRouterResponse], error) {
	var router *mantraev1.Router
	var err error

	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		var params db.UpdateHttpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Enabled = req.Msg.Enabled
		params.Config, err = convert.UnmarshalStruct[schema.HTTPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		// Update DNS Providers
		existing, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, params.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		existingMap := make(map[int64]bool)
		for _, provider := range existing {
			existingMap[provider.ID] = true
		}

		desiredMap := make(map[int64]bool)
		var desiredIDs []int64
		for _, protoProvider := range req.Msg.DnsProviders {
			desiredMap[protoProvider.Id] = true
			desiredIDs = append(desiredIDs, protoProvider.Id)
		}
		// Identify inserts
		for _, id := range desiredIDs {
			if !existingMap[id] {
				if err = s.app.Conn.GetQuery().
					CreateHttpRouterDNSProvider(ctx, db.CreateHttpRouterDNSProviderParams{
						HttpRouterID:  params.ID,
						DnsProviderID: id,
					}); err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
			}
		}

		// Identify deletes
		for id := range existingMap {
			if !desiredMap[id] {
				if err = s.app.Conn.GetQuery().
					DeleteHttpRouterDNSProvider(ctx, db.DeleteHttpRouterDNSProviderParams{
						HttpRouterID:  params.ID,
						DnsProviderID: id,
					}); err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
			}
		}

		result, err := s.app.Conn.GetQuery().UpdateHttpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resultDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, result.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resultDNS)
		router = convert.HTTPRouterToProto(&result, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		var params db.UpdateTcpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Enabled = req.Msg.Enabled
		params.Config, err = convert.UnmarshalStruct[schema.TCPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		// Update DNS Providers
		existing, err := s.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, params.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		existingMap := make(map[int64]bool)
		for _, provider := range existing {
			existingMap[provider.ID] = true
		}

		desiredMap := make(map[int64]bool)
		var desiredIDs []int64
		for _, protoProvider := range req.Msg.DnsProviders {
			desiredMap[protoProvider.Id] = true
			desiredIDs = append(desiredIDs, protoProvider.Id)
		}

		// Identify inserts
		for _, id := range desiredIDs {
			if !existingMap[id] {
				if err = s.app.Conn.GetQuery().
					CreateTcpRouterDNSProvider(ctx, db.CreateTcpRouterDNSProviderParams{
						TcpRouterID:   params.ID,
						DnsProviderID: id,
					}); err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
			}
		}

		// Identify deletes
		for id := range existingMap {
			if !desiredMap[id] {
				if err = s.app.Conn.GetQuery().
					DeleteTcpRouterDNSProvider(ctx, db.DeleteTcpRouterDNSProviderParams{
						TcpRouterID:   params.ID,
						DnsProviderID: id,
					}); err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
			}
		}

		result, err := s.app.Conn.GetQuery().UpdateTcpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		resultDNS, err := s.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, result.ID)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		protoDNS := convert.DNSProvidersToProto(resultDNS)
		router = convert.TCPRouterToProto(&result, protoDNS)

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		var params db.UpdateUdpRouterParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Enabled = req.Msg.Enabled
		params.Config, err = convert.UnmarshalStruct[schema.UDPRouter](req.Msg.Config)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		params.Config.Service = params.Name

		result, err := s.app.Conn.GetQuery().UpdateUdpRouter(ctx, params)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		router = convert.UDPRouterToProto(&result)

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	// Broadcast event
	// s.app.Event.BroadcastProfileEvent(&mantraev1.ProfileEvent{
	// 	EventType:    mantraev1.EventType_EVENT_TYPE_UPDATED,
	// 	ResourceType: mantraev1.ResourceType_RESOURCE_TYPE_ROUTER,
	// 	ProfileId:    router.ProfileId,
	// 	Timestamp:    timestamppb.Now(),
	// 	Resource: &mantraev1.ProfileEvent_Router{
	// 		Router: router,
	// 	},
	// })

	return connect.NewResponse(&mantraev1.UpdateRouterResponse{Router: router}), nil
}

func (s *RouterService) DeleteRouter(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteRouterRequest],
) (*connect.Response[mantraev1.DeleteRouterResponse], error) {
	switch req.Msg.Type {
	case mantraev1.RouterType_ROUTER_TYPE_HTTP:
		result, err := s.app.Conn.GetQuery().GetHttpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if err := s.app.Conn.GetQuery().DeleteHttpRouter(ctx, req.Msg.Id); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if result.Config.Service != "" {
			service, err := s.app.Conn.GetQuery().GetHttpServiceByName(ctx, result.Config.Service)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			if err := s.app.Conn.GetQuery().DeleteHttpService(ctx, service.ID); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
		}

	case mantraev1.RouterType_ROUTER_TYPE_TCP:
		result, err := s.app.Conn.GetQuery().GetTcpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if err := s.app.Conn.GetQuery().DeleteTcpRouter(ctx, req.Msg.Id); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if result.Config.Service != "" {
			service, err := s.app.Conn.GetQuery().GetTcpServiceByName(ctx, result.Config.Service)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			if err := s.app.Conn.GetQuery().DeleteTcpService(ctx, service.ID); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
		}

	case mantraev1.RouterType_ROUTER_TYPE_UDP:
		result, err := s.app.Conn.GetQuery().GetUdpRouter(ctx, req.Msg.Id)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if err := s.app.Conn.GetQuery().DeleteUdpRouter(ctx, req.Msg.Id); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if result.Config.Service != "" {
			service, err := s.app.Conn.GetQuery().GetUdpServiceByName(ctx, result.Config.Service)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			if err := s.app.Conn.GetQuery().DeleteUdpService(ctx, service.ID); err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
	}

	// s.app.Event.BroadcastProfileEvent(&mantraev1.ProfileEvent{
	// 	EventType:    mantraev1.EventType_EVENT_TYPE_DELETED,
	// 	ResourceType: mantraev1.ResourceType_RESOURCE_TYPE_ROUTER,
	// 	ProfileId:    router.ProfileId,
	// 	Timestamp:    timestamppb.Now(),
	// 	Resource: &mantraev1.ProfileEvent_Router{
	// 		Router: router,
	// 	},
	// })

	return connect.NewResponse(&mantraev1.DeleteRouterResponse{}), nil
}

func (s *RouterService) ListRouters(
	ctx context.Context,
	req *connect.Request[mantraev1.ListRoutersRequest],
) (*connect.Response[mantraev1.ListRoutersResponse], error) {
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

	var routers []*mantraev1.Router
	var totalCount int64

	if req.Msg.Type == nil {
		if req.Msg.AgentId == nil {
			result, err := s.app.Conn.GetQuery().
				ListRoutersByProfile(ctx, db.ListRoutersByProfileParams{
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
				CountRoutersByProfile(ctx, db.CountRoutersByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
					ProfileID_3: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			routers = convert.RoutersByProfileToProto(result, s.app.Conn.GetQuery())
		} else {
			result, err := s.app.Conn.GetQuery().
				ListRoutersByAgent(ctx, db.ListRoutersByAgentParams{
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
				CountRoutersByAgent(ctx, db.CountRoutersByAgentParams{
					AgentID:   req.Msg.AgentId,
					AgentID_2: req.Msg.AgentId,
					AgentID_3: req.Msg.AgentId,
				})
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			routers = convert.RoutersByAgentToProto(result, s.app.Conn.GetQuery())
		}
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.RouterType_ROUTER_TYPE_HTTP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountHttpRoutersByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpRouters(ctx, db.ListHttpRoutersParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.HTTPRoutersToProto(result, s.app.Conn.GetQuery())
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountHttpRoutersByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListHttpRoutersByAgent(ctx, db.ListHttpRoutersByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.HTTPRoutersToProto(result, s.app.Conn.GetQuery())
			}

		case mantraev1.RouterType_ROUTER_TYPE_TCP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountTcpRoutersByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpRouters(ctx, db.ListTcpRoutersParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.TCPRoutersToProto(result, s.app.Conn.GetQuery())
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountTcpRoutersByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpRoutersByAgent(ctx, db.ListTcpRoutersByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.TCPRoutersToProto(result, s.app.Conn.GetQuery())
			}

		case mantraev1.RouterType_ROUTER_TYPE_UDP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountUdpRoutersByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListUdpRouters(ctx, db.ListUdpRoutersParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.UDPRoutersToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountUdpRoutersByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListUdpRoutersByAgent(ctx, db.ListUdpRoutersByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				routers = convert.UDPRoutersToProto(result)
			}

		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid router type"))
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}), nil
}
