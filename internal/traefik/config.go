package traefik

import (
	"context"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// BuildDynamicConfig builds a Traefik configuration from the database
func BuildDynamicConfig(
	ctx context.Context,
	q *db.Queries,
	profileName string,
) (*dynamic.Configuration, error) {
	profile, err := q.GetProfileByName(ctx, profileName)
	if err != nil {
		return nil, err
	}

	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:     make(map[string]*dynamic.TCPRouter),
			Middlewares: make(map[string]*dynamic.TCPMiddleware),
			Services:    make(map[string]*dynamic.TCPService),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
	}

	// Routers
	httpRouters, err := q.ListHttpRouters(
		ctx,
		db.ListHttpRoutersParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}
	tcpRouters, err := q.ListTcpRouters(
		ctx,
		db.ListTcpRoutersParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}
	udpRouters, err := q.ListUdpRouters(
		ctx,
		db.ListUdpRoutersParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}

	// Services
	httpServices, err := q.ListHttpServices(
		ctx,
		db.ListHttpServicesParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}
	tcpServices, err := q.ListTcpServices(
		ctx,
		db.ListTcpServicesParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}
	udpServices, err := q.ListUdpServices(
		ctx,
		db.ListUdpServicesParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}

	// Middlewares
	httpMiddlewares, err := q.ListHttpMiddlewares(
		ctx,
		db.ListHttpMiddlewaresParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}
	tcpMiddlewares, err := q.ListTcpMiddlewares(
		ctx,
		db.ListTcpMiddlewaresParams{ProfileID: profile.ID, Limit: -1, Offset: 0},
	)
	if err != nil {
		return nil, err
	}

	for _, r := range httpRouters {
		if !r.Enabled {
			continue
		}
		cfg.HTTP.Routers[r.Name] = r.Config.ToDynamic()
	}
	for _, r := range tcpRouters {
		if !r.Enabled {
			continue
		}
		cfg.TCP.Routers[r.Name] = r.Config.ToDynamic()
	}
	for _, r := range udpRouters {
		if !r.Enabled {
			continue
		}
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
		if !m.Enabled {
			continue
		}
		cfg.HTTP.Middlewares[m.Name] = m.Config.ToDynamic()
	}
	for _, m := range tcpMiddlewares {
		if !m.Enabled {
			continue
		}
		cfg.TCP.Middlewares[m.Name] = m.Config.ToDynamic()
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
