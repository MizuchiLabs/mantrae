package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/mizuchilabs/mantrae/internal/dns"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/mizuchilabs/mantrae/internal/util"
)

// setupBackgroundJobs initiates essential background operations for the application.
func (a *App) setupBackgroundJobs(ctx context.Context) {
	slog.Info("Starting background tasks...")
	go a.syncTraefik(ctx)
	go a.syncDNS(ctx)
	go a.cleanupAgents(ctx)
}

// syncTraefik periodically syncs the Traefik configuration
func (a *App) syncTraefik(ctx context.Context) {
	duration, err := a.SM.Get(ctx, settings.KeyTraefikSyncInterval)
	if err != nil {
		slog.Error("Failed to get Traefik sync interval setting", "error", err)
		return
	}
	ticker := time.NewTicker(duration.Duration(time.Hour))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			instances, err := a.Conn.GetQuery().ListTraefikInstances(ctx)
			if err != nil {
				slog.Error("failed to list traefik instances", "error", err)
				continue
			}

			for _, instance := range instances {
				if err := traefik.UpdateTraefikAPI(a.Conn.Get(), instance.ID); err != nil {
					slog.Error("failed to update traefik api", "error", err)
					continue
				}
			}
		}
	}
}

// syncDNS periodically syncs the DNS records
func (a *App) syncDNS(ctx context.Context) {
	duration, err := a.SM.Get(ctx, settings.KeyDNSSyncInterval)
	if err != nil {
		slog.Error("Failed to get DNS sync interval setting", "error", err)
		return
	}
	ticker := time.NewTicker(duration.Duration(time.Hour))
	defer ticker.Stop()

	if err := dns.UpdateDNS(a.Conn.Get()); err != nil {
		slog.Error("Failed to update DNS", "error", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := dns.UpdateDNS(a.Conn.Get()); err != nil {
				slog.Error("Failed to update DNS", "error", err)
			}
		}
	}
}

func (a *App) cleanupAgents(ctx context.Context) {
	duration, err := a.SM.Get(ctx, settings.KeyAgentCleanupInterval)
	if err != nil {
		slog.Error("failed to get agent cleanup interval setting", "error", err)
		return
	}
	ticker := time.NewTicker(duration.Duration(time.Hour))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			enabled, err := a.SM.Get(ctx, settings.KeyAgentCleanupEnabled)
			if err != nil {
				slog.Error("failed to get agent cleanup enabled setting", "error", err)
				return
			}
			if enabled.Bool(false) {
				return
			}

			// Timeout to delete old agents
			timeout, err := a.SM.Get(ctx, settings.KeyAgentCleanupInterval)
			if err != nil {
				slog.Error("failed to get agent cleanup interval setting", "error", err)
				return
			}

			q := a.Conn.GetQuery()
			agents, err := q.ListAgents(ctx)
			if err != nil {
				slog.Error("failed to list agents", "error", err)
				return
			}

			for _, agent := range agents {
				if agent.UpdatedAt == nil || agent.Hostname == nil {
					continue
				}

				if time.Now().Sub(*agent.UpdatedAt) > timeout.Duration(time.Hour*24) {
					if err := q.DeleteTraefikConfigByAgent(ctx, &agent.ID); err != nil {
						slog.Error(
							"failed to delete agent config",
							"id",
							agent.ID,
							"error",
							err,
						)
						continue
					}
					if err := q.DeleteAgent(ctx, agent.ID); err != nil {
						slog.Error(
							"failed to delete disconnected agent",
							"id",
							agent.ID,
							"error",
							err,
						)
						continue
					} else {
						slog.Info("Deleted disconnected agent", "id", agent.ID)
						util.Broadcast <- util.EventMessage{
							Type:     util.EventTypeDelete,
							Category: util.EventCategoryAgent,
						}
					}
				}
			}
		}
	}
}
