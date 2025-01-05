package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/dns"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func SetupDB() error {
	os.Setenv("SECRET", gofakeit.Password(true, true, true, true, true, 24))

	// Initialize the in-memory database
	if err := db.InitDB(); err != nil {
		return err
	}

	if err := setDefaultSettings(); err != nil {
		return err
	}

	// Seed mock data
	if err := seedMockData(); err != nil {
		return err
	}

	return nil
}

// Seed random data for testing tables
func seedMockData() error {
	// Seed Users
	for i := 0; i < 5; i++ {
		email := gofakeit.Email()
		password, _ := util.HashPassword(gofakeit.Password(true, true, true, true, true, 24))
		_, err := db.Query.UpsertUser(
			context.Background(),
			db.UpsertUserParams{
				Username: gofakeit.Username(),
				Password: password,
				Email:    &email,
				IsAdmin:  gofakeit.Bool(),
			},
		)
		if err != nil {
			return err
		}
	}

	// Additional testuser
	password, err := util.HashPassword("test")
	if err != nil {
		return err
	}
	_, err = db.Query.UpsertUser(
		context.Background(),
		db.UpsertUserParams{
			Username: "test",
			Password: password,
			Email:    nil,
			IsAdmin:  true,
		},
	)
	if err != nil {
		return err
	}

	// Seed Providers
	var dnsProvider []db.Provider
	for i := 0; i < 5; i++ {
		apiUrl := gofakeit.URL()
		zonetype := dns.ZoneTypes[gofakeit.Number(0, len(dns.ZoneTypes)-1)]

		data, err := db.Query.CreateProvider(context.Background(), db.CreateProviderParams{
			Name:       gofakeit.Name(),
			Type:       dns.DNSProviders[gofakeit.Number(0, len(dns.DNSProviders)-1)],
			ExternalIp: gofakeit.IPv4Address(),
			ApiKey:     gofakeit.UUID(),
			ApiUrl:     &apiUrl,
			ZoneType:   &zonetype,
			Proxied:    gofakeit.Bool(),
			IsActive:   gofakeit.Bool(),
		})
		if err != nil {
			return err
		}

		dnsProvider = append(dnsProvider, data)
	}

	// Seed Profiles
	var profiles []db.Profile
	for i := 0; i < 5; i++ {
		username := gofakeit.Username()
		password := gofakeit.Password(true, true, true, true, true, 24)
		data, err := db.Query.CreateProfile(context.Background(), db.CreateProfileParams{
			Name:     gofakeit.Name(),
			Url:      gofakeit.URL(),
			Username: &username,
			Password: &password,
		})
		if err != nil {
			return err
		}

		profiles = append(profiles, data)
	}

	// Seed EntryPoints
	var entrypoints []db.Entrypoint
	for i := 0; i < 5; i++ {
		names := []string{"web", "websecure", "traefik"}
		addresses := []string{"80", "443", ":8080"}
		data, err := db.Query.UpsertEntryPoint(context.Background(), db.UpsertEntryPointParams{
			ProfileID: profiles[gofakeit.Number(0, len(profiles)-1)].ID,
			Name:      names[gofakeit.Number(0, len(names)-1)],
			Address:   addresses[gofakeit.Number(0, len(addresses)-1)],
		})
		if err != nil {
			return err
		}

		entrypoints = append(entrypoints, data)
	}

	// Seed Agents
	var agents []db.Agent
	for i := 0; i < 5; i++ {
		ip := gofakeit.IPv4Address()
		privateIps, _ := json.Marshal([]string{gofakeit.IPv4Address(), gofakeit.IPv4Address()})
		data, err := db.Query.UpsertAgent(context.Background(), db.UpsertAgentParams{
			ID:         uuid.New().String(),
			ProfileID:  profiles[gofakeit.Number(0, len(profiles)-1)].ID,
			Hostname:   gofakeit.DomainName(),
			PublicIp:   &ip,
			PrivateIps: privateIps,
			Containers: nil,
			ActiveIp:   &ip,
			Token:      gofakeit.UUID(),
		})
		if err != nil {
			return err
		}

		agents = append(agents, data)
	}

	// Additional testagent
	testAgentIP := gofakeit.IPv4Address()
	_, err = db.Query.UpsertAgent(
		context.Background(),
		db.UpsertAgentParams{
			ID:         "test",
			ProfileID:  profiles[gofakeit.Number(0, len(profiles)-1)].ID,
			Hostname:   "test",
			PublicIp:   &testAgentIP,
			PrivateIps: nil,
			Containers: nil,
			ActiveIp:   nil,
			Token:      gofakeit.UUID(),
		},
	)
	if err != nil {
		return err
	}

	// Seed Routers, Services, Middlewares
	for i := 0; i < 20; i++ {
		name := gofakeit.DomainName()
		protocols := []string{"http", "tcp", "udp"}
		providers := []string{"http", "docker", "kubernetes"}
		priorities := int64(gofakeit.Number(0, 100))
		router, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams{
			ID:          uuid.New().String(),
			ProfileID:   profiles[gofakeit.Number(0, len(profiles)-1)].ID,
			Name:        name,
			Provider:    providers[gofakeit.Number(0, len(providers)-1)],
			Protocol:    protocols[gofakeit.Number(0, len(protocols)-1)],
			AgentID:     &agents[gofakeit.Number(0, len(agents)-1)].ID,
			EntryPoints: entrypoints[gofakeit.Number(0, len(entrypoints)-1)].Name,
			Rule:        "",
			Service:     name,
			Priority:    &priorities,
			DnsProvider: &dnsProvider[gofakeit.Number(0, len(dnsProvider)-1)].ID,
		})
		if err != nil {
			return err
		}

		_, err = db.Query.UpsertService(context.Background(), db.UpsertServiceParams{
			ID:           uuid.New().String(),
			ProfileID:    router.ProfileID,
			Name:         router.Name,
			Provider:     router.Provider,
			Protocol:     router.Protocol,
			AgentID:      router.AgentID,
			LoadBalancer: nil,
		})
		if err != nil {
			return err
		}

		mwProtocols := []string{"http", "tcp"}
		mwTypes := []string{
			"addprefix",
			"stripprefix",
			"stripprefixregex",
			"replacepath",
			"replacepathregex",
			"chain",
			"ipallowlist",
			"basicauth",
			"digestauth",
			"inflightreq",
		}
		_, err = db.Query.UpsertMiddleware(context.Background(), db.UpsertMiddlewareParams{
			ID:        uuid.New().String(),
			ProfileID: profiles[gofakeit.Number(0, len(profiles)-1)].ID,
			Name:      router.Name,
			Provider:  providers[gofakeit.Number(0, len(providers)-1)],
			Type:      mwTypes[gofakeit.Number(0, len(mwTypes)-1)],
			Protocol:  mwProtocols[gofakeit.Number(0, len(mwProtocols)-1)],
			Content:   nil,
		})
		if err != nil {
			return err
		}

	}

	return nil
}

