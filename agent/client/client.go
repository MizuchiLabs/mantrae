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

var tokenFilename = ".mantrae-token"

func Client(quit chan os.Signal) {
	// Bootstrap token: disk > env
	var token string
	tokenFile, err := os.ReadFile(tokenFilename)
	if err != nil {
		token = os.Getenv("TOKEN")
		if token == "" {
			log.Fatal("No token found in environment or .mantrae-token file")
		}
	} else {
		token = string(tokenFile)
	}

	claims, err := DecodeJWT(token)
	if err != nil {
		log.Fatalf("Failed to decode JWT: %v", err)
	}

	// Initialize client
	client := agentv1connect.NewAgentServiceClient(http.DefaultClient, claims.ServerURL)

	// Initial HealthCheck to confirm & maybe rotate
	token, claims = doHealth(client, token, claims, quit)
	if err = os.WriteFile("token.jwt", []byte(token), 0600); err != nil {
		slog.Error("Failed to write token file", "error", err)
	}
	slog.Info("Connected to server", "server", claims.ServerURL)

	// Prepare tickers
	healthTicker := time.NewTicker(15 * time.Second)
	defer healthTicker.Stop()
	containerTicker := time.NewTicker(10 * time.Second)
	defer containerTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			token, claims = doHealth(client, token, claims, quit)
		case <-containerTicker.C:
			doContainer(client, token, claims, quit)
		case <-quit:
			slog.Info("Shutting down agent...")
			return
		}
	}
}

// doHealth invokes HealthCheck, handles token rotation & agent deletion
func doHealth(
	client agentv1connect.AgentServiceClient,
	token string,
	claims *Claims,
	quit chan os.Signal,
) (newToken string, newClaims *Claims) {
	req := connect.NewRequest(&agentv1.HealthCheckRequest{
		AgentId:   claims.AgentID,
		ProfileId: claims.ProfileID,
	})
	req.Header().Set("Authorization", "Bearer "+token)

	resp, err := client.HealthCheck(context.Background(), req)
	if err != nil {
		slog.Warn("HealthCheck error, skipping rotation", "error", err)
		return token, claims
	}
	if !resp.Msg.Ok {
		slog.Warn("Agent deleted by server, shutting down")
		quit <- os.Interrupt
		return token, claims
	}
	if rt := resp.Msg.GetToken(); rt != "" && rt != token {
		slog.Info("Rotating token...")
		if err = os.WriteFile(tokenFilename, []byte(rt), 0600); err != nil {
			slog.Error("Failed to write token file", "error", err)
		}
		newClaims, err := DecodeJWT(rt)
		if err != nil {
			slog.Error("Invalid rotated token", "error", err)
			return token, claims
		}
		return rt, newClaims
	}
	return token, claims
}

// doContainer invokes GetContainer using current token/claims
func doContainer(
	client agentv1connect.AgentServiceClient,
	token string,
	claims *Claims,
	quit chan os.Signal,
) {
	// build payload
	req := sendContainerRequest(token, claims)
	_, err := client.GetContainer(context.Background(), req)

	switch connect.CodeOf(err) {
	case connect.CodeNotFound:
		slog.Warn("Agent deleted by server, shutting down")
		quit <- os.Interrupt
	case connect.CodeInternal:
		slog.Error("GetContainer server error", "error", err)
	case connect.CodeUnauthenticated:
		slog.Warn("Token invalid, will pick up on next health tick")
	default:
		if err != nil {
			slog.Warn("GetContainer error", "error", err)
		}
	}
}

// sendContainer creates a GetContainerRequest with information about the local machine
func sendContainerRequest(
	token string,
	claims *Claims,
) *connect.Request[agentv1.GetContainerRequest] {
	var request agentv1.GetContainerRequest
	request.AgentId = claims.AgentID
	request.ProfileId = claims.ProfileID

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
	request.Updated = timestamppb.New(time.Now())

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
