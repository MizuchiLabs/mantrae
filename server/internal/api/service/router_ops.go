package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
)

type RouterOps interface {
	Get(
		ctx context.Context,
		req *mantraev1.GetRouterRequest,
	) (*mantraev1.GetRouterResponse, error)
	Create(
		ctx context.Context,
		req *mantraev1.CreateRouterRequest,
	) (*mantraev1.CreateRouterResponse, error)
	Update(
		ctx context.Context,
		req *mantraev1.UpdateRouterRequest,
	) (*mantraev1.UpdateRouterResponse, error)
	Delete(
		ctx context.Context,
		req *mantraev1.DeleteRouterRequest,
	) (*mantraev1.DeleteRouterResponse, error)
	List(
		ctx context.Context,
		req *mantraev1.ListRoutersRequest,
	) (*mantraev1.ListRoutersResponse, error)
}

type HTTPRouterOps struct {
	app *config.App
}

type TCPRouterOps struct {
	app *config.App
}

type UDPRouterOps struct {
	app *config.App
}

func NewHTTPRouterOps(app *config.App) *HTTPRouterOps {
	return &HTTPRouterOps{app: app}
}

func NewTCPRouterOps(app *config.App) *TCPRouterOps {
	return &TCPRouterOps{app: app}
}

func NewUDPRouterOps(app *config.App) *UDPRouterOps {
	return &UDPRouterOps{app: app}
}

// HTTP Router Operations -----------------------------------------------------

func (s *HTTPRouterOps) Get(
	ctx context.Context,
	req *mantraev1.GetRouterRequest,
) (*mantraev1.GetRouterResponse, error) {
	result, err := s.app.Conn.GetQuery().GetHttpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *HTTPRouterOps) Create(
	ctx context.Context,
	req *mantraev1.CreateRouterRequest,
) (*mantraev1.CreateRouterResponse, error) {
	params := db.CreateHttpRouterParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	result, err := s.app.Conn.GetQuery().CreateHttpRouter(ctx, params)
	if err != nil {
		return nil, err
	}
	router := result.ToProto()

	dnsProviders, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, result.ID)
	if err != nil {
		return nil, err
	}
	router.DnsProviders = make([]*mantraev1.DnsProvider, 0, len(dnsProviders))
	for _, p := range dnsProviders {
		router.DnsProviders = append(router.DnsProviders, p.ToProto())
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: router,
		},
	})
	return &mantraev1.CreateRouterResponse{
		Router: router,
	}, nil
}

func (s *HTTPRouterOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateRouterRequest,
) (*mantraev1.UpdateRouterResponse, error) {
	params := db.UpdateHttpRouterParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.HTTPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	// Update DNS Providers
	existing, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	existingMap := make(map[string]bool)
	for _, provider := range existing {
		existingMap[provider.ID] = true
	}

	desiredMap := make(map[string]bool)
	var desiredIDs []string
	for _, protoProvider := range req.DnsProviders {
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
				return nil, err
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
				return nil, err
			}
		}
	}

	result, err := s.app.Conn.GetQuery().UpdateHttpRouter(ctx, params)
	if err != nil {
		return nil, err
	}
	router := result.ToProto()

	dnsProviders, err := s.app.Conn.GetQuery().GetDnsProvidersByHttpRouter(ctx, result.ID)
	if err != nil {
		return nil, err
	}
	router.DnsProviders = make([]*mantraev1.DnsProvider, 0, len(dnsProviders))
	for _, p := range dnsProviders {
		router.DnsProviders = append(router.DnsProviders, p.ToProto())
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: router,
		},
	})
	return &mantraev1.UpdateRouterResponse{
		Router: router,
	}, nil
}

func (s *HTTPRouterOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteRouterRequest,
) (*mantraev1.DeleteRouterResponse, error) {
	router, err := s.app.Conn.GetQuery().GetHttpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if router.Config.Service != "" {
		service, err := s.app.Conn.GetQuery().
			GetHttpServiceByName(ctx, db.GetHttpServiceByNameParams{
				ProfileID: router.ProfileID,
				Name:      router.Config.Service,
			})
		if err != nil {
			slog.Error("failed to get http service", "err", err)
		}
		if err := s.app.Conn.GetQuery().DeleteHttpService(ctx, service.ID); err != nil {
			slog.Error("failed to delete http service", "err", err)
		}
	}

	if err := s.app.Conn.GetQuery().DeleteHttpRouter(ctx, req.Id); err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: router.ToProto(),
		},
	})
	return &mantraev1.DeleteRouterResponse{}, nil
}

