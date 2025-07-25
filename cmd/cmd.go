package cmd

import (
	"context"
	"io"
	"log/slog"

	"github.com/alecthomas/kong"
	"github.com/ohhfishal/resume-wizard/cmd/build"
	"github.com/ohhfishal/resume-wizard/cmd/wizard"
)

type RootCmd struct {
	Build  build.Cmd  `cmd:"" help:"Compile a resume from a input file."`
	Wizard wizard.Cmd `cmd:"" help:"Do some magic"`
}

func Run(ctx context.Context, stdout io.Writer, args []string) error {
	cmd := &RootCmd{}

	var exit bool
	parser, err := kong.New(
		cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.BindTo(ctx, new(context.Context)),
	)
	if err != nil {
		return err
	}

	parser.Stdout = stdout
	parser.Stderr = stdout

	context, err := parser.Parse(args)
	if err != nil || exit {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // TODO: Make configurable
	}))
	slog.SetDefault(logger)

	if err := context.Run(logger); err != nil {
		return err
	}
	return nil
}
