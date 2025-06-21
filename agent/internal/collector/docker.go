package collector

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// GetContainers retrieves all local containers
func GetContainers() ([]types.Container, error) {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New("failed to create Docker client")
	}

	// Get all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list containers")
	}

	var result []types.Container
	// portMap := make(map[int32]int32)

	// Iterate over each container and populate the Container struct
	for _, c := range containers {
		// Skip Traefik
		for _, name := range c.Names {
			if strings.Contains(strings.ToLower(name), "traefik") {
				continue
			}
		}

		// Populate PortInfo
		// for _, p := range c.Ports {
		// 	fmt.Printf("%s:%d -> %d/%s\n", p.IP, p.PublicPort, p.PrivatePort, p.Type)
		// 	// portMap[int32(p.PublicPort)] = int32(p.PrivatePort)
		// }
		result = append(result, c)
	}

	return result, nil
}
