package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func Client(ctx context.Context) {
	t := NewTokenSource()
	t.Refresh(ctx)
	t.PrintConnection()

	// Prepare tickers
	healthTicker := time.NewTicker(15 * time.Second)
	defer healthTicker.Stop()
	containerTicker := time.NewTicker(10 * time.Second)
	defer containerTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			t.Refresh(ctx)
		case <-containerTicker.C:
			t.Update(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (t *TokenSource) Update(ctx context.Context) error {
	if t.activeIP == "" {
		return nil
	}

	containers, err := collector.GetContainers()
	if err != nil {
		return err
	}
	routerClient := mantraev1connect.NewRouterServiceClient(
		http.DefaultClient,
		t.claims.ServerURL,
		connect.WithInterceptors(t.Interceptor()),
	)
	serviceClient := mantraev1connect.NewServiceServiceClient(
		http.DefaultClient,
		t.claims.ServerURL,
		connect.WithInterceptors(t.Interceptor()),
	)
	middlewareClient := mantraev1connect.NewMiddlewareServiceClient(
		http.DefaultClient,
		t.claims.ServerURL,
		connect.WithInterceptors(t.Interceptor()),
	)

	// Track which resources are synced
	syncedRouters := map[string]struct{}{}
	syncedServices := map[string]struct{}{}
	syncedMiddlewares := map[string]struct{}{}

	// Parse labels and upsert
	for _, container := range containers {
		dyn := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}

		if err := parser.Decode(
			container.Labels,
			dyn,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		); err != nil {
			return err
		}

		// Use the first public port
		port := container.Ports[0].PublicPort
		injectServiceAddresses(dyn, t.activeIP, port)

		// Routers ------------------------------------------------------------
		if err := t.upsertRouters(
			ctx,
			routerClient,
			mantraev1.RouterType_ROUTER_TYPE_HTTP,
			ToAnyMap(dyn.HTTP.Routers),
			syncedRouters,
		); err != nil {
			return err
		}
		if err := t.upsertRouters(
			ctx,
			routerClient,
			mantraev1.RouterType_ROUTER_TYPE_TCP,
			ToAnyMap(dyn.TCP.Routers),
			syncedRouters,
		); err != nil {
			return err
		}
		if err := t.upsertRouters(
			ctx,
			routerClient,
			mantraev1.RouterType_ROUTER_TYPE_UDP,
			ToAnyMap(dyn.UDP.Routers),
			syncedRouters,
		); err != nil {
			return err
		}

		// Services -----------------------------------------------------------
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_HTTP,
			ToAnyMap(dyn.HTTP.Services),
			syncedServices,
		); err != nil {
			return err
		}
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_TCP,
			ToAnyMap(dyn.TCP.Services),
			syncedServices,
		); err != nil {
			return err
		}
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_UDP,
			ToAnyMap(dyn.UDP.Services),
			syncedServices,
		); err != nil {
			return err
		}

		// Middlewares --------------------------------------------------------
		if err := t.upsertMiddlewares(
			ctx,
			middlewareClient,
			mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
			ToAnyMap(dyn.HTTP.Middlewares),
			syncedMiddlewares,
		); err != nil {
			return err
		}
		if err := t.upsertMiddlewares(
			ctx,
			middlewareClient,
			mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
			ToAnyMap(dyn.TCP.Middlewares),
			syncedMiddlewares,
		); err != nil {
			return err
		}
	}

	return t.cleanup(
		ctx,
		routerClient,
		serviceClient,
		middlewareClient,
		syncedRouters,
		syncedServices,
		syncedMiddlewares,
	)
}

func injectServiceAddresses(d *dynamic.Configuration, ip string, port uint16) {
	p := strconv.Itoa(int(port))
	for _, svc := range d.HTTP.Services {
		svc.LoadBalancer.Servers = []dynamic.Server{{URL: ip, Port: p}}
	}
	for _, svc := range d.TCP.Services {
		svc.LoadBalancer.Servers = []dynamic.TCPServer{{Address: ip, Port: p}}
	}
	for _, svc := range d.UDP.Services {
		svc.LoadBalancer.Servers = []dynamic.UDPServer{{Address: ip, Port: p}}
	}
}

