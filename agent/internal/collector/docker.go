// Package collector provides functions for collecting data from the host system.
package collector

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID     string
	Name   string
	Labels map[string]string
	Action events.Action
}

var CleanupActions = []events.Action{
	events.ActionStop,
	events.ActionDestroy,
	events.ActionKill,
	events.ActionRemove,
}

var SyncActions = []events.Action{
	events.ActionStart,
	events.ActionRestart,
	events.ActionUnPause,
}

func dockerClient() (*client.Client, error) {
	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		slog.Warn("Docker socket not found", "path", "/var/run/docker.sock")
	}
	if _, ok := os.LookupEnv("DOCKER_HOST"); !ok {
		_ = os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return nil, err
	}
	return cli, nil
}

// GetContainers retrieves all local containers
func GetContainers() ([]ContainerInfo, error) {
	cli, err := dockerClient()
	if err != nil {
		return nil, errors.New("failed to create Docker client")
	}
	defer func() {
		_ = cli.Close()
	}()

	containers, err := cli.ContainerList(
		context.Background(),
		container.ListOptions{All: false},
	)
	if err != nil {
		return nil, errors.New("failed to list containers")
	}

	var result []ContainerInfo
	for _, c := range containers {
		name := strings.TrimPrefix(c.Names[0], "/")

		// Skip Traefik reverse proxy itself
		if isTraefikProxy(c.Image, c.Labels) {
			slog.Debug("Skipping Traefik proxy", "name", name)
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

		result = append(result, ContainerInfo{
			ID:     c.ID,
			Name:   name,
			Labels: c.Labels,
		})
	}

	return result, nil
}

func WatchContainers(ctx context.Context, ch chan<- ContainerInfo) {
	cli, err := dockerClient()
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return
	}
	defer func() {
		_ = cli.Close()
	}()

	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")
	for _, action := range slices.Concat(SyncActions, CleanupActions) {
		eventFilter.Add("event", string(action))
	}

	msgs, errs := cli.Events(ctx, events.ListOptions{Filters: eventFilter})

	for {
		select {
		case event := <-msgs:
			slog.Debug("Received Docker event", "action", event.Action, "id", event.Actor.ID[:12])
			image := event.Actor.Attributes["image"]
			name := event.Actor.Attributes["name"]

			containerInfo := ContainerInfo{
				ID:     event.Actor.ID,
				Name:   name,
				Action: event.Action,
			}
			inspect, err := cli.ContainerInspect(ctx, event.Actor.ID)
			if err != nil {
				if slices.Contains(CleanupActions, event.Action) {
					ch <- containerInfo
				} else {
					slog.Error("Failed to inspect container", "id", event.Actor.ID, "error", err)
				}
				continue
			}

			// Skip Traefik reverse proxy itself (not whoami or other services)
			if isTraefikProxy(image, containerInfo.Labels) {
				slog.Debug("Skipping Traefik proxy", "name", name)
				continue
			}

			containerInfo.Name = inspect.Name
			containerInfo.Labels = inspect.Config.Labels
			slog.Debug(
				"Sending container info",
				"name",
				containerInfo.Name,
				"action",
				containerInfo.Action,
			)
			ch <- containerInfo

		case err := <-errs:
			slog.Error("Error listening for Docker events", "error", err)
			return
		case <-ctx.Done():
			slog.Info("Stopping container watcher...")
			return
		}
	}
}

// isTraefikProxy checks if a container is the Traefik reverse proxy itself
func isTraefikProxy(image string, labels map[string]string) bool {
	// Check for official Traefik image
	if strings.HasPrefix(image, "traefik:") ||
		strings.HasPrefix(image, "docker.io/traefik:") {
		return true
	}

	// Check for Traefik-specific labels that only the proxy has
	if _, ok := labels["traefik.http.routers.api.rule"]; ok {
		return true
	}
	return false
}
