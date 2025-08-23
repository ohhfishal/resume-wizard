package wizard

import (
	"context"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
)

type Wizard struct {
	Enabled bool         `default:"true" negatable:"wizard-disable" help:"Enable using LLM annotation (default=true)"`
	URL     string       `help:"LLM API url"`
	ApiKey  string       `help:"API Key to access LLM (KEEP SECRET!)" env:"API_KEY"`
	logger  *slog.Logger `kong:"-"`
}

func (wizard *Wizard) Init(logger *slog.Logger) error {
	wizard.logger = logger
	return nil
}

type AnnotationContext struct {
	Base        db.BaseResume
	Company     string
	Position    string
	Description string
}

func (wizard *Wizard) Annotate(ctx context.Context, args AnnotationContext) (*resume.Resume, error) {
	if !wizard.Enabled {
		return args.Base.Resume, nil
	}
	wizard.logger.Info("creating annotation",
		slog.String("base", args.Base.Name),
		slog.String("position", args.Position),
	)
	// TODO: IMPLEMENT THIS
	return args.Base.Resume, nil
}
