package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func HTTPServersTransportToProto(in *db.HttpServersTransport) *mantraev1.ServersTransport {
	config, err := MarshalStruct(in.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.ServersTransport{
		Id:        in.ID,
		ProfileId: in.ProfileID,
		AgentId:   SafeString(in.AgentID),
		Name:      in.Name,
		Config:    config,
		Enabled:   in.Enabled,
		Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP,
		CreatedAt: SafeTimestamp(in.CreatedAt),
		UpdatedAt: SafeTimestamp(in.UpdatedAt),
	}
}

func TCPServersTransportToProto(in *db.TcpServersTransport) *mantraev1.ServersTransport {
	config, err := MarshalStruct(in.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.ServersTransport{
		Id:        in.ID,
		ProfileId: in.ProfileID,
		AgentId:   SafeString(in.AgentID),
		Name:      in.Name,
		Config:    config,
		Enabled:   in.Enabled,
		Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP,
		CreatedAt: SafeTimestamp(in.CreatedAt),
		UpdatedAt: SafeTimestamp(in.UpdatedAt),
	}
}

func HTTPServersTransportsToProto(in []db.HttpServersTransport) []*mantraev1.ServersTransport {
	var out []*mantraev1.ServersTransport
	for _, s := range in {
		out = append(out, HTTPServersTransportToProto(&s))
	}
	return out
}

func TCPServersTransportsToProto(in []db.TcpServersTransport) []*mantraev1.ServersTransport {
	var out []*mantraev1.ServersTransport
	for _, s := range in {
		out = append(out, TCPServersTransportToProto(&s))
	}
	return out
}

func ServersTransportsByProfileToProto(
	in []db.ListServersTransportsByProfileRow,
	q *db.Queries,
) []*mantraev1.ServersTransport {
	var out []*mantraev1.ServersTransport
	for _, s := range in {
		switch s.Type {
		case "http":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}

			out = append(out, &mantraev1.ServersTransport{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}

			out = append(out, &mantraev1.ServersTransport{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return out
}

func ServersTransportsByAgentToProto(
	in []db.ListServersTransportsByAgentRow,
	q *db.Queries,
) []*mantraev1.ServersTransport {
	var out []*mantraev1.ServersTransport
	for _, s := range in {
		switch s.Type {
		case "http":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}

			out = append(out, &mantraev1.ServersTransport{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}

			out = append(out, &mantraev1.ServersTransport{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return out
}
