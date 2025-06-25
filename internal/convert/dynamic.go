package convert

import (
	"context"
	"log/slog"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// DynamicToDB saves a Traefik dynamic configuration to the database
func DynamicToDB(
	ctx context.Context,
	q db.Queries,
	profileID int64,
	dynamic *dynamic.Configuration,
) error {
	for k, v := range dynamic.HTTP.Routers {
		if _, err := q.CreateHttpRouter(ctx, db.CreateHttpRouterParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapRouter(v),
		}); err != nil {
			slog.Error("failed to create http router", "err", err)
			continue
		}
	}
	for k, v := range dynamic.TCP.Routers {
		if _, err := q.CreateTcpRouter(ctx, db.CreateTcpRouterParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapTCPRouter(v),
		}); err != nil {
			slog.Error("failed to create tcp router", "err", err)
			continue
		}
	}
	for k, v := range dynamic.UDP.Routers {
		if _, err := q.CreateUdpRouter(ctx, db.CreateUdpRouterParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapUDPRouter(v),
		}); err != nil {
			slog.Error("failed to create udp router", "err", err)
			continue
		}
	}

	// Services
	for k, v := range dynamic.HTTP.Services {
		if _, err := q.CreateHttpService(ctx, db.CreateHttpServiceParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapService(v),
		}); err != nil {
			slog.Error("failed to create http service", "err", err)
			continue
		}
	}
	for k, v := range dynamic.TCP.Services {
		if _, err := q.CreateTcpService(ctx, db.CreateTcpServiceParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapTCPService(v),
		}); err != nil {
			slog.Error("failed to create tcp service", "err", err)
			continue
		}
	}
	for k, v := range dynamic.UDP.Services {
		if _, err := q.CreateUdpService(ctx, db.CreateUdpServiceParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapUDPService(v),
		}); err != nil {
			slog.Error("failed to create udp service", "err", err)
			continue
		}
	}

	// Middlewares
	for k, v := range dynamic.HTTP.Middlewares {
		if _, err := q.CreateHttpMiddleware(ctx, db.CreateHttpMiddlewareParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapMiddleware(v),
		}); err != nil {
			slog.Error("failed to create http middleware", "err", err)
			continue
		}
	}
	for k, v := range dynamic.TCP.Middlewares {
		if _, err := q.CreateTcpMiddleware(ctx, db.CreateTcpMiddlewareParams{
			ProfileID: profileID,
			Name:      k,
			Config:    schema.WrapTCPMiddleware(v),
		}); err != nil {
			slog.Error("failed to create tcp middleware", "err", err)
			continue
		}
	}
	return nil
}
