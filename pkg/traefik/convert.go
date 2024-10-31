package traefik

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// DecodeFromDB decodes the config from the database into our Dynamic struct
func DecodeFromDB(id int64) (*Dynamic, error) {
	if id == 0 {
		return nil, nil
	}

	config, err := db.Query.GetConfigByProfileID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	data := &Dynamic{
		ProfileID:   id,
		Overview:    nil,
		Entrypoints: make([]Entrypoint, 0),
		Version:     "",
		Mutex:       sync.Mutex{},
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
	if config.Version != nil {
		data.Version = *config.Version
	}
	return data, nil
}

// EncodeToDB encodes the config from our Dynamic struct and saves it to the database
func EncodeToDB(data *Dynamic) (*Dynamic, error) {
	if data.ProfileID == 0 {
		return nil, nil
	}
	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	overview, err := json.Marshal(data.Overview)
	if err != nil {
		return nil, err
	}
	entrypoints, err := json.Marshal(data.Entrypoints)
	if err != nil {
		return nil, err
	}
	_, err = db.Query.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ProfileID:   data.ProfileID,
		Overview:    overview,
		Entrypoints: entrypoints,
		Version:     &data.Version,
	})
	if err != nil {
		slog.Error("Failed to update config", "error", err)
		return nil, err
	}

	return DecodeFromDB(data.ProfileID)
}

// DecodeFromLabels uses the traefik parses to decode the config from the labels into our Dynamic struct
func DecodeFromLabels(id string) (*Dynamic, error) {
	// 	agent, err := db.Query.GetAgentByID(context.Background(), id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	var containers []agentv1.Container
	// 	if err := json.Unmarshal(agent.Containers.([]byte), &containers); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	dbConfig, err := db.Query.GetConfigByProfileID(context.Background(), 1)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	dynamicConfig, err := DecodeFromDB(dbConfig.ProfileID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	// delete all by agent id
	// 	for _, router := range dynamicConfig.Routers {
	// 		if router.AgentID == id {
	// 			delete(dynamicConfig.Routers, router.Name)
	// 		}
	// 	}
	// 	for _, service := range dynamicConfig.Services {
	// 		if service.AgentID == id {
	// 			delete(dynamicConfig.Services, service.Name)
	// 		}
	// 	}
	// 	for _, middleware := range dynamicConfig.Middlewares {
	// 		if middleware.AgentID == id {
	// 			delete(dynamicConfig.Middlewares, middleware.Name)
	// 		}
	// 	}
	//
	// 	for i := range containers {
	// 		// Convert labels to official traefik types
	// 		config := &dynamic.Configuration{
	// 			HTTP: &dynamic.HTTPConfiguration{},
	// 			TCP:  &dynamic.TCPConfiguration{},
	// 			UDP:  &dynamic.UDPConfiguration{},
	// 			TLS:  &dynamic.TLSConfiguration{},
	// 		}
	// 		if err := parser.Decode(
	// 			containers[i].Labels,
	// 			config,
	// 			parser.DefaultRootName,
	// 			"traefik.http",
	// 			"traefik.tcp",
	// 			"traefik.udp",
	// 			"traefik.tls.stores.default",
	// 		); err != nil {
	// 			return nil, err
	// 		}
	//
	// 		// Add to our dynamic config
	// 		for i, router := range config.HTTP.Routers {
	// 			var tlsConfig *dynamic.RouterTCPTLSConfig
	// 			if router.TLS != nil {
	// 				tlsConfig = &dynamic.RouterTCPTLSConfig{
	// 					Options:      router.TLS.Options,
	// 					CertResolver: router.TLS.CertResolver,
	// 					Domains:      router.TLS.Domains,
	// 				}
	// 			}
	// 			name := strings.Split(i, "@")[0] + "@http"
	// 			dynamicConfig.Routers[name] = Router{
	// 				Name:       name,
	// 				Provider:   "http",
	// 				RouterType: "http",
	// 				// DNSProvider: new(int64),
	// 				AgentID:     id,
	// 				Entrypoints: router.EntryPoints,
	// 				Middlewares: router.Middlewares,
	// 				Rule:        router.Rule,
	// 				RuleSyntax:  router.RuleSyntax,
	// 				Service:     router.Service,
	// 				Priority:    router.Priority,
	// 				TLS:         tlsConfig,
	// 			}
	// 		}
	//
	// 		for i, router := range config.TCP.Routers {
	// 			name := strings.Split(i, "@")[0] + "@http"
	// 			dynamicConfig.Routers[name] = Router{
	// 				Name:       name,
	// 				Provider:   "http",
	// 				RouterType: "tcp",
	// 				// DNSProvider: new(int64),
	// 				AgentID:     id,
	// 				Entrypoints: router.EntryPoints,
	// 				Middlewares: router.Middlewares,
	// 				Rule:        router.Rule,
	// 				RuleSyntax:  router.RuleSyntax,
	// 				Service:     router.Service,
	// 				Priority:    router.Priority,
	// 				TLS:         router.TLS,
	// 			}
	// 		}
	//
	// 		for i, router := range config.UDP.Routers {
	// 			name := strings.Split(i, "@")[0] + "@http"
	// 			dynamicConfig.Routers[name] = Router{
	// 				Name:       name,
	// 				Provider:   "http",
	// 				RouterType: "udp",
	// 				// DNSProvider: new(int64),
	// 				AgentID:     id,
	// 				Entrypoints: router.EntryPoints,
	// 				Service:     router.Service,
	// 			}
	// 		}
	// 	}
	//
	return nil, nil
}
