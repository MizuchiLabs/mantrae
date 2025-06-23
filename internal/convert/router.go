package convert

import (
	"context"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func HTTPRouterToProto(r *db.HttpRouter, d []*mantraev1.DnsProvider) *mantraev1.Router {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Router{
		Id:           r.ID,
		ProfileId:    r.ProfileID,
		AgentId:      SafeString(r.AgentID),
		Name:         r.Name,
		Config:       config,
		Enabled:      r.Enabled,
		Type:         mantraev1.RouterType_ROUTER_TYPE_HTTP,
		DnsProviders: d,
		CreatedAt:    SafeTimestamp(r.CreatedAt),
		UpdatedAt:    SafeTimestamp(r.UpdatedAt),
	}
}

func TCPRouterToProto(r *db.TcpRouter, d []*mantraev1.DnsProvider) *mantraev1.Router {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Router{
		Id:           r.ID,
		ProfileId:    r.ProfileID,
		AgentId:      SafeString(r.AgentID),
		Name:         r.Name,
		Config:       config,
		Enabled:      r.Enabled,
		Type:         mantraev1.RouterType_ROUTER_TYPE_TCP,
		DnsProviders: d,
		CreatedAt:    SafeTimestamp(r.CreatedAt),
		UpdatedAt:    SafeTimestamp(r.UpdatedAt),
	}
}

func UDPRouterToProto(r *db.UdpRouter) *mantraev1.Router {
	config, err := MarshalStruct(r.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Router{
		Id:        r.ID,
		ProfileId: r.ProfileID,
		AgentId:   SafeString(r.AgentID),
		Name:      r.Name,
		Config:    config,
		Enabled:   r.Enabled,
		Type:      mantraev1.RouterType_ROUTER_TYPE_UDP,
		CreatedAt: SafeTimestamp(r.CreatedAt),
		UpdatedAt: SafeTimestamp(r.UpdatedAt),
	}
}

func HTTPRoutersToProto(routers []db.HttpRouter, q *db.Queries) []*mantraev1.Router {
	var routersProto []*mantraev1.Router
	for _, r := range routers {
		result, _ := q.GetDnsProvidersByHttpRouter(context.Background(), r.ID)
		routersProto = append(routersProto, HTTPRouterToProto(&r, DNSProvidersToProto(result)))
	}
	return routersProto
}

func TCPRoutersToProto(routers []db.TcpRouter, q *db.Queries) []*mantraev1.Router {
	var routersProto []*mantraev1.Router
	for _, r := range routers {
		result, _ := q.GetDnsProvidersByTcpRouter(context.Background(), r.ID)
		routersProto = append(routersProto, TCPRouterToProto(&r, DNSProvidersToProto(result)))
	}
	return routersProto
}

func UDPRoutersToProto(routers []db.UdpRouter) []*mantraev1.Router {
	var routersProto []*mantraev1.Router
	for _, r := range routers {
		routersProto = append(routersProto, UDPRouterToProto(&r))
	}
	return routersProto
}

// Specialized batch conversion functions
func RoutersByProfileToProto(
	routers []db.ListRoutersByProfileRow,
	q *db.Queries,
) []*mantraev1.Router {
	var routersProto []*mantraev1.Router
	for _, r := range routers {
		switch r.Type {
		case "http":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}

			result, _ := q.GetDnsProvidersByHttpRouter(context.Background(), r.ID)
			routersProto = append(routersProto, &mantraev1.Router{
				Id:           r.ID,
				ProfileId:    r.ProfileID,
				AgentId:      SafeString(r.AgentID),
				Name:         r.Name,
				Config:       config,
				Enabled:      r.Enabled,
				Type:         mantraev1.RouterType_ROUTER_TYPE_HTTP,
				DnsProviders: DNSProvidersToProto(result),
				CreatedAt:    SafeTimestamp(r.CreatedAt),
				UpdatedAt:    SafeTimestamp(r.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}

			result, _ := q.GetDnsProvidersByTcpRouter(context.Background(), r.ID)
			routersProto = append(routersProto, &mantraev1.Router{
				Id:           r.ID,
				ProfileId:    r.ProfileID,
				AgentId:      SafeString(r.AgentID),
				Name:         r.Name,
				Config:       config,
				Enabled:      r.Enabled,
				Type:         mantraev1.RouterType_ROUTER_TYPE_TCP,
				DnsProviders: DNSProvidersToProto(result),
				CreatedAt:    SafeTimestamp(r.CreatedAt),
				UpdatedAt:    SafeTimestamp(r.UpdatedAt),
			})
		case "udp":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}
			routersProto = append(routersProto, &mantraev1.Router{
				Id:        r.ID,
				ProfileId: r.ProfileID,
				AgentId:   SafeString(r.AgentID),
				Name:      r.Name,
				Config:    config,
				Enabled:   r.Enabled,
				Type:      mantraev1.RouterType_ROUTER_TYPE_UDP,
				CreatedAt: SafeTimestamp(r.CreatedAt),
				UpdatedAt: SafeTimestamp(r.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return routersProto
}

func RoutersByAgentToProto(routers []db.ListRoutersByAgentRow, q *db.Queries) []*mantraev1.Router {
	var routersProto []*mantraev1.Router
	for _, r := range routers {
		switch r.Type {
		case "http":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}

			result, _ := q.GetDnsProvidersByHttpRouter(context.Background(), r.ID)
			routersProto = append(routersProto, &mantraev1.Router{
				Id:           r.ID,
				ProfileId:    r.ProfileID,
				AgentId:      SafeString(r.AgentID),
				Name:         r.Name,
				Config:       config,
				Enabled:      r.Enabled,
				Type:         mantraev1.RouterType_ROUTER_TYPE_HTTP,
				DnsProviders: DNSProvidersToProto(result),
				CreatedAt:    SafeTimestamp(r.CreatedAt),
				UpdatedAt:    SafeTimestamp(r.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}

			result, _ := q.GetDnsProvidersByTcpRouter(context.Background(), r.ID)
			routersProto = append(routersProto, &mantraev1.Router{
				Id:           r.ID,
				ProfileId:    r.ProfileID,
				AgentId:      SafeString(r.AgentID),
				Name:         r.Name,
				Config:       config,
				Enabled:      r.Enabled,
				Type:         mantraev1.RouterType_ROUTER_TYPE_TCP,
				DnsProviders: DNSProvidersToProto(result),
				CreatedAt:    SafeTimestamp(r.CreatedAt),
				UpdatedAt:    SafeTimestamp(r.UpdatedAt),
			})
		case "udp":
			config, err := MarshalStruct(r.Config)
			if err != nil {
				return nil
			}
			routersProto = append(routersProto, &mantraev1.Router{
				Id:        r.ID,
				ProfileId: r.ProfileID,
				AgentId:   SafeString(r.AgentID),
				Name:      r.Name,
				Config:    config,
				Enabled:   r.Enabled,
				Type:      mantraev1.RouterType_ROUTER_TYPE_UDP,
				CreatedAt: SafeTimestamp(r.CreatedAt),
				UpdatedAt: SafeTimestamp(r.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return routersProto
}
