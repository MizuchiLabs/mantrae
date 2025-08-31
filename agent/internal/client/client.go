// Package client provides the agent's main gRPC client logic.
package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
)

type Agent struct {
	config    *Config
	job       *SyncJob
	client    mantraev1connect.AgentServiceClient
	connected bool
}

func NewAgent(cfg *Config) *Agent {
	httpClient := &http.Client{Timeout: cfg.ConnectionTimeout}

	client := mantraev1connect.NewAgentServiceClient(
		httpClient,
		cfg.ServerURL,
		connect.WithInterceptors(authInterceptor(cfg)),
	)

	return &Agent{
		config: cfg,
		client: client,
		job:    NewJob(cfg),
	}
}

func (a *Agent) Run(ctx context.Context) {
	// Run initial health check
	a.healthCheck(ctx)

	// Sync existing containers on startup
	if a.connected {
		if err := a.initialSync(ctx); err != nil {
			slog.Error("Failed to sync containers", "error", err)
		}
	}

	healthTicker := time.NewTicker(a.config.HealthCheckInterval)
	defer healthTicker.Stop()

	// Start container watcher
	containerEvents := make(chan collector.ContainerInfo, 100)
	go collector.WatchContainers(ctx, containerEvents)

	for {
		select {
		case <-healthTicker.C:
			a.healthCheck(ctx)
		case container := <-containerEvents:
			a.handleEvents(ctx, container)
		case <-ctx.Done():
			slog.Info("Agent stopping...")
			return
		}
	}
}

func (a *Agent) healthCheck(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, a.config.HealthTimeout)
	defer cancel()

	info := collector.GetMachineInfo()
	req := connect.NewRequest(&mantraev1.HealthCheckRequest{
		Hostname:  info.Hostname,
		PrivateIp: info.PrivateIP,
		PublicIp:  info.PublicIPs.IPv4,
	})

	resp, err := a.client.HealthCheck(ctx, req)
	if err != nil {
		slog.Error("Health check failed", "error", err)
		a.connected = false
		return
	}

	if resp.Msg.Agent == nil {
		slog.Error("Agent not found on server")
		a.connected = false
		return
	}

	a.connected = true
	if ip := resp.Msg.Agent.ActiveIp; ip != "" {
		a.config.ActiveIP = ip
	}

	slog.Debug("Health check successful", "ip", a.config.ActiveIP)
}

func (a *Agent) initialSync(ctx context.Context) error {
	containers, err := collector.GetContainers()
	if err != nil {
		return fmt.Errorf("failed to get containers: %w", err)
	}

	for _, container := range containers {
		if err := a.job.processContainer(ctx, container); err != nil {
			return fmt.Errorf("failed to process container %s: %w", container.ID[:12], err)
		}
	}

	return a.job.cleanup(ctx)
}

func (a *Agent) handleEvents(ctx context.Context, container collector.ContainerInfo) {
	if container.Action == "" || !a.connected {
		return
	}
	if slices.Contains(collector.SyncActions, container.Action) {
		if err := a.job.processContainer(ctx, container); err != nil {
			slog.Error(
				"Failed to sync container on start",
				"container",
				container.ID[:12],
				"error",
				err,
			)
		}
	}
	if slices.Contains(collector.CleanupActions, container.Action) {
		a.job.removeContainer(container.ID)
		if err := a.job.cleanup(ctx); err != nil {
			slog.Error(
				"Failed to cleanup container on stop/destroy",
				"container",
				container.ID[:12],
				"error",
				err,
			)
		}
	}
}

func authInterceptor(cfg *Config) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("Authorization", "Bearer "+cfg.Token)
			req.Header().Set(meta.HeaderAgentID, cfg.AgentID)
			return next(ctx, req)
		}
	}
}
