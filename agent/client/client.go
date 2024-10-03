package client

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Client(quit chan os.Signal) {
	token := LoadToken()

	// Decode the JWT to get auth token and server URL
	claims, err := util.DecodeJWT(token)
	if err != nil {
		DeleteToken()
		log.Fatalf("Failed to decode JWT: %v", err)
	}

	// Create a new client instance
	client := agentv1connect.NewAgentServiceClient(
		http.DefaultClient,
		claims.ServerURL,
		connect.WithGRPC(),
	)
	slog.Info("Connected to", "server", claims.ServerURL)

	// Start a goroutine for sending container data
	go func() {
		for {
			// Send machine/container info
			if _, err := client.GetContainer(context.Background(), sendContainer(claims.Secret)); err != nil {
				slog.Error(
					"Failed to send container info",
					"server",
					claims.ServerURL,
					"error",
					err,
				)
			}
			// Wait 10 seconds before sending the next container info
			time.Sleep(10 * time.Second)
		}
	}()

	// Start a separate goroutine for refreshing the token
	go func() {
		for {
			// Refresh token
			tokenRequest := connect.NewRequest(&agentv1.RefreshTokenRequest{Token: token})
			tokenRequest.Header().Set("Authorization", "Bearer "+claims.Secret)
			newToken, err := client.RefreshToken(context.Background(), tokenRequest)
			if err != nil {
				slog.Error("Failed to refresh token", "server", claims.ServerURL, "error", err)
			} else {
				SaveToken(newToken.Msg.Token)
				token = newToken.Msg.Token
			}

			time.Sleep(1 * time.Hour)
		}
	}()

	// Wait for the main loop to finish
	<-quit
}

func LoadToken() string {
	token, err := os.ReadFile("token")
	if err != nil {
		slog.Error("Failed to read token", "error", err)
	}
	return strings.TrimSpace(string(token))
}

func SaveToken(token string) {
	err := os.WriteFile("token", []byte(token), 0644)
	if err != nil {
		slog.Error("Failed to write token", "error", err)
	}
}

func DeleteToken() {
	err := os.Remove("token")
	if err != nil {
		slog.Error("Failed to delete token", "error", err)
	}
}

// sendContainer creates a GetContainerRequest with information about the local machine
func sendContainer(secret string) *connect.Request[agentv1.GetContainerRequest] {
	var request agentv1.GetContainerRequest

	// Get machine ID
	machineID, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		request.Id = "unknown"
	}

	if len(machineID) > 0 {
		request.Id = strings.TrimSpace(string(machineID))
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		request.Hostname = "unknown"
	}
	request.Hostname = hostname

	request.PublicIp, err = util.GetPublicIP()
	if err != nil {
		slog.Error("Failed to get public IP", "error", err)
	}
	request.PrivateIps, err = util.GetPrivateIP()
	if err != nil {
		slog.Error("Failed to get local IP", "error", err)
	}
	request.Containers, err = getContainers()
	if err != nil {
		slog.Error("Failed to get containers", "error", err)
	}
	request.LastSeen = timestamppb.New(time.Now())

	req := connect.NewRequest(&request)
	req.Header().Set("authorization", "Bearer "+secret)
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

	// Iterate over each container and populate the Container struct
	for _, c := range containers {
		// Retrieve the container details
		containerJSON, err := cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			slog.Error("Failed to inspect container", "container", c.ID, "error", err)
			continue
		}

		// Populate PortInfo
		var ports []int32
		for _, portmap := range containerJSON.NetworkSettings.Ports {
			for _, binding := range portmap {
				port, err := strconv.ParseInt(binding.HostPort, 10, 32)
				if err != nil {
					slog.Error("Failed to parse port", "port", port, "error", err)
					continue
				}
				ports = append(ports, int32(port))
			}
		}

		// Remove duplicates
		slices.Sort(ports)
		ports = slices.Compact(ports)

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
			Ports:   ports,
			Status:  containerJSON.State.Status,
			Created: timestamppb.New(created),
		}

		result = append(result, container)
	}

	return result, nil
}
