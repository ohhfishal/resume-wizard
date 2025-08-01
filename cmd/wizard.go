package cmd

import (
	// "fmt"
	// "io"
	"log/slog"
	"os"
)

type WizardCmd struct {
	File *os.File `arg:"" help:"File to read"`
	AnthropicKey string `env:"ANTHROPIC_API_KEY" help:"API key to use Claude."`
}

func (cmd *WizardCmd) Run(logger *slog.Logger) error {
	defer cmd.File.Close()
	return nil
}
