package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/mizuchilabs/mantrae/internal/dns"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

// setupBackgroundJobs initiates essential background operations for the application.
func (a *App) setupBackgroundJobs(ctx context.Context) {
	go a.syncDNS(ctx)
	go a.cleanupAgents(ctx)
}

// syncDNS periodically syncs the DNS records
func (a *App) syncDNS(ctx context.Context) {
	duration, ok := a.SM.Get(ctx, settings.KeyDNSSyncInterval)
	if !ok {
		slog.Error("Failed to get DNS sync interval setting")
		return
	}

	ticker := time.NewTicker(settings.AsDuration(duration))
	defer ticker.Stop()

	if err := dns.UpdateDNS(ctx, a.Conn.GetQuery()); err != nil {
		slog.Error("Failed to update DNS", "error", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := dns.UpdateDNS(ctx, a.Conn.GetQuery()); err != nil {
				slog.Error("Failed to update DNS", "error", err)
			}
		}
	}
}

func (a *App) cleanupAgents(ctx context.Context) {
	duration, ok := a.SM.Get(ctx, settings.KeyAgentCleanupInterval)
	if !ok {
		slog.Error("failed to get agent cleanup interval setting")
		return
	}

	ticker := time.NewTicker(settings.AsDuration(duration))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			enabled, ok := a.SM.Get(ctx, settings.KeyAgentCleanupEnabled)
			if !ok {
				slog.Error("failed to get agent cleanup enabled setting")
				return
			}
			if enabled != "true" {
				slog.Info("agent cleanup is disabled, skipping")
				return
			}

			// Timeout to delete old agents
			timeout, ok := a.SM.Get(ctx, settings.KeyAgentCleanupInterval)
			if !ok {
				slog.Error("failed to get agent cleanup interval setting")
				return
			}

			// List profiles
			profiles, err := a.Conn.GetQuery().ListProfiles(ctx, db.ListProfilesParams{
				Limit:  -1,
				Offset: 0,
			})
			if err != nil {
				slog.Error("failed to list profiles", "error", err)
				continue
			}

			var agents []db.Agent
			for _, profile := range profiles {
				a, err := a.Conn.GetQuery().ListAgents(ctx, db.ListAgentsParams{
					ProfileID: profile.ID,
					Limit:     -1,
					Offset:    0,
				})
				if err != nil {
					slog.Error("failed to list agents", "error", err)
					continue
				}
				agents = append(agents, a...)
			}

			for _, agent := range agents {
				if agent.UpdatedAt == nil || agent.Hostname == nil {
					continue
				}

				if time.Since(*agent.UpdatedAt) > settings.AsDuration(timeout) {
					if err := a.Conn.GetQuery().DeleteAgent(ctx, agent.ID); err != nil {
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
					}
				}
			}
		}
	}
}
