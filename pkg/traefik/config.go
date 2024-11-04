package traefik

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	ttls "github.com/traefik/traefik/v3/pkg/tls"
)

func GenerateConfig(id int64) (*dynamic.Configuration, error) {
	config := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:           make(map[string]*dynamic.Router),
			Middlewares:       make(map[string]*dynamic.Middleware),
			Services:          make(map[string]*dynamic.Service),
			ServersTransports: make(map[string]*dynamic.ServersTransport),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:     make(map[string]*dynamic.TCPRouter),
			Services:    make(map[string]*dynamic.TCPService),
			Middlewares: make(map[string]*dynamic.TCPMiddleware),
		},
		UDP: &dynamic.UDPConfiguration{
			Routers:  make(map[string]*dynamic.UDPRouter),
			Services: make(map[string]*dynamic.UDPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]ttls.Store),
			Options: make(map[string]ttls.Options),
		},
	}

	routers, err := db.Query.ListRoutersByProfileID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	for _, router := range routers {
		if err := router.DecodeFields(); err != nil {
			continue
		}

		// Only add routers by our provider
		if router.Provider == "http" {
			name := strings.Split(router.Name, "@")[0]
			switch router.Protocol {
			case "http":
				var rtr *dynamic.Router
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.HTTP.Routers[name] = rtr
			case "tcp":
				var rtr *dynamic.TCPRouter
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}

				config.TCP.Routers[name] = rtr
			case "udp":
				var rtr *dynamic.UDPRouter
				rtrBytes, err := json.Marshal(router)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal router: %w", err)
				}
				if err := json.Unmarshal(rtrBytes, &rtr); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.UDP.Routers[name] = rtr
			}
		}
	}

	services, err := db.Query.ListServicesByProfileID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if err := service.DecodeFields(); err != nil {
			continue
		}

		// Only add services by our provider
		if service.Provider == "http" {
			name := strings.Split(service.Name, "@")[0]
			switch service.Protocol {
			case "http":
				var svc *dynamic.Service
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.HTTP.Services[name] = svc
			case "tcp":
				var svc *dynamic.TCPService
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.TCP.Services[name] = svc
			case "udp":
				var svc *dynamic.UDPService
				svcBytes, err := json.Marshal(service)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal service: %w", err)
				}
				if err := json.Unmarshal(svcBytes, &svc); err != nil {
					return nil, fmt.Errorf("failed to unmarshal service: %w", err)
				}
				config.UDP.Services[name] = svc
			}
		}
	}

	middlewares, err := db.Query.ListMiddlewaresByProfileID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	for _, middleware := range middlewares {
		if err := middleware.DecodeFields(); err != nil {
			continue
		}

		if middleware.Provider == "http" {
			name := strings.Split(middleware.Name, "@")[0]
			switch middleware.Protocol {
			case "http":
				var mw *dynamic.Middleware
				mwBytes, err := json.Marshal(middleware)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal middleware: %w", err)
				}
				if err := json.Unmarshal(mwBytes, &mw); err != nil {
					return nil, fmt.Errorf("failed to unmarshal middleware: %w", err)
				}

				// Handle Content based on the Type field
				var contentBytes []byte
				if middleware.Content != nil {
					contentBytes, err = json.Marshal(middleware.Content)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal middleware content: %w", err)
					}
				}

				// Use reflection to assign Content to the correct field based on Type
				contentValue := reflect.New(reflect.TypeOf(middleware.Content)).Interface()
				if err := json.Unmarshal(contentBytes, contentValue); err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal content for type %s: %w",
						middleware.Type,
						err,
					)
				}

				pluginContent := make(map[string]dynamic.PluginConf)

				// Set the content in the right place using reflection
				mwValue := reflect.ValueOf(mw).Elem()
				for i := 0; i < mwValue.NumField(); i++ {
					field := mwValue.Type().Field(i)
					// Compare using EqualFold for case-insensitive matching
					if strings.EqualFold(field.Name, middleware.Type) {
						// Create an instance of the correct type
						fieldType := field.Type
						if fieldType.Kind() == reflect.Ptr {
							// Dereference pointer to get underlying type
							fieldType = fieldType.Elem()
						}

						// Handle plugin case separately
						if strings.EqualFold(middleware.Type, "plugin") {
							// Initialize pluginContent as a map to hold the dynamic.PluginConf
							pluginConf := make(dynamic.PluginConf)

							// Unmarshal contentBytes directly into the pluginContent map
							if err := json.Unmarshal(contentBytes, &pluginConf); err != nil {
								return nil, fmt.Errorf(
									"failed to unmarshal content for type %s: %w",
									middleware.Type,
									err,
								)
							}
							pluginContent[middleware.Name] = pluginConf
							mwValue.Field(i).Set(reflect.ValueOf(pluginContent))
						} else {
							// Create an instance of the correct type
							contentValue = reflect.New(fieldType).Interface()

							// Unmarshal contentBytes into the correct type
							if err := json.Unmarshal(contentBytes, contentValue); err != nil {
								return nil, fmt.Errorf("failed to unmarshal content for type %s: %w", middleware.Type, err)
							}

							// Set the value in the Middleware struct
							mwValue.Field(i).Set(reflect.ValueOf(contentValue))
						}

						break
					}
				}

				config.HTTP.Middlewares[name] = mw
			case "tcp":
				var mw *dynamic.TCPMiddleware
				mwBytes, err := json.Marshal(middleware)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal middleware: %w", err)
				}
				if err := json.Unmarshal(mwBytes, &mw); err != nil {
					return nil, fmt.Errorf("failed to unmarshal middleware: %w", err)
				}
				// Handle Content based on the Type field
				var contentBytes []byte
				if middleware.Content != nil {
					contentBytes, err = json.Marshal(middleware.Content)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal middleware content: %w", err)
					}
				}

				// Use reflection to assign Content to the correct field based on Type
				contentValue := reflect.New(reflect.TypeOf(middleware.Content)).Interface()
				if err := json.Unmarshal(contentBytes, contentValue); err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal content for type %s: %w",
						middleware.Type,
						err,
					)
				}

				// Set the content in the right place using reflection
				mwValue := reflect.ValueOf(mw).Elem()
				for i := 0; i < mwValue.NumField(); i++ {
					field := mwValue.Type().Field(i)
					// Compare using EqualFold for case-insensitive matching
					if strings.EqualFold(field.Name, middleware.Type) {
						// Create an instance of the correct type
						fieldType := field.Type
						if fieldType.Kind() == reflect.Ptr {
							// Dereference pointer to get underlying type
							fieldType = fieldType.Elem()
						}

						// Create an instance of the correct type
						contentValue = reflect.New(fieldType).Interface()

						// Unmarshal contentBytes into the correct type
						if err := json.Unmarshal(contentBytes, contentValue); err != nil {
							return nil, fmt.Errorf(
								"failed to unmarshal content for type %s: %w",
								middleware.Type,
								err,
							)
						}

						// Set the value in the Middleware struct
						mwValue.Field(i).Set(reflect.ValueOf(contentValue))
						break
					}
				}

				config.TCP.Middlewares[name] = mw
			}
		}
	}
	cleanConfig(config)

	return config, nil
}

// cleanConfig removes empty configurations
func cleanConfig(config *dynamic.Configuration) {
	if len(config.HTTP.Routers) == 0 && len(config.HTTP.Services) == 0 &&
		len(config.HTTP.Middlewares) == 0 {
		config.HTTP = nil
	}
	if len(config.TCP.Routers) == 0 && len(config.TCP.Services) == 0 &&
		len(config.TCP.Middlewares) == 0 {
		config.TCP = nil
	}
	if len(config.UDP.Routers) == 0 && len(config.UDP.Services) == 0 {
		config.UDP = nil
	}
	if len(config.TLS.Stores) == 0 && len(config.TLS.Options) == 0 {
		config.TLS = nil
	}
}
