package build

import (
	"fmt"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"os"
)

type Cmd struct {
	Input            *os.File `arg:"" required:"" help:"Input file to convert (Must match \"*.yaml\", \"*.json\" or \"-\")"`
	HidePersonalInfo bool     `short:"q" help:"Hide personal info"`
	// TODO: Implement this instead of defaulting to os.Stdout
	// Output io.Writer
}

// TODO: Careful where we log to since this can emit to stdout
func (cmd *Cmd) Run(logger *slog.Logger) error {
	defer cmd.Input.Close()
	entry, err := resume.FromFile(cmd.Input)
	if err != nil {
		return fmt.Errorf("creating resume from input: %w", err)
	}

	if cmd.HidePersonalInfo {
		entry.HidePersonalInfo()
	}

	if err := entry.ToHTML(os.Stdout); err != nil {
		return fmt.Errorf("converting to HTML: %w", err)
	}
	return nil
}
