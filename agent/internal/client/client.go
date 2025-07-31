// Package client provides the agent's main gRPC client logic.
package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
)

type Agent struct {
	config    *Config
	client    mantraev1connect.AgentServiceClient
	activeIP  string
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
	}
}

func (a *Agent) Run(ctx context.Context) {
	a.healthCheck(ctx)

	healthTicker := time.NewTicker(a.config.HealthCheckInterval)
	defer healthTicker.Stop()

	updateTicker := time.NewTicker(a.config.UpdateInterval)
	defer updateTicker.Stop()

	for {
		select {
		case <-healthTicker.C:
			a.healthCheck(ctx)
		case <-updateTicker.C:
			if a.connected {
				if err := a.syncContainers(ctx); err != nil {
					slog.Error("Failed to sync containers", "error", err)
					a.connected = false
				}
			}
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
		a.activeIP = ip
	}

	slog.Debug("Health check successful", "active_ip", a.activeIP)
}

func (a *Agent) syncContainers(ctx context.Context) error {
	containers, err := collector.GetContainers()
	if err != nil {
		return fmt.Errorf("failed to get containers: %w", err)
	}

	syncer := newResourceSyncer(a.config, a.activeIP)

	for _, container := range containers {
		if err := syncer.processContainer(ctx, container); err != nil {
			return fmt.Errorf("failed to process container %s: %w", container.ID[:12], err)
		}
	}

	return syncer.cleanup(ctx)
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
