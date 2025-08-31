package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type SyncJob struct {
	config           *Config
	routerClient     mantraev1connect.RouterServiceClient
	serviceClient    mantraev1connect.ServiceServiceClient
	middlewareClient mantraev1connect.MiddlewareServiceClient
	syncedResources  map[string][]string
}

func NewJob(cfg *Config) *SyncJob {
	httpClient := &http.Client{Timeout: cfg.ConnectionTimeout}
	interceptor := authInterceptor(cfg)

	return &SyncJob{
		config: cfg,
		routerClient: mantraev1connect.NewRouterServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		serviceClient: mantraev1connect.NewServiceServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		middlewareClient: mantraev1connect.NewMiddlewareServiceClient(
			httpClient, cfg.ServerURL, connect.WithInterceptors(interceptor)),
		syncedResources: make(map[string][]string),
	}
}

func (s *SyncJob) processContainer(ctx context.Context, container collector.ContainerInfo) error {
	dynamic := s.parseTraefikConfig(container.Labels)
	if dynamic == nil {
		return nil
	}

	// Process all resource types
	resourceTypes := []struct {
		protocol    mantraev1.ProtocolType
		routers     map[string]any
		services    map[string]any
		middlewares map[string]any
	}{
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_HTTP,
			toAnyMap(dynamic.HTTP.Routers),
			toAnyMap(dynamic.HTTP.Services),
			toAnyMap(dynamic.HTTP.Middlewares),
		},
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_TCP,
			toAnyMap(dynamic.TCP.Routers),
			toAnyMap(dynamic.TCP.Services),
			toAnyMap(dynamic.TCP.Middlewares),
		},
		{
			mantraev1.ProtocolType_PROTOCOL_TYPE_UDP,
			toAnyMap(dynamic.UDP.Routers),
			toAnyMap(dynamic.UDP.Services),
			nil,
		},
	}

	for _, rt := range resourceTypes {
		if err := s.syncResources(ctx, container.ID, rt.protocol, rt.routers, rt.services, rt.middlewares); err != nil {
			return err
		}
	}

	if err := s.cleanup(ctx); err != nil {
		return err
	}

	return nil
}

func (s *SyncJob) parseTraefikConfig(labels map[string]string) *dynamic.Configuration {
	// Use parser from traefik to parse labels
	d := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{},
		TCP:  &dynamic.TCPConfiguration{},
		UDP:  &dynamic.UDPConfiguration{},
		TLS:  &dynamic.TLSConfiguration{},
	}

	if err := parser.Decode(labels, d, parser.DefaultRootName,
		"traefik.http", "traefik.tcp", "traefik.udp", "traefik.tls.stores.default"); err != nil {
		slog.Error("Failed to parse Traefik config", "error", err)
		return nil
	}

	// Add active IP to load balancer servers
	for _, svc := range d.HTTP.Services {
		for i := range svc.LoadBalancer.Servers {
			svc.LoadBalancer.Servers[i].URL = fmt.Sprintf(
				"http://%s:%s",
				s.config.ActiveIP,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}
	for _, svc := range d.TCP.Services {
		for i := range svc.LoadBalancer.Servers {
			svc.LoadBalancer.Servers[i].Address = fmt.Sprintf(
				"%s:%s",
				s.config.ActiveIP,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}
	for _, svc := range d.UDP.Services {
		for i := range svc.LoadBalancer.Servers {
			svc.LoadBalancer.Servers[i].Address = fmt.Sprintf(
				"%s:%s",
				s.config.ActiveIP,
				svc.LoadBalancer.Servers[i].Port,
			)
		}
	}

	return d
}

func (s *SyncJob) syncResources(
	ctx context.Context,
	containerID string,
	protocol mantraev1.ProtocolType,
	routers, services, middlewares map[string]any,
) error {
	if routers != nil {
		if err := s.upsertResources(ctx, containerID, "router", protocol, routers); err != nil {
			return err
		}
		if err := s.upsertResources(ctx, containerID, "service", protocol, services); err != nil {
			return err
		}
	}
	if middlewares != nil {
		if err := s.upsertResources(ctx, containerID, "middleware", protocol, middlewares); err != nil {
			return err
		}
	}
	return nil
}

// Generic resource upsert function
func (s *SyncJob) upsertResources(ctx context.Context, containerID, resourceType string,
	protocol mantraev1.ProtocolType, resources map[string]any,
) error {
	for name, cfg := range resources {
		key := fmt.Sprintf("%s_%s_%s", resourceType, protocol.String(), name)
		s.syncedResources[containerID] = append(s.syncedResources[containerID], key)

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
	// Get all currently synced resource keys
	syncedKeys := make(map[string]struct{})
	for _, resourceKeys := range s.syncedResources {
		for _, key := range resourceKeys {
			syncedKeys[key] = struct{}{}
		}
	}

	// Cleanup stale routers
	resRouter, err := s.routerClient.ListRouters(
		ctx,
		connect.NewRequest(
			&mantraev1.ListRoutersRequest{
				ProfileId: s.config.ProfileID,
				AgentId:   &s.config.AgentID,
			},
		),
	)
	if err != nil {
		slog.Error("Failed to list routers", "error", err)
		return err
	}

	for _, r := range resRouter.Msg.Routers {
		key := fmt.Sprintf("router_%s_%s", r.Type.String(), r.Name)

		if _, synced := syncedKeys[key]; !synced {
			if _, err = s.routerClient.DeleteRouter(
				ctx,
				connect.NewRequest(&mantraev1.DeleteRouterRequest{Id: r.Id, Type: r.Type}),
			); err != nil {
				slog.Error("Failed to delete stale router", "name", r.Name, "error", err)
				continue
			}

			slog.Info("Deleted stale router", "name", r.Name)
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

		if _, synced := syncedKeys[key]; !synced {
			if _, err = s.serviceClient.DeleteService(
				ctx,
				connect.NewRequest(&mantraev1.DeleteServiceRequest{Id: svc.Id, Type: svc.Type}),
			); err != nil {
				slog.Error("Failed to delete stale service", "name", svc.Name, "error", err)
				continue
			}

			slog.Info("Deleted stale service", "name", svc.Name)
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

	for _, m := range resMW.Msg.Middlewares {
		key := fmt.Sprintf("middleware_%s_%s", m.Type.String(), m.Name)

		if _, synced := syncedKeys[key]; !synced {
			if _, err := s.middlewareClient.DeleteMiddleware(
				ctx,
				connect.NewRequest(&mantraev1.DeleteMiddlewareRequest{Id: m.Id, Type: m.Type}),
			); err != nil {
				slog.Error("Failed to delete stale middleware", "name", m.Name, "error", err)
				continue
			}

			slog.Info("Deleted stale middleware", "name", m.Name)
		}
	}

	return nil
}

func (s *SyncJob) removeContainer(containerID string) {
	delete(s.syncedResources, containerID)
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
