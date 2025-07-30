package traefik

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mizuchilabs/mantrae/server/internal/storage"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/server/internal/store/schema"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
)

// BuildDynamicConfig builds a Traefik configuration from the database
func BuildDynamicConfig(
	ctx context.Context,
	q *db.Queries,
	profile db.Profile,
) (*dynamic.Configuration, error) {
	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:           make(map[string]*dynamic.Router),
			Middlewares:       make(map[string]*dynamic.Middleware),
			Services:          make(map[string]*dynamic.Service),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:           make(map[string]*dynamic.TCPRouter),
			Middlewares:       make(map[string]*dynamic.TCPMiddleware),
			Services:          make(map[string]*dynamic.TCPService),
			ServersTransports: make(map[string]*dynamic.TCPServersTransport),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
	}

	// Routers
	httpRouters, err := q.ListHttpRoutersEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	tcpRouters, err := q.ListTcpRoutersEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	udpRouters, err := q.ListUdpRoutersEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	// Services
	httpServices, err := q.ListHttpServicesEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	tcpServices, err := q.ListTcpServicesEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	udpServices, err := q.ListUdpServicesEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	// Middlewares
	httpMiddlewares, err := q.ListHttpMiddlewaresEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	tcpMiddlewares, err := q.ListTcpMiddlewaresEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	// Servers Transports
	httpServersTransports, err := q.ListHttpServersTransportsEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}
	tcpServersTransports, err := q.ListTcpServersTransportsEnabled(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	for _, r := range httpRouters {
		cfg.HTTP.Routers[r.Name] = r.Config.ToDynamic()
	}
	for _, r := range tcpRouters {
		cfg.TCP.Routers[r.Name] = r.Config.ToDynamic()
	}
	for _, r := range udpRouters {
		cfg.UDP.Routers[r.Name] = r.Config.ToDynamic()
	}

	for _, s := range httpServices {
		cfg.HTTP.Services[s.Name] = s.Config.ToDynamic()
	}
	for _, s := range tcpServices {
		cfg.TCP.Services[s.Name] = s.Config.ToDynamic()
	}
	for _, s := range udpServices {
		cfg.UDP.Services[s.Name] = s.Config.ToDynamic()
	}

	for _, m := range httpMiddlewares {
		cfg.HTTP.Middlewares[m.Name] = m.Config.ToDynamic()
	}
	for _, m := range tcpMiddlewares {
		cfg.TCP.Middlewares[m.Name] = m.Config.ToDynamic()
	}

	for _, s := range httpServersTransports {
		cfg.HTTP.ServersTransports[s.Name] = s.Config.ToDynamic()
	}
	for _, s := range tcpServersTransports {
		cfg.TCP.ServersTransports[s.Name] = s.Config.ToDynamic()
	}

	// Cleanup empty sections (to avoid Traefik {} block warnings)
	if len(cfg.HTTP.Routers) == 0 && len(cfg.HTTP.Middlewares) == 0 &&
		len(cfg.HTTP.Services) == 0 {
		cfg.HTTP = nil
	}
	if len(cfg.TCP.Routers) == 0 && len(cfg.TCP.Middlewares) == 0 &&
		len(cfg.TCP.Services) == 0 {
		cfg.TCP = nil
	}
	if len(cfg.UDP.Routers) == 0 && len(cfg.UDP.Services) == 0 {
		cfg.UDP = nil
	}

	return cfg, nil
}

func BackupDynamicConfigs(ctx context.Context, q *db.Queries, store storage.Backend) error {
	profiles, err := q.ListProfiles(ctx, db.ListProfilesParams{})
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		cfg, err := BuildDynamicConfig(ctx, q, profile)
		if err != nil {
			return err
		}
		if cfg == nil || (cfg.HTTP == nil && cfg.TCP == nil && cfg.UDP == nil) {
			continue
		}
		backupName := fmt.Sprintf(
			"backup_%s_%s.yaml",
			profile.Name,
			time.Now().UTC().Format("20060102_150405"),
		)
		tmpFile, err := os.CreateTemp("", "traefik_backup_*")
		if err != nil {
			return err
		}
		defer func() {
			if err = os.Remove(tmpFile.Name()); err != nil {
				slog.Error("failed to remove temp file", "error", err)
			}
			if err = tmpFile.Close(); err != nil {
				slog.Error("failed to close temp file", "error", err)
			}
		}()
		enc := yaml.NewEncoder(tmpFile)
		enc.SetIndent(2)
		if err = enc.Encode(cfg); err != nil {
			return err
		}
		if _, err = tmpFile.Seek(0, 0); err != nil {
			return err
		}
		if err = store.Store(ctx, backupName, tmpFile); err != nil {
			return err
		}
	}
	return nil
}