func (t *TokenSource) upsertRouters(
	ctx context.Context,
	client mantraev1connect.RouterServiceClient,
	typ mantraev1.RouterType,
	routers map[string]any,
	synced map[string]struct{},
) error {
	res, err := client.ListRouters(ctx, connect.NewRequest(&mantraev1.ListRoutersRequest{
		ProfileId: t.claims.ProfileID,
		AgentId:   &t.claims.AgentID,
		Type:      &typ,
	}))
	if err != nil {
		return err
	}

	existing := make(map[string]*mantraev1.Router, len(res.Msg.Routers))
	for _, r := range res.Msg.Routers {
		existing[r.Name] = r
	}

	for name, cfg := range routers {
		synced[name] = struct{}{}
		newConfig, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}

		if r, found := existing[name]; found {
			if proto.Equal(r.Config, newConfig) {
				slog.Debug("Skipped updating router", "name", name, "id", r.Id)
				continue
			}
			params := &mantraev1.UpdateRouterRequest{
				Id:      r.Id,
				Name:    name,
				Config:  newConfig,
				Enabled: true,
				Type:    typ,
			}
			if _, err := client.UpdateRouter(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Updated router", "name", name, "id", r.Id)
		} else {
			params := &mantraev1.CreateRouterRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    newConfig,
				Enabled:   true,
				Type:      typ,
			}
			if _, err := client.CreateRouter(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Created router", "name", name, "id", r.Id)
		}
	}

	return nil
}

func (t *TokenSource) upsertServices(
	ctx context.Context,
	client mantraev1connect.ServiceServiceClient,
	typ mantraev1.ServiceType,
	services map[string]any,
	synced map[string]struct{},
) error {
	res, err := client.ListServices(ctx, connect.NewRequest(&mantraev1.ListServicesRequest{
		ProfileId: t.claims.ProfileID,
		AgentId:   &t.claims.AgentID,
		Type:      &typ,
	}))
	if err != nil {
		return err
	}

	existing := make(map[string]*mantraev1.Service, len(res.Msg.Services))
	for _, s := range res.Msg.Services {
		existing[s.Name] = s
	}

	for name, cfg := range services {
		synced[name] = struct{}{}
		newConfig, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}

		if s, found := existing[name]; found {
			if proto.Equal(s.Config, newConfig) {
				slog.Debug("Skipped updating service", "name", name, "id", s.Id)
				continue
			}
			params := &mantraev1.UpdateServiceRequest{
				Id:     s.Id,
				Name:   name,
				Config: newConfig,
				Type:   typ,
			}
			if _, err := client.UpdateService(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Updated service", "name", name, "id", s.Id)
		} else {
			params := &mantraev1.CreateServiceRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    newConfig,
				Type:      typ,
			}
			if _, err := client.CreateService(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Created service", "name", name, "id", s.Id)
		}
	}

	return nil
}

func (t *TokenSource) upsertMiddlewares(
	ctx context.Context,
	client mantraev1connect.MiddlewareServiceClient,
	typ mantraev1.MiddlewareType,
	middlewares map[string]any,
	synced map[string]struct{},
) error {
	res, err := client.ListMiddlewares(ctx, connect.NewRequest(&mantraev1.ListMiddlewaresRequest{
		ProfileId: t.claims.ProfileID,
		AgentId:   &t.claims.AgentID,
		Type:      &typ,
	}))
	if err != nil {
		return err
	}

	existing := make(map[string]*mantraev1.Middleware, len(res.Msg.Middlewares))
	for _, m := range res.Msg.Middlewares {
		existing[m.Name] = m
	}

	for name, cfg := range middlewares {
		synced[name] = struct{}{}
		newConfig, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}

		if m, found := existing[name]; found {
			if proto.Equal(m.Config, newConfig) {
				slog.Debug("Skipped updating middleware", "name", name, "id", m.Id)
				continue
			}
			params := &mantraev1.UpdateMiddlewareRequest{
				Id:     m.Id,
				Name:   name,
				Config: newConfig,
				Type:   typ,
			}
			if _, err := client.UpdateMiddleware(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Updated middleware", "name", name, "id", m.Id)
		} else {
			params := &mantraev1.CreateMiddlewareRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    newConfig,
				Type:      typ,
			}
			if _, err := client.CreateMiddleware(ctx, connect.NewRequest(params)); err != nil {
				return err
			}
			slog.Debug("Created middleware", "name", name, "id", m.Id)
		}
	}

	return nil
}