func setDefaultSettings() error {
	baseSettings := []db.Setting{
		{
			Key:   "server-url",
			Value: util.App.ServerURL,
		},
		{
			Key:   "backup-enabled",
			Value: "true",
		},
		{
			Key:   "backup-schedule",
			Value: "0 2 * * 1", // Weekly at 02:00 AM on Monday
		},
		{
			Key:   "backup-keep",
			Value: "3", // Keep 3 backups
		},
		{
			Key:   "agent-cleanup-enabled",
			Value: "true",
		},
		{
			Key:   "agent-cleanup-timeout",
			Value: "168h",
		},
		{
			Key:   "email-host",
			Value: util.App.EmailHost,
		},
		{
			Key:   "email-port",
			Value: util.App.EmailPort,
		},
		{
			Key:   "email-username",
			Value: util.App.EmailUsername,
		},
		{
			Key:   "email-password",
			Value: util.App.EmailPassword,
		},
		{
			Key:   "email-from",
			Value: util.App.EmailFrom,
		},
	}

	for _, setting := range baseSettings {
		if _, err := db.Query.GetSettingByKey(context.Background(), setting.Key); err != nil {
			if _, err := db.Query.CreateSetting(context.Background(), db.CreateSettingParams{
				Key:   setting.Key,
				Value: setting.Value,
			}); err != nil {
				return fmt.Errorf("failed to create setting: %w", err)
			}
		}
	}
	return nil
}
