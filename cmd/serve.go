package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ohhfishal/resume-wizard/feature"
	"github.com/ohhfishal/resume-wizard/server"
)

type ServeCmd struct {
	Config server.Config `embed:""`
}

func (cmd *ServeCmd) Run(ctx context.Context, logger *slog.Logger) error {
	s, err := server.New(ctx, cmd.Config, logger)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}

	feature.SetFeatures(cmd.Config.Features)

	if err := s.Run(ctx); err != nil {
		logger.Error(
			"running server",
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}
