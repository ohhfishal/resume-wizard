package build

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"github.com/ohhfishal/resume-wizard/ir"
)

type Cmd struct {
	Input Input `embed:""`
}

type Input struct {
	File *os.File `arg:"" required:"" help:"Input file to convert (Must match \"*.yaml\", \"*.json\" or \"-\")"`
}

func (input Input) Validate() error {
	return errors.New("testing")
}

func (cmd *Cmd) Run(logger *slog.Logger) error {
	defer cmd.Input.File.Close()
	_, err := ir.FromFile(cmd.Input.File)
	if err != nil {
		return fmt.Errorf("creating resume from input: %w", err)
	}
	return nil
}
