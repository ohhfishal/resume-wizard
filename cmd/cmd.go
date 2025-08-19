package cmd

import (
	"context"
	"io"
	"log/slog"

	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
)

type RootCmd struct {
	Config kong.ConfigFlag `short:"c" help:"Path to config file to load." type:"path"`
	Test string `short:"t"`
	Build  BuildCmd  `cmd:"" help:"Compile a resume from a input file."`
	Log    LogCmd    `cmd:"" help:"Log usuage of resumes."`
	Serve  ServeCmd  `cmd:"" help:"Run resume-wizard as a local HTTP serve."`
	Wizard WizardCmd `cmd:"" help:"Do some magic."`
}

func Run(ctx context.Context, stdout io.Writer, args []string) error {
	cmd := &RootCmd{}

	var exit bool
	parser, err := kong.New(
		cmd,
		kong.Exit(func(_ int) { exit = true }),
		kong.BindTo(ctx, new(context.Context)),
		kong.Configuration(kongyaml.Loader),
	)
	if err != nil {
		return err
	}

	parser.Stdout = stdout
	parser.Stderr = stdout

	context, err := parser.Parse(
		args,
	)
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
