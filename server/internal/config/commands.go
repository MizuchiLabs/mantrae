package config

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/urfave/cli/v3"
)

func (a *App) ResetPassword(ctx context.Context, cmd *cli.Command) {
	if cmd.String("password") == "" {
		return
	}

	user, err := a.Conn.Query.GetUserByUsername(ctx, cmd.String("user"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("failed to get user", "user", cmd.String("user"))
		} else {
			slog.Error("failed to get user", "error", err)
		}
		os.Exit(1)
	}
	hash, err := util.HashPassword(cmd.String("password"))
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		os.Exit(1)
	}
	if err = a.Conn.Query.UpdateUserPassword(ctx, &db.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: hash,
	}); err != nil {
		slog.Error("failed to update password for user", "user", cmd.String("user"), "error", err)
		os.Exit(1)
	}

	slog.Info("Reset successful!", "user", cmd.String("user"), "password", cmd.String("password"))
	os.Exit(1)
}
