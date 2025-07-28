package serve

import (
	"context"
	"fmt"
	"github.com/ohhfishal/resume-wizard/server"
	"log/slog"
)

type Cmd struct {
	// TODO: Add port
}

func (cmd *Cmd) Run(ctx context.Context, logger *slog.Logger) error {
	s, err := server.New(logger)
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
