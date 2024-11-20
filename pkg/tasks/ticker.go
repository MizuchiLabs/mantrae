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
		if err := router.DecodeFields(); err != nil {
			continue
		}
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
				if err := router.DecodeFields(); err != nil {
					continue
				}
				router.SSLCheck()
			}
		}
	}
}
