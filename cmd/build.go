package cmd

import (
	"fmt"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"os"
)

type BuildCmd struct {
	Inputs            []*os.File `arg:"" required:"" help:"Input file to convert (Must match \"*.yaml\", \"*.json\" or \"-\")"`
	Output            string     `short:"o" optional:"" enum:"html,yaml,json" default:"html" help:"Output format (\"html\",\"yaml\",\"json\")"`
	ApplyPersonalInfo *os.File   `short:"p" group:"Personal Info:" placeholder:"FILE" help:"Apply personal info from a file to remove redactions."`
	HidePersonalInfo  bool       `short:"q" group:"Personal Info:" help:"Hide personal info"`
	// TODO: Implement this instead of defaulting to os.Stdout
	// Output io.Writer
}

// TODO: Careful where we log to since this can emit to stdout
func (cmd *BuildCmd) Run(logger *slog.Logger) error {
	for _, file := range cmd.Inputs {
		defer file.Close()
	}

	entry, err := resume.FromFiles(cmd.Inputs)
	if err != nil {
		return fmt.Errorf("creating resume from input: %w", err)
	}

	if cmd.ApplyPersonalInfo != nil {
		err := entry.ApplyPatch(cmd.ApplyPersonalInfo)
		if err != nil {
			return fmt.Errorf("applying personal info patch: %w", err)
		}
	}

	if cmd.HidePersonalInfo {
		entry.HidePersonalInfo()
	}

	switch cmd.Output {
	case "html":
		if err := entry.ToHTML(os.Stdout); err != nil {
			return fmt.Errorf("converting to HTML: %w", err)
		}
	case "yaml":
		if err := entry.ToYAML(os.Stdout); err != nil {
			return fmt.Errorf("converting to YAML: %w", err)
		}
	case "json":
		if err := entry.ToJSON(os.Stdout); err != nil {
			return fmt.Errorf("converting to YAML: %w", err)
		}
	default:
		return fmt.Errorf("unknown output format: %s", cmd.Output)
	}

	return nil
}
