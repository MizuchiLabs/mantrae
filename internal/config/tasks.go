package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/dns"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
)

func (a *App) setupBackgroundJobs(ctx context.Context) {
	slog.Info("Starting background tasks...")
	go a.traefikSync(ctx)
	go a.syncDNS(ctx)
	// go sslCheck(ctx)
	// go cleanupAgents(ctx)
	// go cleanupRouters(ctx)
}

// traefikSync periodically syncs the Traefik configuration
func (a *App) traefikSync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(a.Config.Background.Traefik))
	defer ticker.Stop()

	traefik.GetTraefikConfig(a.DB)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			traefik.GetTraefikConfig(a.DB)
		}
	}
}

// syncDNS periodically syncs the DNS records
func (a *App) syncDNS(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(a.Config.Background.DNS))
	defer ticker.Stop()

	if err := dns.UpdateDNS(a.DB); err != nil {
		slog.Error("Failed to update DNS", "error", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := dns.UpdateDNS(a.DB); err != nil {
				slog.Error("Failed to update DNS", "error", err)
			}
		}
	}
}

// func sslCheck(ctx context.Context) {
// 	ticker := time.NewTicker(time.Second * time.Duration(util.App.SSLInterval))
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			// Fetch new router list
// 			routers, err := db.Query.ListRouters(context.Background())
// 			if err != nil {
// 				slog.Error("Failed to get routers", "error", err)
// 			}

// 			for _, router := range routers {
// 				router.SSLCheck()
// 			}
// 		}
// 	}
// }

// func cleanupAgents(ctx context.Context) {
// 	enabled, err := db.Query.GetSetting(context.Background(), "agent-cleanup-enabled")
// 	if err != nil {
// 		slog.Error("failed to get agent cleanup timeout", "error", err)
// 		return
// 	}

// 	if enabled.Value != "true" {
// 		return
// 	}

// 	// Timeout to delete old agents
// 	timeout, err := db.Query.GetSetting(context.Background(), "agent-cleanup-timeout")
// 	if err != nil {
// 		slog.Error("failed to get agent cleanup timeout", "error", err)
// 		return
// 	}

// 	timeoutDuration, err := time.ParseDuration(timeout.Value)
// 	if err != nil {
// 		slog.Error("failed to parse timeout cleanup duration", "error", err)
// 	}

// 	ticker := time.NewTicker(timeoutDuration)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			now := time.Now()
// 			agents, err := db.Query.ListAgents(context.Background())
// 			if err != nil {
// 				slog.Error("failed to query disconnected agents", "error", err)
// 				continue
// 			}

// 			for _, agent := range agents {
// 				if agent.LastSeen == nil {
// 					continue
// 				}

// 				if now.Sub(*agent.LastSeen) > timeoutDuration {
// 					if err := db.Query.DeleteAgent(context.Background(), agent.ID); err != nil {
// 						slog.Error(
// 							"failed to delete disconnected agent",
// 							"id",
// 							agent.ID,
// 							"error",
// 							err,
// 						)
// 					} else {
// 						slog.Info("Deleted disconnected agent", "id", agent.ID)
// 						util.Broadcast <- util.EventMessage{
// 							Type:    "agent_updated",
// 							Message: "Deleted disconnected agent",
// 						}
// 					}

// 					// Delete all connected routers
// 					routers, err := db.Query.ListRoutersByAgentID(context.Background(), &agent.ID)
// 					if err != nil {
// 						slog.Error("Failed to get routers", "error", err)
// 						continue
// 					}
// 					for _, router := range routers {
// 						if err := db.Query.DeleteRouter(context.Background(), router.ID); err != nil {
// 							slog.Error("Failed to delete router", "id", router.ID, "error", err)
// 							return
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// cleanupRouters periodically deletes routers from offline/deleted agents
// func (a *App) cleanupRouters(ctx context.Context) {
// 	q := db.New(a.DB)
// 	// Timeout to delete old agents
// 	timeout, err := q.GetSetting(context.Background(), "agent-cleanup-timeout")
// 	if err != nil {
// 		slog.Error("failed to get agent cleanup timeout", "error", err)
// 		return
// 	}
//
// 	timeoutDuration, err := time.ParseDuration(timeout.Value)
// 	if err != nil {
// 		slog.Error("failed to parse timeout cleanup duration", "error", err)
// 	}
//
// 	ticker := time.NewTicker(timeoutDuration)
// 	defer ticker.Stop()
//
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			q := db.New(a.DB)
//
// 			profiles, err := q.ListProfiles(ctx)
// 			if err != nil {
// 				slog.Error("failed to query disconnected agents", "error", err)
// 				continue
// 			}
//
// 			for _, profile := range profiles {
// 				config, err := q.GetTraefikConfigBySource(ctx, db.GetTraefikConfigBySourceParams{
// 					ProfileID: profile.ID,
// 					Source:    source.Agent,
// 				})
// 				if err != nil {
// 					slog.Error("failed to get agent config", "error", err)
// 					continue
// 				}
// 				if config.Config == nil {
// 					continue
// 				}
// 				// TODO
// 			}
// 		}
// 	}
// }
