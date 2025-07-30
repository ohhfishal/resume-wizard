package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
)

// Date,Company,Position,Resume,Date,Status,Notes
type LogCmd struct {
	// Subcommands
	Apply ApplyCmd `cmd:"" help:"Submitted a job applicaton."`
	// TODO: Commands
	// export - Dump out as CSV or JSON
	// update - Update fields
	// delete - Delete a row
}

type ApplyCmd struct {
	File *os.File `arg:"" required:"" help:"File used (Must match \"*.yaml\", \"*.json\" or \"-\")"`
	// TODO: Improve defaults
	DatabaseSource string `short:"s" default:":memory:" help:"Database connection string (sqlite)."`
	Company        string `arg:"" required:"" help:"Company applied to."`
	Position       string `arg:"" required:"" help:"Position applied to."`
}

func (cmd *ApplyCmd) Run(ctx context.Context, logger *slog.Logger) error {
	defer cmd.File.Close()

	entry, err := resume.FromFiles(cmd.File)
	if err != nil {
		return fmt.Errorf("getting resume from input: %w", err)
	}

	database, err := db.Open(ctx, "sqlite3", cmd.DatabaseSource)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	// TODO: Transaction this?
	id, err := database.InsertResume(ctx, db.InsertResumeParams{
		Name: "TODO",
		Body: &entry,
	})
	if err != nil {
		return fmt.Errorf("inserting resume into database: %w", err)
	}

	application, err := database.InsertLog(ctx, db.InsertLogParams{
		ResumeID: id,
		Company:  cmd.Company,
		Position: cmd.Position,
	})
	if err != nil {
		return fmt.Errorf("inserting log into database: %w", err)
	}
	logger.Info("done", slog.Any("log", application))

	return nil
}
