package traefik

import (
	"context"
	"encoding/json"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

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

func UpdateConfig(profileID int64, data *Dynamic) error {
	for _, r := range data.Routers {
		if err := r.Verify(); err != nil {
			return err
		}
	}
	for _, s := range data.Services {
		if err := s.Verify(); err != nil {
			return err
		}
	}
	for _, m := range data.Middlewares {
		if err := m.Verify(); err != nil {
			return err
		}
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