func (s *HTTPRouterOps) List(
	ctx context.Context,
	req *mantraev1.ListRoutersRequest,
) (*mantraev1.ListRoutersResponse, error) {
	result, err := s.app.Conn.GetQuery().
		ListHttpRouters(ctx, db.ListHttpRoutersParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.GetQuery().CountHttpRouters(ctx, db.CountHttpRoutersParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	routers := make([]*mantraev1.Router, 0, len(result))
	for _, r := range result {
		routers = append(routers, r.ToProto())
	}
	return &mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}, nil
}

// TCP Router Operations ------------------------------------------------------

func (s *TCPRouterOps) Get(
	ctx context.Context,
	req *mantraev1.GetRouterRequest,
) (*mantraev1.GetRouterResponse, error) {
	result, err := s.app.Conn.GetQuery().GetTcpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *TCPRouterOps) Create(
	ctx context.Context,
	req *mantraev1.CreateRouterRequest,
) (*mantraev1.CreateRouterResponse, error) {
	params := db.CreateTcpRouterParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	result, err := s.app.Conn.GetQuery().CreateTcpRouter(ctx, params)
	if err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: result.ToProto(),
		},
	})
	return &mantraev1.CreateRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *TCPRouterOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateRouterRequest,
) (*mantraev1.UpdateRouterResponse, error) {
	params := db.UpdateTcpRouterParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.TCPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	// Update DNS Providers
	existing, err := s.app.Conn.GetQuery().GetDnsProvidersByTcpRouter(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	existingMap := make(map[string]bool)
	for _, provider := range existing {
		existingMap[provider.ID] = true
	}

	desiredMap := make(map[string]bool)
	var desiredIDs []string
	for _, protoProvider := range req.DnsProviders {
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
				return nil, err
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
				return nil, err
			}
		}
	}

	result, err := s.app.Conn.GetQuery().UpdateTcpRouter(ctx, params)
	if err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: result.ToProto(),
		},
	})
	return &mantraev1.UpdateRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *TCPRouterOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteRouterRequest,
) (*mantraev1.DeleteRouterResponse, error) {
	router, err := s.app.Conn.GetQuery().GetTcpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if router.Config.Service != "" {
		service, err := s.app.Conn.GetQuery().GetTcpServiceByName(ctx, db.GetTcpServiceByNameParams{
			ProfileID: router.ProfileID,
			Name:      router.Config.Service,
		})
		if err != nil {
			slog.Error("failed to get tcp service", "err", err)
		}
		if err := s.app.Conn.GetQuery().DeleteTcpService(ctx, service.ID); err != nil {
			slog.Error("failed to delete tcp service", "err", err)
		}
	}

	if err := s.app.Conn.GetQuery().DeleteTcpRouter(ctx, req.Id); err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: router.ToProto(),
		},
	})
	return &mantraev1.DeleteRouterResponse{}, nil
}

func (s *TCPRouterOps) List(
	ctx context.Context,
	req *mantraev1.ListRoutersRequest,
) (*mantraev1.ListRoutersResponse, error) {
	result, err := s.app.Conn.GetQuery().
		ListTcpRouters(ctx, db.ListTcpRoutersParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.GetQuery().CountTcpRouters(ctx, db.CountTcpRoutersParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	routers := make([]*mantraev1.Router, 0, len(result))
	for _, r := range result {
		routers = append(routers, r.ToProto())
	}
	return &mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}, nil
}

// UDP Router Operations ------------------------------------------------------

func (s *UDPRouterOps) Get(
	ctx context.Context,
	req *mantraev1.GetRouterRequest,
) (*mantraev1.GetRouterResponse, error) {
	result, err := s.app.Conn.GetQuery().GetUdpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *UDPRouterOps) Create(
	ctx context.Context,
	req *mantraev1.CreateRouterRequest,
) (*mantraev1.CreateRouterResponse, error) {
	params := db.CreateUdpRouterParams{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileId,
		Name:      req.Name,
		AgentID:   req.AgentId,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.UDPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	result, err := s.app.Conn.GetQuery().CreateUdpRouter(ctx, params)
	if err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: result.ToProto(),
		},
	})
	return &mantraev1.CreateRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *UDPRouterOps) Update(
	ctx context.Context,
	req *mantraev1.UpdateRouterRequest,
) (*mantraev1.UpdateRouterResponse, error) {
	params := db.UpdateUdpRouterParams{
		ID:      req.Id,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	var err error
	params.Config, err = db.UnmarshalStruct[schema.UDPRouter](req.Config)
	if err != nil {
		return nil, err
	}
	params.Config.Service = params.Name

	result, err := s.app.Conn.GetQuery().UpdateUdpRouter(ctx, params)
	if err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: result.ToProto(),
		},
	})
	return &mantraev1.UpdateRouterResponse{
		Router: result.ToProto(),
	}, nil
}

func (s *UDPRouterOps) Delete(
	ctx context.Context,
	req *mantraev1.DeleteRouterRequest,
) (*mantraev1.DeleteRouterResponse, error) {
	router, err := s.app.Conn.GetQuery().GetUdpRouter(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if router.Config.Service != "" {
		service, err := s.app.Conn.GetQuery().GetUdpServiceByName(ctx, db.GetUdpServiceByNameParams{
			ProfileID: router.ProfileID,
			Name:      router.Config.Service,
		})
		if err != nil {
			slog.Error("failed to get udp service", "err", err)
		}
		if err := s.app.Conn.GetQuery().DeleteUdpService(ctx, service.ID); err != nil {
			slog.Error("failed to delete udp service", "err", err)
		}
	}

	if err := s.app.Conn.GetQuery().DeleteUdpRouter(ctx, req.Id); err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Router{
			Router: router.ToProto(),
		},
	})
	return &mantraev1.DeleteRouterResponse{}, nil
}

func (s *UDPRouterOps) List(
	ctx context.Context,
	req *mantraev1.ListRoutersRequest,
) (*mantraev1.ListRoutersResponse, error) {
	result, err := s.app.Conn.GetQuery().
		ListUdpRouters(ctx, db.ListUdpRoutersParams{
			ProfileID: req.ProfileId,
			AgentID:   req.AgentId,
			Limit:     req.Limit,
			Offset:    req.Offset,
		})
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.GetQuery().CountUdpRouters(ctx, db.CountUdpRoutersParams{
		ProfileID: req.ProfileId,
		AgentID:   req.AgentId,
	})
	if err != nil {
		return nil, err
	}

	routers := make([]*mantraev1.Router, 0, len(result))
	for _, r := range result {
		routers = append(routers, r.ToProto())
	}
	return &mantraev1.ListRoutersResponse{
		Routers:    routers,
		TotalCount: totalCount,
	}, nil
}
