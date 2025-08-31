// Package collector provides functions for collecting data from the host system.
package collector

import (
	"context"
	"errors"
	"log/slog"
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
	events.ActionDestroy,
	events.ActionKill,
	events.ActionRemove,
}

var SyncActions = []events.Action{
	events.ActionStart,
	events.ActionRestart,
	events.ActionUnPause,
}

// GetContainers retrieves all local containers
func GetContainers() ([]ContainerInfo, error) {
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

	var result []ContainerInfo
	for _, c := range containers {
		if strings.Contains(c.Image, "traefik") {
			slog.Debug("Skipping Traefik container", "name", c.Names[0])
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
			Name:   c.Names[0],
			Labels: c.Labels,
		})
	}

	return result, nil
}

func WatchContainers(ctx context.Context, ch chan<- ContainerInfo) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create Docker client", "error", err)
		return
	}

	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")
	eventFilter.Add("event", string(events.ActionStart))
	eventFilter.Add("event", string(events.ActionRestart))
	eventFilter.Add("event", string(events.ActionDie))
	eventFilter.Add("event", string(events.ActionDestroy))
	eventFilter.Add("event", string(events.ActionKill))
	eventFilter.Add("event", string(events.ActionPause))
	eventFilter.Add("event", string(events.ActionUnPause))
	eventFilter.Add("event", string(events.ActionStop))
	eventFilter.Add("event", string(events.ActionRemove))

	msgs, errs := cli.Events(ctx, events.ListOptions{Filters: eventFilter})

	for {
		select {
		case event := <-msgs:
			image := event.Actor.Attributes["image"]
			name := event.Actor.Attributes["name"]

			// Skip traefik containers
			if strings.Contains(image, "traefik") {
				continue
			}

			slog.Debug("Container event",
				"action", event.Action,
				"id", event.Actor.ID[:12],
				"image", image,
				"name", name,
			)
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

			containerInfo.Name = inspect.Name
			containerInfo.Labels = inspect.Config.Labels
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
