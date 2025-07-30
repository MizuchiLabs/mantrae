// Package collector provides functions for collecting data from the host system.
package collector

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// GetContainers retrieves all local containers
func GetContainers() ([]container.Summary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New("failed to create Docker client")
	}

	containers, err := cli.ContainerList(
		context.Background(),
		container.ListOptions{All: false},
	)
	if err != nil {
		return nil, errors.New("failed to list containers")
	}

	var result []container.Summary
	for _, c := range containers {
		// Skip Traefik by name
		skip := false
		for _, name := range c.Names {
			if strings.Contains(strings.ToLower(name), "traefik") {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Must have at least one exposed port
		if len(c.Ports) == 0 {
			continue
		}

		// Must have at least one label
		if len(c.Labels) == 0 {
			continue
		}
		result = append(result, c)
	}

	return result, nil
}
