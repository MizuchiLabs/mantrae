package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/internal/api/server"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/tasks"
	"github.com/urfave/cli/v3"
)

var (
	Version = "debug"
	Commit  string
	Date    string
	Dirty   string
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Suggest:               true,
		Name:                  "mantrae",
		Version:               Version,
		Usage:                 "mantrae [command]",
		Description: `Mantrae simplifies the management of Traefik reverse proxy configurations through an intuitive web interface. Manage routers, middleware, services, and DNS providers with ease.

See https://github.com/mizuchilabs/mantrae for more information.`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			app, err := config.New(ctx, cmd)
			if err != nil {
				slog.Error("Setup failed", "error", err)
				return err
			}

			// Start background jobs
			tasks.NewScheduler(ctx, app).Start()
			return server.NewServer(app).Start(ctx)
		},
		Commands: []*cli.Command{
			{
				Name:  "reset",
				Usage: "Reset a user's password",
				Description: `Reset the password for a specific user account.
By default, resets the admin user's password. Use the --user flag
to specify a different username.`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					app, err := config.New(ctx, cmd)
					if err != nil {
						slog.Error("Setup failed", "error", err)
						return err
					}

					app.ResetPassword(ctx, cmd)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Display version information and exit",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "New password to set for the user (used with reset command)",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Username for password reset operations",
				Value:   "admin",
			},
		},
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
