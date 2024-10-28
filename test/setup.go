package test

import (
	"context"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/dns"
	"github.com/brianvoe/gofakeit/v7"
)

func SetupDB() error {
	// Initialize the in-memory database
	if err := db.InitDB(); err != nil {
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
		_, err := db.Query.CreateUser(
			context.Background(),
			db.CreateUserParams{
				Username: gofakeit.Username(),
				Password: gofakeit.Password(true, true, true, true, true, 24),
				Email:    &email,
				Type:     "user",
			},
		)
		if err != nil {
			return err
		}
	}

	// Seed Providers
	for i := 0; i < 5; i++ {
		apiUrl := gofakeit.URL()
		zonetype := dns.ZoneTypes[gofakeit.Number(0, len(dns.ZoneTypes)-1)]

		_, err := db.Query.CreateProvider(context.Background(), db.CreateProviderParams{
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
	}

	return nil
}
