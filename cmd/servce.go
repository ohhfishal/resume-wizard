package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/server"
)

type ServeCmd struct {
	// TODO: Add port
	DatabaseSource string `short:"s" default:":memory:" help:"Database connection string (sqlite)."`
}

func (cmd *ServeCmd) Run(ctx context.Context, logger *slog.Logger) error {
	database, err := db.Open(ctx, "sqlite3", cmd.DatabaseSource)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	s, err := server.New(logger, database)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	if err := s.Run(ctx); err != nil {
		logger.Error(
			"running server",
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}
