package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func AgentToProto(a *db.Agent) *mantraev1.Agent {
	var containers []*mantraev1.Container
	if a.Containers != nil {
		for _, c := range *a.Containers {
			containers = append(containers, &mantraev1.Container{
				Id:      c.ID,
				Name:    c.Name,
				Labels:  c.Labels,
				Image:   c.Image,
				Portmap: c.Portmap,
				Status:  c.Status,
				Created: SafeTimestamp(c.Created),
			})
		}
	}

	return &mantraev1.Agent{
		Id:         a.ID,
		ProfileId:  a.ProfileID,
		Hostname:   SafeString(a.Hostname),
		PublicIp:   SafeString(a.PublicIp),
		PrivateIp:  SafeString(a.PrivateIp),
		ActiveIp:   SafeString(a.ActiveIp),
		Token:      a.Token,
		Containers: containers,
		CreatedAt:  SafeTimestamp(a.CreatedAt),
		UpdatedAt:  SafeTimestamp(a.UpdatedAt),
	}
}

func AgentsToProto(agents []db.Agent) []*mantraev1.Agent {
	var agentsProto []*mantraev1.Agent
	for _, a := range agents {
		agentsProto = append(agentsProto, AgentToProto(&a))
	}
	return agentsProto
}
