package build

import (
	"fmt"
	"log/slog"
	"os"
	"github.com/ohhfishal/resume-wizard/ir"
)

type Cmd struct {
	Input *os.File `arg:"" required:"" help:"Input file to convert (Must match \"*.yaml\", \"*.json\" or \"-\")"`
	// TODO: Implement this instead of defaulting to os.Stdout
	// Output io.Writer
}


// TODO: Careful where we log to since this can emit to stdout
func (cmd *Cmd) Run(logger *slog.Logger) error {
	defer cmd.Input.Close()
	resume, err := ir.FromFile(cmd.Input)
	if err != nil {
		return fmt.Errorf("creating resume from input: %w", err)
	}
	if err := resume.ToHTML(os.Stdout); err != nil {
		return fmt.Errorf("converting to HTML: %w", err)
	}
	return nil
}