// cleanup removes all stale resources
func (t *TokenSource) cleanup(
	ctx context.Context,
	routerClient mantraev1connect.RouterServiceClient,
	serviceClient mantraev1connect.ServiceServiceClient,
	middlewareClient mantraev1connect.MiddlewareServiceClient,
	syncedRouters map[string]struct{},
	syncedServices map[string]struct{},
	syncedMiddlewares map[string]struct{},
) error {
	// Cleanup Routers
	routers, err := routerClient.ListRouters(ctx, connect.NewRequest(&mantraev1.ListRoutersRequest{
		ProfileId: t.claims.ProfileID,
		AgentId:   &t.claims.AgentID,
	}))
	if err != nil {
		return err
	}
	for _, r := range routers.Msg.Routers {
		if _, ok := syncedRouters[r.Name]; !ok {
			if _, err := routerClient.DeleteRouter(
				ctx,
				connect.NewRequest(&mantraev1.DeleteRouterRequest{Id: r.Id, Type: r.Type}),
			); err != nil {
				slog.Error("Failed to delete stale router", "name", r.Name, "err", err)
			} else {
				slog.Info("Deleted stale router", "name", r.Name)
			}
		}
	}

	// Cleanup Services
	services, err := serviceClient.ListServices(
		ctx,
		connect.NewRequest(&mantraev1.ListServicesRequest{
			ProfileId: t.claims.ProfileID,
			AgentId:   &t.claims.AgentID,
		}),
	)
	if err != nil {
		return err
	}

	for _, s := range services.Msg.Services {
		if _, ok := syncedServices[s.Name]; !ok {
			if _, err := serviceClient.DeleteService(
				ctx,
				connect.NewRequest(&mantraev1.DeleteServiceRequest{Id: s.Id, Type: s.Type}),
			); err != nil {
				slog.Error("Failed to delete stale service", "name", s.Name, "err", err)
			} else {
				slog.Info("Deleted stale service", "name", s.Name)
			}
		}
	}

	// Cleanup Middlewares
	middlewares, err := middlewareClient.ListMiddlewares(
		ctx,
		connect.NewRequest(&mantraev1.ListMiddlewaresRequest{
			ProfileId: t.claims.ProfileID,
			AgentId:   &t.claims.AgentID,
		}),
	)
	if err != nil {
		return err
	}

	for _, m := range middlewares.Msg.Middlewares {
		if _, ok := syncedMiddlewares[m.Name]; !ok {
			if _, err := middlewareClient.DeleteMiddleware(
				ctx,
				connect.NewRequest(&mantraev1.DeleteMiddlewareRequest{Id: m.Id, Type: m.Type}),
			); err != nil {
				slog.Error("Failed to delete stale middleware", "name", m.Name, "err", err)
			} else {
				slog.Info("Deleted stale middleware", "name", m.Name)
			}
		}
	}

	return nil
}

// ToAnyMap converts a map[string]T to map[string]any
func ToAnyMap[T any](in map[string]T) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// ToProtoStruct converts any Go struct to *structpb.Struct
func ToProtoStruct(v any) (*structpb.Struct, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var mapData map[string]interface{}
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, err
	}

	return structpb.NewStruct(mapData)
}
