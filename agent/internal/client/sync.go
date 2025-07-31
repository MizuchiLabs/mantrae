package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/docker/docker/api/types/container"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type SyncJob struct {
	config           *Config
	activeIP         string
	routerClient     mantraev1connect.RouterServiceClient
	serviceClient    mantraev1connect.ServiceServiceClient
	middlewareClient mantraev1connect.MiddlewareServiceClient
	syncedResources  map[string]struct{}
}

func newResourceSyncer(cfg *Config, activeIP string) *SyncJob {
	httpClient := &http.Client{Timeout: cfg.ConnectionTimeout}
	interceptor := authInterceptor(cfg)

	return &SyncJob{
		config:   cfg,
		activeIP: activeIP,
		routerClient: mantraev1connect.NewRouterServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		serviceClient: mantraev1connect.NewServiceServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		middlewareClient: mantraev1connect.NewMiddlewareServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		syncedResources: make(map[string]struct{}),
	}
}

func (s *SyncJob) processContainer(ctx context.Context, container container.Summary) error {
	dyn := parseTraefikConfig(container.Labels)
	if dyn == nil {
		return nil
	}
	port := container.Ports[0].PublicPort
	enhanceData(dyn, s.activeIP, port)

	// Process all resource types
	resourceTypes := []struct {
		protocol    mantraev1.ProtocolType
		routers     map[string]any
		services    map[string]any
		middlewares map[string]any
	}{
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
			toAnyMap(dyn.HTTP.Routers),
			toAnyMap(dyn.HTTP.Services),
			toAnyMap(dyn.HTTP.Middlewares),
		},
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
			toAnyMap(dyn.TCP.Routers),
			toAnyMap(dyn.TCP.Services),
			toAnyMap(dyn.TCP.Middlewares),
		},
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_UDP,
			toAnyMap(dyn.UDP.Routers),
			toAnyMap(dyn.UDP.Services),
			nil,
		},
	}

	for _, rt := range resourceTypes {
		if err := s.syncResources(ctx, rt.protocol, rt.routers, rt.services, rt.middlewares); err != nil {
			return err
		}
	}

	return nil
}

func parseTraefikConfig(labels map[string]string) *dynamic.Configuration {
	dyn := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{},
		TCP:  &dynamic.TCPConfiguration{},
		UDP:  &dynamic.UDPConfiguration{},
		TLS:  &dynamic.TLSConfiguration{},
	}

	if err := parser.Decode(labels, dyn, parser.DefaultRootName,
		"traefik.http", "traefik.tcp", "traefik.udp", "traefik.tls.stores.default"); err != nil {
		slog.Error("Failed to parse Traefik config", "error", err)
		return nil
	}

	return dyn
}

func (s *SyncJob) syncResources(ctx context.Context, protocol mantraev1.ProtocolType,
	routers, services, middlewares map[string]any,
) error {
	if routers != nil {
		if err := s.upsertResources(ctx, "router", protocol, routers); err != nil {
			return err
		}
		if err := s.upsertResources(ctx, "service", protocol, services); err != nil {
			return err
		}
	}
	if middlewares != nil {
		if err := s.upsertResources(ctx, "middleware", protocol, middlewares); err != nil {
			return err
		}
	}
	return nil
}

// Generic resource upsert function
func (s *SyncJob) upsertResources(ctx context.Context, resourceType string,
	protocol mantraev1.ProtocolType, resources map[string]any,
) error {
	for name, cfg := range resources {
		key := fmt.Sprintf("%s_%s_%s", resourceType, protocol.String(), name)
		s.syncedResources[key] = struct{}{}

		config, err := toProtoStruct(cfg)
		if err != nil {
			return err
		}
		if resourceType == "router" {
			config.Fields["service"] = &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: name},
			}
		}

		if err := s.upsertResource(ctx, resourceType, protocol, name, config); err != nil {
			return err
		}
	}
	return nil
}

