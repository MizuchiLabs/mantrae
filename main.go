package main

import (
	"context"
	"fmt"
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
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Suggest:               true,
		Name:                  "mantrae",
		Version:               fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date),
		Usage:                 "traefik configuration manager",
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			level := slog.LevelInfo
			if cmd.Bool("debug") {
				level = slog.LevelDebug
			}
			slog.SetDefault(
				slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})),
			)
			return ctx, nil
		},
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
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug logging",
				Sources: cli.EnvVars("MANTRAE_DEBUG"),
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
