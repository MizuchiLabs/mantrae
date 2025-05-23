package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/dns"
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

func (a *App) setupBackgroundJobs(ctx context.Context) {
	slog.Info("Starting background tasks...")
	go a.traefikSync(ctx)
	go a.syncDNS(ctx)
	go a.cleanupAgents(ctx)
}

// traefikSync periodically syncs the Traefik configuration
func (a *App) traefikSync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(a.Config.Background.Traefik))
	defer ticker.Stop()

	traefik.GetTraefikConfig(a.Conn.Get())
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			traefik.GetTraefikConfig(a.Conn.Get())
		}
	}
}

// syncDNS periodically syncs the DNS records
func (a *App) syncDNS(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(a.Config.Background.DNS))
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
	ticker := time.NewTicker(time.Second * time.Duration(a.Config.Background.Agent))
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