func (s *SyncJob) upsertResource(ctx context.Context, resourceType string,
	protocol mantraev1.ProtocolType, name string, config *structpb.Struct,
) error {
	switch resourceType {
	case "router":
		return s.upsertRouter(ctx, protocol, name, config)
	case "service":
		return s.upsertService(ctx, protocol, name, config)
	case "middleware":
		return s.upsertMiddleware(ctx, protocol, name, config)
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

func (s *SyncJob) upsertRouter(
	ctx context.Context,
	protocol mantraev1.ProtocolType,
	name string,
	config *structpb.Struct,
) error {
	// List existing routers
	res, err := s.routerClient.ListRouters(ctx, connect.NewRequest(&mantraev1.ListRoutersRequest{
		ProfileId: s.config.ProfileID,
		AgentId:   &s.config.AgentID,
		Type:      &protocol,
	}))
	if err != nil {
		return err
	}

	// Check if router exists
	for _, r := range res.Msg.Routers {
		if r.Name == name {
			if proto.Equal(r.Config, config) {
				slog.Debug("Skipped updating router", "name", name, "id", r.Id)
				return nil
			}
			// Update existing router
			if _, err = s.routerClient.UpdateRouter(
				ctx,
				connect.NewRequest(&mantraev1.UpdateRouterRequest{
					Id:      r.Id,
					Name:    name,
					Config:  config,
					Type:    protocol,
					Enabled: true,
				}),
			); err != nil {
				return err
			}
			slog.Debug("Updated router", "name", name, "id", r.Id)
			return nil
		}
	}

	// Create new router
	createRes, err := s.routerClient.CreateRouter(
		ctx,
		connect.NewRequest(&mantraev1.CreateRouterRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
			Name:      name,
			Config:    config,
			Type:      protocol,
			Enabled:   true,
		}),
	)
	if err != nil {
		return err
	}
	slog.Debug("Created router", "name", createRes.Msg.Router.Name, "id", createRes.Msg.Router.Id)
	return nil
}

func (s *SyncJob) upsertService(
	ctx context.Context,
	protocol mantraev1.ProtocolType,
	name string,
	config *structpb.Struct,
) error {
	// List existing services
	res, err := s.serviceClient.ListServices(ctx, connect.NewRequest(&mantraev1.ListServicesRequest{
		ProfileId: s.config.ProfileID,
		AgentId:   &s.config.AgentID,
		Type:      &protocol,
	}))
	if err != nil {
		return err
	}

	// Check if service exists
	for _, svc := range res.Msg.Services {
		if svc.Name == name {
			if proto.Equal(svc.Config, config) {
				slog.Debug("Skipped updating service", "name", name, "id", svc.Id)
				return nil
			}
			// Update existing service
			if _, err = s.serviceClient.UpdateService(
				ctx,
				connect.NewRequest(&mantraev1.UpdateServiceRequest{
					Id:      svc.Id,
					Name:    name,
					Config:  config,
					Type:    protocol,
					Enabled: true,
				}),
			); err != nil {
				return err
			}
			slog.Debug("Updated service", "name", name, "id", svc.Id)
			return nil
		}
	}

	// Create new service
	createRes, err := s.serviceClient.CreateService(
		ctx,
		connect.NewRequest(&mantraev1.CreateServiceRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
			Name:      name,
			Config:    config,
			Type:      protocol,
			Enabled:   true,
		}),
	)
	if err != nil {
		return err
	}
	slog.Debug(
		"Created service",
		"name",
		createRes.Msg.Service.Name,
		"id",
		createRes.Msg.Service.Id,
	)
	return nil
}

func (s *SyncJob) upsertMiddleware(
	ctx context.Context,
	protocol mantraev1.ProtocolType,
	name string,
	config *structpb.Struct,
) error {
	// List existing middlewares
	res, err := s.middlewareClient.ListMiddlewares(
		ctx,
		connect.NewRequest(&mantraev1.ListMiddlewaresRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
			Type:      &protocol,
		}),
	)
	if err != nil {
		return err
	}

	// Check if middleware exists
	for _, mw := range res.Msg.Middlewares {
		if mw.Name == name {
			if proto.Equal(mw.Config, config) {
				slog.Debug("Skipped updating middleware", "name", name, "id", mw.Id)
				return nil
			}
			// Update existing middleware
			if _, err = s.middlewareClient.UpdateMiddleware(
				ctx,
				connect.NewRequest(&mantraev1.UpdateMiddlewareRequest{
					Id:      mw.Id,
					Name:    name,
					Config:  config,
					Type:    protocol,
					Enabled: true,
				}),
			); err != nil {
				return err
			}
			slog.Debug("Updated middleware", "name", name, "id", mw.Id)
			return nil
		}
	}

	// Create new middleware
	createRes, err := s.middlewareClient.CreateMiddleware(
		ctx,
		connect.NewRequest(&mantraev1.CreateMiddlewareRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
			Name:      name,
			Config:    config,
			Type:      protocol,
		}),
	)
	if err != nil {
		return err
	}
	slog.Debug(
		"Created middleware",
		"name",
		createRes.Msg.Middleware.Name,
		"id",
		createRes.Msg.Middleware.Id,
	)
	return nil
}