// DynamicToDB saves a Traefik dynamic configuration to the database
func DynamicToDB(
	ctx context.Context,
	q *db.Queries,
	profileID int64,
	dynamic *dynamic.Configuration,
) {
	if dynamic == nil {
		return
	}

	entryPointSet := make(map[string]db.CreateEntryPointParams)
	addEP := func(pts []string) {
		for _, ep := range pts {
			if ep == "" {
				continue
			}
			entryPointSet[ep] = db.CreateEntryPointParams{
				ProfileID: profileID,
				Name:      ep,
				IsDefault: false,
			}
		}
	}

	// Routers
	if dynamic.HTTP != nil {
		for k, v := range dynamic.HTTP.Routers {
			if v == nil {
				continue
			}
			if _, err := q.CreateHttpRouter(ctx, db.CreateHttpRouterParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapRouter(v),
			}); err != nil {
				slog.Warn("failed to create http router", "err", err)
			}
			addEP(v.EntryPoints)
		}
		for k, v := range dynamic.HTTP.Services {
			if v == nil {
				continue
			}
			if _, err := q.CreateHttpService(ctx, db.CreateHttpServiceParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapService(v),
			}); err != nil {
				slog.Warn("failed to create http service", "err", err)
			}
		}
		for k, v := range dynamic.HTTP.Middlewares {
			if v == nil {
				continue
			}
			if _, err := q.CreateHttpMiddleware(ctx, db.CreateHttpMiddlewareParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapMiddleware(v),
			}); err != nil {
				slog.Warn("failed to create http middleware", "err", err)
			}
		}
	}
	if dynamic.TCP != nil {
		for k, v := range dynamic.TCP.Routers {
			if v == nil {
				continue
			}
			if _, err := q.CreateTcpRouter(ctx, db.CreateTcpRouterParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapTCPRouter(v),
			}); err != nil {
				slog.Warn("failed to create tcp router", "err", err)
			}
			addEP(v.EntryPoints)
		}
		for k, v := range dynamic.TCP.Services {
			if v == nil {
				continue
			}
			if _, err := q.CreateTcpService(ctx, db.CreateTcpServiceParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapTCPService(v),
			}); err != nil {
				slog.Warn("failed to create tcp service", "err", err)
			}
		}
		for k, v := range dynamic.TCP.Middlewares {
			if v == nil {
				continue
			}
			if _, err := q.CreateTcpMiddleware(ctx, db.CreateTcpMiddlewareParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapTCPMiddleware(v),
			}); err != nil {
				slog.Warn("failed to create tcp middleware", "err", err)
			}
		}
	}
	if dynamic.UDP != nil {
		for k, v := range dynamic.UDP.Routers {
			if v == nil {
				continue
			}
			if _, err := q.CreateUdpRouter(ctx, db.CreateUdpRouterParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapUDPRouter(v),
			}); err != nil {
				slog.Warn("failed to create udp router", "err", err)
			}
			addEP(v.EntryPoints)
		}
		for k, v := range dynamic.UDP.Services {
			if v == nil {
				continue
			}
			if _, err := q.CreateUdpService(ctx, db.CreateUdpServiceParams{
				ProfileID: profileID,
				Name:      k,
				Config:    schema.WrapUDPService(v),
			}); err != nil {
				slog.Warn("failed to create udp service", "err", err)
			}
		}
	}

	if len(entryPointSet) > 0 {
		for _, ep := range entryPointSet {
			if _, err := q.CreateEntryPoint(ctx, ep); err != nil {
				slog.Warn("failed to create entry point", "err", err)
			}
		}
	}
}
