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
	LoggerConfig struct {
		Level string `enum="disabled,error,warn,info,debug" placeholder:"info" default="info" help:"Log level to use. ('disabled','error','warn','info','debug')"`
	} `embed:"" prefix:"logger-"`
	Build  BuildCmd        `cmd:"" help:"Compile a resume from a input file."`
	Serve  ServeCmd        `cmd:"" help:"Run resume-wizard as a local HTTP serve."`
	Log    LogCmd          `cmd:"" help:"Log usuage of resumes."`
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

	logger := cmd.NewLogger(stdout)
	slog.SetDefault(logger)

	if err := context.Run(logger); err != nil {
		return err
	}
	return nil
}

func (cmd RootCmd) NewLogger(stdout io.Writer) *slog.Logger {
	var level slog.Level
	switch cmd.LoggerConfig.Level {
	case "disabled":
		return slog.New(slog.DiscardHandler)
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "debug":
		level = slog.LevelDebug
	case "info":
		fallthrough
	default:
		level = slog.LevelInfo
	}
	return  slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