func (s *SyncJob) cleanup(ctx context.Context) error {
	// Cleanup stale routers
	resRouter, err := s.routerClient.ListRouters(
		ctx,
		connect.NewRequest(&mantraev1.ListRoutersRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
		}),
	)
	if err != nil {
		return err
	}

	for _, r := range resRouter.Msg.Routers {
		key := fmt.Sprintf("router_%s_%s", r.Type.String(), r.Name)
		if _, synced := s.syncedResources[key]; !synced {
			if _, err = s.routerClient.DeleteRouter(ctx, connect.NewRequest(&mantraev1.DeleteRouterRequest{
				Id:   r.Id,
				Type: r.Type,
			})); err != nil {
				slog.Error("Failed to delete stale router", "name", r.Name, "error", err)
			} else {
				slog.Info("Deleted stale router", "name", r.Name)
			}
		}
	}

	// Cleanup stale services
	resService, err := s.serviceClient.ListServices(
		ctx,
		connect.NewRequest(&mantraev1.ListServicesRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
		}),
	)
	if err != nil {
		return err
	}

	for _, svc := range resService.Msg.Services {
		key := fmt.Sprintf("service_%s_%s", svc.Type.String(), svc.Name)
		if _, synced := s.syncedResources[key]; !synced {
			if _, err = s.serviceClient.DeleteService(ctx, connect.NewRequest(&mantraev1.DeleteServiceRequest{
				Id:   svc.Id,
				Type: svc.Type,
			})); err != nil {
				slog.Error("Failed to delete stale service", "name", svc.Name, "error", err)
			} else {
				slog.Info("Deleted stale service", "name", svc.Name)
			}
		}
	}

	// Cleanup stale middlewares
	resMW, err := s.middlewareClient.ListMiddlewares(
		ctx,
		connect.NewRequest(&mantraev1.ListMiddlewaresRequest{
			ProfileId: s.config.ProfileID,
			AgentId:   &s.config.AgentID,
		}),
	)
	if err != nil {
		return err
	}

	for _, mw := range resMW.Msg.Middlewares {
		key := fmt.Sprintf("middleware_%s_%s", mw.Type.String(), mw.Name)
		if _, synced := s.syncedResources[key]; !synced {
			if _, err := s.middlewareClient.DeleteMiddleware(ctx, connect.NewRequest(&mantraev1.DeleteMiddlewareRequest{
				Id:   mw.Id,
				Type: mw.Type,
			})); err != nil {
				slog.Error("Failed to delete stale middleware", "name", mw.Name, "error", err)
			} else {
				slog.Info("Deleted stale middleware", "name", mw.Name)
			}
		}
	}

	return nil
}

// enhanceData adds extra data to the dynamic configuration
func enhanceData(d *dynamic.Configuration, ip string, port uint16) {
	fallbackPort := strconv.Itoa(int(port))

	for _, svc := range d.HTTP.Services {
		for i := range svc.LoadBalancer.Servers {
			if svc.LoadBalancer.Servers[i].Port == "" {
				svc.LoadBalancer.Servers[i].Port = fallbackPort
			}
			svc.LoadBalancer.Servers[i].URL = fmt.Sprintf(
				"http://%s:%s",
				ip,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}
	for _, svc := range d.TCP.Services {
		for i := range svc.LoadBalancer.Servers {
			if svc.LoadBalancer.Servers[i].Port == "" {
				svc.LoadBalancer.Servers[i].Port = fallbackPort
			}
			svc.LoadBalancer.Servers[i].Address = fmt.Sprintf(
				"%s:%s",
				ip,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}
	for _, svc := range d.UDP.Services {
		for i := range svc.LoadBalancer.Servers {
			if svc.LoadBalancer.Servers[i].Port == "" {
				svc.LoadBalancer.Servers[i].Port = fallbackPort
			}
			svc.LoadBalancer.Servers[i].Address = fmt.Sprintf(
				"%s:%s",
				ip,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}
}

// Helper functions remain the same but simplified
func toAnyMap[T any](in map[string]T) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func toProtoStruct(v any) (*structpb.Struct, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var mapData map[string]any
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}

	return structpb.NewStruct(mapData)
}
