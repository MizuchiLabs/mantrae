package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
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

func Client(ctx context.Context, quit chan os.Signal) {
	t := NewTokenSource()
	if err := t.SetToken(ctx); err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}
	t.PrintConnection()

	// Prepare tickers
	healthTicker := time.NewTicker(15 * time.Second)
	defer healthTicker.Stop()
	containerTicker := time.NewTicker(10 * time.Second)
	defer containerTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			if err := t.Refresh(ctx); err != nil {
				slog.Error("Failed to refresh token", "error", err)
				return
			}
		case <-containerTicker.C:
			t.Update(ctx)
		case <-quit:
			slog.Info("Shutting down agent...")
			return
		}
	}
}

func (t *TokenSource) Update(ctx context.Context) error {
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
		dynConfig := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}

		err := parser.Decode(
			container.Labels,
			dynConfig,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		)
		if err != nil {
			return err
		}

		// Use the first public port
		publicPort := container.Ports[0].PublicPort

		// Routers ------------------------------------------------------------
		for name, config := range dynConfig.HTTP.Routers {
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateRouterRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.RouterType_ROUTER_TYPE_HTTP,
			})
			if _, err := routerClient.CreateRouter(ctx, req); err != nil {
				return err
			}
		}
		for name, config := range dynConfig.TCP.Routers {
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateRouterRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.RouterType_ROUTER_TYPE_TCP,
			})
			if _, err := routerClient.CreateRouter(ctx, req); err != nil {
				return err
			}
		}
		for name, config := range dynConfig.UDP.Routers {
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateRouterRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.RouterType_ROUTER_TYPE_UDP,
			})
			if _, err := routerClient.CreateRouter(ctx, req); err != nil {
				return err
			}
		}

		// Services -----------------------------------------------------------
		for name, config := range dynConfig.HTTP.Services {
			config.LoadBalancer.Servers = []dynamic.Server{{
				URL:  t.activeIP,
				Port: strconv.Itoa(int(publicPort)),
			}}

			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateServiceRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
			})
			if _, err := serviceClient.CreateService(ctx, req); err != nil {
				return err
			}
		}
		for name, config := range dynConfig.TCP.Services {
			config.LoadBalancer.Servers = []dynamic.TCPServer{{
				Address: t.activeIP,
				Port:    strconv.Itoa(int(publicPort)),
			}}
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateServiceRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
			})
			if _, err := serviceClient.CreateService(ctx, req); err != nil {
				return err
			}
		}
		for name, config := range dynConfig.UDP.Services {
			config.LoadBalancer.Servers = []dynamic.UDPServer{{
				Address: t.activeIP,
				Port:    strconv.Itoa(int(publicPort)),
			}}
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateServiceRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
			})
			if _, err := serviceClient.CreateService(ctx, req); err != nil {
				return err
			}
		}

		// Middlewares --------------------------------------------------------
		for name, config := range dynConfig.HTTP.Middlewares {
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateMiddlewareRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
			})
			if _, err := middlewareClient.CreateMiddleware(ctx, req); err != nil {
				return err
			}
		}
		for name, config := range dynConfig.TCP.Middlewares {
			wrapped, err := ToProtoStruct(config)
			if err != nil {
				return err
			}
			req := connect.NewRequest(&mantraev1.CreateMiddlewareRequest{
				Name:      name,
				Config:    wrapped,
				ProfileId: t.claims.ProfileID,
				AgentId:   t.claims.AgentID,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
			})
			if _, err := middlewareClient.CreateMiddleware(ctx, req); err != nil {
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
