package client

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Client(quit chan os.Signal) {
	token := os.Getenv("TOKEN")

	// Decode the JWT to get auth token and server URL
	claims, err := DecodeJWT(token)
	if err != nil {
		log.Fatalf("Failed to decode JWT: %v", err)
	}

	// Create a new client instance
	client := agentv1connect.NewAgentServiceClient(
		http.DefaultClient,
		claims.ServerURL+"/grpc",
		connect.WithGRPC(),
	)

	// Test connection
	healthCheckRequest := connect.NewRequest(&agentv1.HealthCheckRequest{Id: getID()})
	healthCheckRequest.Header().Set("Authorization", "Bearer "+token)
	result, err := client.HealthCheck(context.Background(), healthCheckRequest)
	if err != nil {
		slog.Error("Failed to connect to server", "server", claims.ServerURL, "error", err)
		os.Exit(1)
	}
	if !result.Msg.Ok {
		slog.Warn("Agent was deleted, cleaning up...", "server", claims.ServerURL)
		quit <- os.Interrupt
	} else {
		slog.Info("Connected to", "server", claims.ServerURL)
	}

	tickerContainer := time.NewTicker(10 * time.Second)
	defer tickerContainer.Stop()

	// Start a goroutine for sending container data
	go func() {
		connected := true
		for range tickerContainer.C {
			// Send machine/container info
			_, err := client.GetContainer(context.Background(), sendContainer(token))

			if err != nil && connected {
				slog.Warn("Lost connection to server, retrying...", "server", claims.ServerURL)
				connected = false
			} else if err == nil && !connected {
				slog.Info("Reconnected to server", "server", claims.ServerURL)
				connected = true
			}
		}
	}()

	healthTicker := time.NewTicker(5 * time.Second)
	defer healthTicker.Stop()

	// Start a goroutine for sending health data
	go func() {
		for range healthTicker.C {
			healthCheckRequest := connect.NewRequest(&agentv1.HealthCheckRequest{Id: getID()})
			healthCheckRequest.Header().Set("Authorization", "Bearer "+token)
			result, err := client.HealthCheck(context.Background(), healthCheckRequest)
			if err != nil {
				continue
			}
			if !result.Msg.Ok {
				slog.Warn("Agent was deleted, cleaning up...", "server", claims.ServerURL)
				quit <- os.Interrupt
			}
		}
	}()

	// Wait for the main loop to finish
	<-quit
}

// sendContainer creates a GetContainerRequest with information about the local machine
func sendContainer(token string) *connect.Request[agentv1.GetContainerRequest] {
	var request agentv1.GetContainerRequest

	// Get machine ID
	request.Id = getID()

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		request.Hostname = "unknown"
	}
	request.Hostname = hostname
	request.PublicIp, err = GetPublicIP()
	if err != nil {
		slog.Error("Failed to get public IP", "error", err)
	}
	request.PrivateIps, err = GetPrivateIP()
	if err != nil {
		slog.Error("Failed to get local IP", "error", err)
	}
	request.Containers, err = getContainers()
	if err != nil {
		slog.Error("Failed to get containers", "error", err)
	}
	request.LastSeen = timestamppb.New(time.Now())

	req := connect.NewRequest(&request)
	req.Header().Set("authorization", "Bearer "+token)
	return req
}

// getContainers retrieves all containers and their info on the local machine
func getContainers() ([]*agentv1.Container, error) {
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

	var result []*agentv1.Container
	portMap := make(map[int32]int32)

	// Iterate over each container and populate the Container struct
	for _, c := range containers {
		// Retrieve the container details
		containerJSON, err := cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			slog.Error("Failed to inspect container", "container", c.ID, "error", err)
			continue
		}

		// Skip Traefik
		skipTraefik := os.Getenv("SKIP_TRAEFIK")
		if skipTraefik == "true" {
			if strings.Contains(strings.ToLower(containerJSON.Config.Image), "traefik") ||
				(len(c.Names) > 0 && strings.Contains(strings.ToLower(c.Names[0]), "traefik")) {
				continue
			}
		}

		// Populate PortInfo
		for port, bindings := range containerJSON.NetworkSettings.Ports {
			for _, binding := range bindings {
				// Get external port
				externalPort, err := strconv.ParseInt(binding.HostPort, 10, 32)
				if err != nil {
					slog.Error(
						"Failed to parse external port",
						"port",
						binding.HostPort,
						"error",
						err,
					)
					continue
				}

				// Get internal port from the port key
				internalPort, err := strconv.ParseInt(
					port.Port(),
					10,
					32,
				) // port is of type nat.Port
				if err != nil {
					slog.Error("Failed to parse internal port", "port", port.Port(), "error", err)
					continue
				}

				// Map internal port to external port
				portMap[int32(internalPort)] = int32(externalPort)
			}
		}

		created, err := time.Parse(time.RFC3339, containerJSON.Created)
		if err != nil {
			slog.Error("Failed to parse created time", "time", containerJSON.Created, "error", err)
		}

		// Populate the Container struct
		container := &agentv1.Container{
			Id:      c.ID,
			Name:    c.Names[0], // Take the first name if multiple exist
			Labels:  containerJSON.Config.Labels,
			Image:   containerJSON.Config.Image,
			Portmap: portMap,
			Status:  containerJSON.State.Status,
			Created: timestamppb.New(created),
		}

		result = append(result, container)
	}

	return result, nil
}

// Get machine ID
func getID() string {
	var id string
	machineID, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		id = "unknown"
	}

	if len(machineID) > 0 {
		id = strings.TrimSpace(string(machineID))
	}
	return id
}
