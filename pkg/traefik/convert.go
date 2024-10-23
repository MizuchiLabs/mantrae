package traefik

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/db"

	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func VerifyConfig(config *Dynamic) {
	for i, r := range config.Routers {
		if err := r.Verify(); err != nil {
			slog.Error("Failed to verify router", "name", r.Name, "error", err)
			delete(config.Routers, i)
		}
	}
	for i, s := range config.Services {
		if err := s.Verify(); err != nil {
			slog.Error("Failed to verify service", "name", s.Name, "error", err)
			delete(config.Services, i)
		}
	}
	for i, m := range config.Middlewares {
		if err := m.Verify(); err != nil {
			slog.Error("Failed to verify middleware", "name", m.Name, "error", err)
			delete(config.Middlewares, i)
		}
	}
}

// DecodeFromDB decodes the config from the database into our Dynamic struct
func DecodeFromDB(config db.Config) (*Dynamic, error) {
	data := &Dynamic{
		ProfileID:   config.ProfileID,
		Overview:    nil,
		Entrypoints: make([]Entrypoint, 0),
		Routers:     make(map[string]Router),
		Services:    make(map[string]Service),
		Middlewares: make(map[string]Middleware),
		TLS:         nil,
		Version:     "",
	}

	if config.Overview != nil {
		if err := json.Unmarshal(config.Overview.([]byte), &data.Overview); err != nil {
			return nil, err
		}
	}
	if config.Entrypoints != nil {
		if err := json.Unmarshal(config.Entrypoints.([]byte), &data.Entrypoints); err != nil {
			return nil, err
		}
	}
	if config.Routers != nil {
		if err := json.Unmarshal(config.Routers.([]byte), &data.Routers); err != nil {
			return nil, err
		}
	}
	if config.Services != nil {
		if err := json.Unmarshal(config.Services.([]byte), &data.Services); err != nil {
			return nil, err
		}
	}
	if config.Middlewares != nil {
		if err := json.Unmarshal(config.Middlewares.([]byte), &data.Middlewares); err != nil {
			return nil, err
		}
	}
	if config.Tls != nil {
		if err := json.Unmarshal(config.Tls.([]byte), &data.TLS); err != nil {
			return nil, err
		}
	}
	if config.Version != nil {
		data.Version = *config.Version
	}
	return data, nil
}

// EncodeToDB encodes the config from our Dynamic struct into the database
func EncodeToDB(profileID int64, data *Dynamic) (*Dynamic, error) {
	VerifyConfig(data)

	overview, err := json.Marshal(data.Overview)
	if err != nil {
		return nil, err
	}
	entrypoints, err := json.Marshal(data.Entrypoints)
	if err != nil {
		return nil, err
	}
	routers, err := json.Marshal(data.Routers)
	if err != nil {
		return nil, err
	}
	services, err := json.Marshal(data.Services)
	if err != nil {
		return nil, err
	}
	middlewares, err := json.Marshal(data.Middlewares)
	if err != nil {
		return nil, err
	}
	tls, err := json.Marshal(data.TLS)
	if err != nil {
		return nil, err
	}
	config, err := db.Query.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ProfileID:   profileID,
		Overview:    overview,
		Entrypoints: entrypoints,
		Routers:     routers,
		Services:    services,
		Middlewares: middlewares,
		Tls:         tls,
		Version:     &data.Version,
	})
	if err != nil {
		return nil, err
	}

	return DecodeFromDB(config)
}

// DecodeFromLabels uses the traefik parses to decode the config from the labels into our Dynamic struct
func DecodeFromLabels(id string) (*Dynamic, error) {
	agent, err := db.Query.GetAgentByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	var containers []agentv1.Container
	if err := json.Unmarshal(agent.Containers.([]byte), &containers); err != nil {
		return nil, err
	}

	dbConfig, err := db.Query.GetConfigByProfileID(context.Background(), 1)
	if err != nil {
		return nil, err
	}

	dynamicConfig, err := DecodeFromDB(dbConfig)
	if err != nil {
		return nil, err
	}

	// delete all by agent id
	for _, router := range dynamicConfig.Routers {
		if router.AgentID == id {
			delete(dynamicConfig.Routers, router.Name)
		}
	}
	for _, service := range dynamicConfig.Services {
		if service.AgentID == id {
			delete(dynamicConfig.Services, service.Name)
		}
	}
	for _, middleware := range dynamicConfig.Middlewares {
		if middleware.AgentID == id {
			delete(dynamicConfig.Middlewares, middleware.Name)
		}
	}

	for i := range containers {
		// Convert labels to official traefik types
		config := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}
		if err := parser.Decode(
			containers[i].Labels,
			config,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		); err != nil {
			return nil, err
		}

		// Add to our dynamic config
		for i, router := range config.HTTP.Routers {
			var tlsConfig *dynamic.RouterTCPTLSConfig
			if router.TLS != nil {
				tlsConfig = &dynamic.RouterTCPTLSConfig{
					Options:      router.TLS.Options,
					CertResolver: router.TLS.CertResolver,
					Domains:      router.TLS.Domains,
				}
			}
			name := strings.Split(i, "@")[0] + "@http"
			dynamicConfig.Routers[name] = Router{
				Name:       name,
				Provider:   "http",
				RouterType: "http",
				// DNSProvider: new(int64),
				AgentID:     id,
				Entrypoints: router.EntryPoints,
				Middlewares: router.Middlewares,
				Rule:        router.Rule,
				RuleSyntax:  router.RuleSyntax,
				Service:     router.Service,
				Priority:    router.Priority,
				TLS:         tlsConfig,
			}
		}

		for i, router := range config.TCP.Routers {
			name := strings.Split(i, "@")[0] + "@http"
			dynamicConfig.Routers[name] = Router{
				Name:       name,
				Provider:   "http",
				RouterType: "tcp",
				// DNSProvider: new(int64),
				AgentID:     id,
				Entrypoints: router.EntryPoints,
				Middlewares: router.Middlewares,
				Rule:        router.Rule,
				RuleSyntax:  router.RuleSyntax,
				Service:     router.Service,
				Priority:    router.Priority,
				TLS:         router.TLS,
			}
		}

		for i, router := range config.UDP.Routers {
			name := strings.Split(i, "@")[0] + "@http"
			dynamicConfig.Routers[name] = Router{
				Name:       name,
				Provider:   "http",
				RouterType: "udp",
				// DNSProvider: new(int64),
				AgentID:     id,
				Entrypoints: router.EntryPoints,
				Service:     router.Service,
			}
		}
	}

	return nil, nil
}
