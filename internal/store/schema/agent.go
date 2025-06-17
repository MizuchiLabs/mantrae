package schema

import "time"

type AgentPrivateIPs struct {
	IPs []string `json:"privateIps,omitempty"`
}

type AgentContainer struct {
	ID      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Image   string            `json:"image,omitempty"`
	Portmap map[int32]int32   `json:"portmap,omitempty"`
	Status  string            `json:"status,omitempty"`
	Created *time.Time        `json:"created,omitempty"`
}

type AgentContainers []AgentContainer
