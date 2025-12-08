package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/agent/internal/client"
	"github.com/mizuchilabs/mantrae/pkg/build"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Suggest:               true,
		Name:                  "mantrae-agent",
		Version:               meta.Version,
		Usage:                 "mantrae-agent [command]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := client.Load(cmd)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Start agent
			client.NewAgent(cfg).Run(ctx)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Check for updates or update mantrae-agent to the latest version",
				Description: `Check if a newer version of mantrae-agent is available.
Use the --install flag to download and install the latest version.

Note: Automatic installation does not work inside Docker containers.`,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					build.Update(cmd.Bool("install"))
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
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug logging",
				Sources: cli.EnvVars("DEBUG"),
			},
			&cli.BoolFlag{
				Name:  "install",
				Usage: "Download and install the latest version (used with update command, does not work in Docker)",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "Mantrae API token",
				Value:   "",
				Sources: cli.EnvVars("TOKEN"),
			},
			&cli.StringFlag{
				Name:    "host",
				Usage:   "Mantrae API host",
				Value:   "",
				Sources: cli.EnvVars("HOST"),
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
