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
		Entrypoints: make([]Entrypoint, 0),
		Routers:     make(map[string]Router),
		Services:    make(map[string]Service),
		Middlewares: make(map[string]Middleware),
		Version:     "",
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
	if config.Version != nil {
		data.Version = *config.Version
	}
	return data, nil
}

// UpdateConfig updates and verifies the data coming in
func UpdateConfig(profileID int64, data *Dynamic) error {
	for key, r := range data.Routers {
		if err := r.Verify(); err != nil {
			return fmt.Errorf("router %s: %w", r.Name, err)
		}
		if key != r.Name {
			delete(data.Routers, key)
		}
		data.Routers[r.Name] = r
	}
	for key, s := range data.Services {
		if err := s.Verify(); err != nil {
			return fmt.Errorf("service %s: %w", s.Name, err)
		}
		if key != s.Name {
			delete(data.Services, key)
		}
		data.Services[s.Name] = s
	}
	for key, m := range data.Middlewares {
		if err := m.Verify(); err != nil {
			return fmt.Errorf("middleware %s: %w", m.Name, err)
		}
		if key != m.Name {
			delete(data.Middlewares, key)
		}
		data.Middlewares[m.Name] = m
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
	if _, err := db.Query.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ProfileID:   profileID,
		Entrypoints: entrypoints,
		Routers:     routers,
		Services:    services,
		Middlewares: middlewares,
		Version:     &data.Version,
	}); err != nil {
		return err
	}

	return nil
}
