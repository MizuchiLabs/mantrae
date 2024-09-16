package traefik

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// DecodeConfig decodes the config from the database into our Dynamic struct
func DecodeConfig(config db.Config) (*Dynamic, error) {
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

// UpdateConfig updates and verifies the data coming in
func UpdateConfig(profileID int64, data *Dynamic) error {
	// Verify and handle routers
	for _, r := range data.Routers {
		if err := r.Verify(); err != nil {
			return fmt.Errorf("router %s: %w", r.Name, err)
		}
		data.Routers[r.Name] = r
	}

	// Verify and handle services
	for _, s := range data.Services {
		if err := s.Verify(); err != nil {
			return fmt.Errorf("service %s: %w", s.Name, err)
		}
		data.Services[s.Name] = s
	}

	// Verify and handle middlewares
	for _, m := range data.Middlewares {
		if err := m.Verify(); err != nil {
			return fmt.Errorf("middleware %s: %w", m.Name, err)
		}
		data.Middlewares[m.Name] = m
	}

	overview, err := json.Marshal(data.Overview)
	if err != nil {
		return err
	}
	entrypoints, err := json.Marshal(data.Entrypoints)
	if err != nil {
		return err
	}
	routers, err := json.Marshal(data.Routers)
	if err != nil {
		return err
	}
	services, err := json.Marshal(data.Services)
	if err != nil {
		return err
	}
	middlewares, err := json.Marshal(data.Middlewares)
	if err != nil {
		return err
	}
	tls, err := json.Marshal(data.TLS)
	if err != nil {
		return err
	}
	if _, err := db.Query.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ProfileID:   profileID,
		Overview:    overview,
		Entrypoints: entrypoints,
		Routers:     routers,
		Services:    services,
		Middlewares: middlewares,
		Tls:         tls,
		Version:     &data.Version,
	}); err != nil {
		return err
	}

	return nil
}
