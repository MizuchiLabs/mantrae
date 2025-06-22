package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
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
			func() map[string]any {
				m := make(map[string]any, len(dyn.HTTP.Routers))
				for k, v := range dyn.HTTP.Routers {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
		if err := t.upsertRouters(
			ctx,
			routerClient,
			mantraev1.RouterType_ROUTER_TYPE_TCP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.TCP.Routers))
				for k, v := range dyn.TCP.Routers {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
		if err := t.upsertRouters(
			ctx,
			routerClient,
			mantraev1.RouterType_ROUTER_TYPE_UDP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.UDP.Routers))
				for k, v := range dyn.UDP.Routers {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}

		// Services -----------------------------------------------------------
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_HTTP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.HTTP.Services))
				for k, v := range dyn.HTTP.Services {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_TCP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.TCP.Services))
				for k, v := range dyn.TCP.Services {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
		if err := t.upsertServices(
			ctx,
			serviceClient,
			mantraev1.ServiceType_SERVICE_TYPE_UDP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.UDP.Services))
				for k, v := range dyn.UDP.Services {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}

		// Middlewares --------------------------------------------------------
		if err := t.upsertMiddlewares(
			ctx,
			middlewareClient,
			mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.HTTP.Middlewares))
				for k, v := range dyn.HTTP.Middlewares {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
		if err := t.upsertMiddlewares(
			ctx,
			middlewareClient,
			mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
			func() map[string]any {
				m := make(map[string]any, len(dyn.TCP.Middlewares))
				for k, v := range dyn.TCP.Middlewares {
					m[k] = v
				}
				return m
			}(),
		); err != nil {
			return err
		}
	}

	return nil
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
) error {
	listReq := connect.NewRequest(
		&mantraev1.ListRoutersRequest{ProfileId: t.claims.ProfileID, Type: &typ},
	)
	listResp, err := client.ListRouters(ctx, listReq)
	if err != nil {
		return err
	}
	existing := make(map[string]int64, len(listResp.Msg.Routers))
	for _, r := range listResp.Msg.Routers {
		existing[r.Name] = r.Id
	}

	for name, cfg := range routers {
		s, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}
		if id, found := existing[name]; found {
			upd := &mantraev1.UpdateRouterRequest{
				Id:      id,
				Name:    name,
				Config:  s,
				Enabled: true,
				Type:    typ,
			}
			if _, err := client.UpdateRouter(ctx, connect.NewRequest(upd)); err != nil {
				return err
			}
		} else {
			cr := &mantraev1.CreateRouterRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    s,
				Enabled:   true,
				Type:      typ,
			}
			if _, err := client.CreateRouter(ctx, connect.NewRequest(cr)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TokenSource) upsertServices(
	ctx context.Context,
	client mantraev1connect.ServiceServiceClient,
	typ mantraev1.ServiceType,
	services map[string]any,
) error {
	listReq := connect.NewRequest(
		&mantraev1.ListServicesRequest{ProfileId: t.claims.ProfileID, Type: &typ},
	)
	listResp, err := client.ListServices(ctx, listReq)
	if err != nil {
		return err
	}
	existing := make(map[string]int64, len(listResp.Msg.Services))
	for _, s := range listResp.Msg.Services {
		existing[s.Name] = s.Id
	}
	for name, cfg := range services {
		s, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}
		if id, found := existing[name]; found {
			upd := &mantraev1.UpdateServiceRequest{Id: id, Name: name, Config: s, Type: typ}
			if _, err := client.UpdateService(ctx, connect.NewRequest(upd)); err != nil {
				return err
			}
		} else {
			cr := &mantraev1.CreateServiceRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    s,
				Type:      typ,
			}
			if _, err := client.CreateService(ctx, connect.NewRequest(cr)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TokenSource) upsertMiddlewares(
	ctx context.Context,
	client mantraev1connect.MiddlewareServiceClient,
	typ mantraev1.MiddlewareType,
	middlewares map[string]any,
) error {
	listReq := connect.NewRequest(
		&mantraev1.ListMiddlewaresRequest{ProfileId: t.claims.ProfileID, Type: &typ},
	)
	listResp, err := client.ListMiddlewares(ctx, listReq)
	if err != nil {
		return err
	}
	existing := make(map[string]int64, len(listResp.Msg.Middlewares))
	for _, m := range listResp.Msg.Middlewares {
		existing[m.Name] = m.Id
	}
	for name, cfg := range middlewares {
		s, err := ToProtoStruct(cfg)
		if err != nil {
			return err
		}
		if id, found := existing[name]; found {
			upd := &mantraev1.UpdateMiddlewareRequest{Id: id, Name: name, Config: s, Type: typ}
			if _, err := client.UpdateMiddleware(ctx, connect.NewRequest(upd)); err != nil {
				return err
			}
		} else {
			cr := &mantraev1.CreateMiddlewareRequest{
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Name:      name,
				Config:    s,
				Type:      typ,
			}
			if _, err := client.CreateMiddleware(ctx, connect.NewRequest(cr)); err != nil {
				return err
			}
		}
	}
	return nil
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
