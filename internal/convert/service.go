package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func HTTPServiceToProto(s *db.HttpService) *mantraev1.Service {
	config, err := MarshalStruct(s.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    config,
		Enabled:   s.Enabled,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func TCPServiceToProto(s *db.TcpService) *mantraev1.Service {
	config, err := MarshalStruct(s.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    config,
		Enabled:   s.Enabled,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func UDPServiceToProto(s *db.UdpService) *mantraev1.Service {
	config, err := MarshalStruct(s.Config)
	if err != nil {
		return nil
	}

	return &mantraev1.Service{
		Id:        s.ID,
		ProfileId: s.ProfileID,
		AgentId:   SafeString(s.AgentID),
		Name:      s.Name,
		Config:    config,
		Enabled:   s.Enabled,
		Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
		CreatedAt: SafeTimestamp(s.CreatedAt),
		UpdatedAt: SafeTimestamp(s.UpdatedAt),
	}
}

func HTTPServicesToProto(services []db.HttpService) []*mantraev1.Service {
	var servicesProto []*mantraev1.Service
	for _, s := range services {
		servicesProto = append(servicesProto, HTTPServiceToProto(&s))
	}
	return servicesProto
}

func TCPServicesToProto(services []db.TcpService) []*mantraev1.Service {
	var servicesProto []*mantraev1.Service
	for _, s := range services {
		servicesProto = append(servicesProto, TCPServiceToProto(&s))
	}
	return servicesProto
}

func UDPServicesToProto(services []db.UdpService) []*mantraev1.Service {
	var servicesProto []*mantraev1.Service
	for _, s := range services {
		servicesProto = append(servicesProto, UDPServiceToProto(&s))
	}
	return servicesProto
}

// Specialized batch conversion functions
func ServicesByProfileToProto(services []db.ListServicesByProfileRow) []*mantraev1.Service {
	var servicesProto []*mantraev1.Service
	for _, s := range services {
		switch s.Type {
		case "http":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "udp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return servicesProto
}

func ServicesByAgentToProto(services []db.ListServicesByAgentRow) []*mantraev1.Service {
	var servicesProto []*mantraev1.Service
	for _, s := range services {
		switch s.Type {
		case "http":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_HTTP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "tcp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_TCP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		case "udp":
			config, err := MarshalStruct(s.Config)
			if err != nil {
				return nil
			}
			servicesProto = append(servicesProto, &mantraev1.Service{
				Id:        s.ID,
				ProfileId: s.ProfileID,
				AgentId:   SafeString(s.AgentID),
				Name:      s.Name,
				Config:    config,
				Enabled:   s.Enabled,
				Type:      mantraev1.ServiceType_SERVICE_TYPE_UDP,
				CreatedAt: SafeTimestamp(s.CreatedAt),
				UpdatedAt: SafeTimestamp(s.UpdatedAt),
			})
		default:
			return nil
		}
	}
	return servicesProto
}
