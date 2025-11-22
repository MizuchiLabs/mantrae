// package tasks provides functionality for running periodic tasks.
package tasks

import (
	"context"
	"log/slog"
	"time"

	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/dns"
	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
)

type Scheduler struct {
	ctx context.Context
	cfg *config.App
}

func NewScheduler(ctx context.Context, cfg *config.App) *Scheduler {
	return &Scheduler{ctx: ctx, cfg: cfg}
}

func (s *Scheduler) Start() {
	go s.syncDNS()
	go s.cleanupAgents()
}

// syncDNS periodically syncs the DNS records
func (s *Scheduler) syncDNS() {
	duration, ok := s.cfg.SM.Get(s.ctx, settings.KeyDNSSyncInterval)
	if !ok {
		slog.Error("Failed to get DNS sync interval setting")
		return
	}

	ticker := time.NewTicker(settings.AsDuration(duration))
	defer ticker.Stop()

	manager := dns.NewManager(s.cfg.Conn, s.cfg.Secret)
	if err := manager.UpdateDNS(s.ctx); err != nil {
		slog.Error("Failed to update DNS", "error", err)
	}
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := manager.UpdateDNS(s.ctx); err != nil {
				slog.Error("Failed to update DNS", "error", err)
			}
		}
	}
}

func (s *Scheduler) cleanupAgents() {
	duration, ok := s.cfg.SM.Get(s.ctx, settings.KeyAgentCleanupInterval)
	if !ok {
		slog.Error("failed to get agent cleanup interval setting")
		return
	}

	ticker := time.NewTicker(settings.AsDuration(duration))
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			enabled, ok := s.cfg.SM.Get(s.ctx, settings.KeyAgentCleanupEnabled)
			if !ok {
				slog.Error("failed to get agent cleanup enabled setting")
				return
			}
			if enabled != "true" {
				slog.Info("agent cleanup is disabled, skipping")
				return
			}

			// Timeout to delete old agents
			timeout, ok := s.cfg.SM.Get(s.ctx, settings.KeyAgentCleanupInterval)
			if !ok {
				slog.Error("failed to get agent cleanup interval setting")
				return
			}

			// List profiles
			profiles, err := s.cfg.Conn.GetQuery().ListProfiles(s.ctx, db.ListProfilesParams{})
			if err != nil {
				slog.Error("failed to list profiles", "error", err)
				continue
			}

			var agents []db.Agent
			for _, profile := range profiles {
				a, err := s.cfg.Conn.GetQuery().
					ListAgents(s.ctx, db.ListAgentsParams{ProfileID: profile.ID})
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
					if err := s.cfg.Conn.GetQuery().DeleteAgent(s.ctx, agent.ID); err != nil {
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
