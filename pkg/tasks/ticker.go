package tasks

import (
	"context"
	"log/slog"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/dns"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func StartSync(ctx context.Context) {
	slog.Info("Starting background tasks...")
	go traefikSync(ctx)
	go syncDNS(ctx)
	go sslCheck(ctx)
	go cleanupAgents(ctx)
}

// Refresh forces a refresh of all tasks
func Refresh() {
	traefik.GetTraefikConfig()
	dns.UpdateDNS()

	routers, err := db.Query.ListRouters(context.Background())
	if err != nil {
		slog.Error("Failed to get routers", "error", err)
	}
	for _, router := range routers {
		router.SSLCheck()
	}
}

// traefikSync periodically syncs the Traefik configuration
func traefikSync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(util.App.TraefikInterval))
	defer ticker.Stop()

	traefik.GetTraefikConfig()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			traefik.GetTraefikConfig()
		}
	}
}

// syncDNS periodically syncs the DNS records
func syncDNS(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(util.App.DNSInterval))
	defer ticker.Stop()

	dns.UpdateDNS()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dns.UpdateDNS()
		}
	}
}

func sslCheck(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(util.App.SSLInterval))
	defer ticker.Stop()

	routers, err := db.Query.ListRouters(context.Background())
	if err != nil {
		slog.Error("Failed to get routers", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, router := range routers {
				router.SSLCheck()
			}
		}
	}
}

func cleanupAgents(ctx context.Context) {
	enabled, err := db.Query.GetSettingByKey(context.Background(), "agent-cleanup-enabled")
	if err != nil {
		slog.Error("failed to get agent cleanup timeout", "error", err)
		return
	}

	if enabled.Value != "true" {
		return
	}

	// Timeout to delete old agents
	timeout, err := db.Query.GetSettingByKey(context.Background(), "agent-cleanup-timeout")
	if err != nil {
		slog.Error("failed to get agent cleanup timeout", "error", err)
		return
	}

	timeoutDuration, err := time.ParseDuration(timeout.Value)
	if err != nil {
		slog.Error("failed to parse timeout cleanup duration", "error", err)
	}

	ticker := time.NewTicker(timeoutDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			agents, err := db.Query.ListAgents(context.Background())
			if err != nil {
				slog.Error("failed to query disconnected agents", "error", err)
				continue
			}

			for _, agent := range agents {
				if agent.LastSeen == nil {
					continue
				}

				if now.Sub(*agent.LastSeen) > timeoutDuration {
					if err := db.Query.DeleteAgentByID(context.Background(), agent.ID); err != nil {
						slog.Error(
							"failed to delete disconnected agent",
							"id",
							agent.ID,
							"error",
							err,
						)
					} else {
						slog.Info("Deleted disconnected agent", "id", agent.ID)
						util.Broadcast <- util.EventMessage{
							Type:    "agent_updated",
							Message: "Deleted disconnected agent",
						}
					}
				}
			}
		}
	}
}
