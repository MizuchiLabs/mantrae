package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func HTTPMiddlewareToProto(m *db.HttpMiddleware) *mantraev1.Middleware {
	config, err := MarshalStruct(m.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Middleware{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    config,
		Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
		Enabled:   m.Enabled,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func TCPMiddlewareToProto(m *db.TcpMiddleware) *mantraev1.Middleware {
	config, err := MarshalStruct(m.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Middleware{
		Id:        m.ID,
		ProfileId: m.ProfileID,
		AgentId:   SafeString(m.AgentID),
		Name:      m.Name,
		Config:    config,
		Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
		Enabled:   m.Enabled,
		CreatedAt: SafeTimestamp(m.CreatedAt),
		UpdatedAt: SafeTimestamp(m.UpdatedAt),
	}
}

func HTTPMiddlewaresToProto(middlewares []db.HttpMiddleware) []*mantraev1.Middleware {
	var middlewaresProto []*mantraev1.Middleware
	for _, m := range middlewares {
		middlewaresProto = append(middlewaresProto, HTTPMiddlewareToProto(&m))
	}
	return middlewaresProto
}

func TCPMiddlewaresToProto(middlewares []db.TcpMiddleware) []*mantraev1.Middleware {
	var middlewaresProto []*mantraev1.Middleware
	for _, m := range middlewares {
		middlewaresProto = append(middlewaresProto, TCPMiddlewareToProto(&m))
	}
	return middlewaresProto
}

// Specialized batch conversion functions

func MiddlewaresByProfileToProto(
	middlewares []db.ListMiddlewaresByProfileRow,
) []*mantraev1.Middleware {
	var middlewaresProto []*mantraev1.Middleware
	for _, m := range middlewares {
		switch m.Type {
		case "http":
			config, err := MarshalStruct(m.Config)
			if err != nil {
				return nil
			}
			middlewaresProto = append(middlewaresProto, &mantraev1.Middleware{
				Id:        m.ID,
				ProfileId: m.ProfileID,
				AgentId:   SafeString(m.AgentID),
				Name:      m.Name,
				Config:    config,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
				Enabled:   m.Enabled,
				CreatedAt: SafeTimestamp(m.CreatedAt),
				UpdatedAt: SafeTimestamp(m.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(m.Config)
			if err != nil {
				return nil
			}
			middlewaresProto = append(middlewaresProto, &mantraev1.Middleware{
				Id:        m.ID,
				ProfileId: m.ProfileID,
				AgentId:   SafeString(m.AgentID),
				Name:      m.Name,
				Config:    config,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
				Enabled:   m.Enabled,
				CreatedAt: SafeTimestamp(m.CreatedAt),
				UpdatedAt: SafeTimestamp(m.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return middlewaresProto
}

func MiddlewaresByAgentToProto(middlewares []db.ListMiddlewaresByAgentRow) []*mantraev1.Middleware {
	var middlewaresProto []*mantraev1.Middleware
	for _, m := range middlewares {
		switch m.Type {
		case "http":
			config, err := MarshalStruct(m.Config)
			if err != nil {
				return nil
			}
			middlewaresProto = append(middlewaresProto, &mantraev1.Middleware{
				Id:        m.ID,
				ProfileId: m.ProfileID,
				AgentId:   SafeString(m.AgentID),
				Name:      m.Name,
				Config:    config,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_HTTP,
				Enabled:   m.Enabled,
				CreatedAt: SafeTimestamp(m.CreatedAt),
				UpdatedAt: SafeTimestamp(m.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(m.Config)
			if err != nil {
				return nil
			}
			middlewaresProto = append(middlewaresProto, &mantraev1.Middleware{
				Id:        m.ID,
				ProfileId: m.ProfileID,
				AgentId:   SafeString(m.AgentID),
				Name:      m.Name,
				Config:    config,
				Type:      mantraev1.MiddlewareType_MIDDLEWARE_TYPE_TCP,
				Enabled:   m.Enabled,
				CreatedAt: SafeTimestamp(m.CreatedAt),
				UpdatedAt: SafeTimestamp(m.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return middlewaresProto
}
