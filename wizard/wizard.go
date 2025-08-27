package wizard

import (
	"context"
	"errors"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"time"
)

type Wizard struct {
	Backend string `enum:"disabled,anthropic,sleep" default:"sleep" help:"Wizard backend"`
	Sleep   struct {
		Duration time.Duration `default:"1s" help:"How long to sleep to mimic doing work"`
	} `embed:""`
	Claude struct {
		ApiKey string            `help:"API Key to access LLM (KEEP SECRET!)" env:"API_KEY"`
		client *anthropic.Client `kong:"-"`
	} `embed:"anthropic-" envprefix:"ANTHROPIC_"`
	logger *slog.Logger `kong:"-"`
}

func (wizard *Wizard) Init(logger *slog.Logger) error {
	wizard.logger = logger
	switch wizard.Backend {
	case "anthropic":
		if wizard.Claude.ApiKey == "" {
			return errors.New("missing field: ApiKey")
		}
		client := anthropic.NewClient(
			option.WithAPIKey(wizard.Claude.ApiKey),
		)
		wizard.Claude.client = &client
		wizard.Claude.ApiKey = ""
	}
	return nil
}

type AnnotationContext struct {
	Base        db.BaseResume
	Company     string
	Position    string
	Description string
}

func (wizard *Wizard) Annotate(ctx context.Context, args AnnotationContext) (*resume.Resume, error) {
	start := time.Now()
	wizard.logger.Info("creating annotation",
		slog.String("base", args.Base.Name),
		slog.String("position", args.Position),
		slog.Time("start", time.Now()),
	)

	defer func() {
		wizard.logger.Info("returning",
			slog.Time("start", start),
			slog.Time("end", time.Now()),
			slog.Duration("duration", time.Since(start)),
		)
	}()
	switch wizard.Backend {
	case "anthropic":
		return wizard.annotateClaude(ctx, args)
	case "sleep":
		time.Sleep(wizard.Sleep.Duration)
		fallthrough
	case "disabled":
		return args.Base.Resume, nil
	default:
		return nil, fmt.Errorf("unknown backend: %s", wizard.Backend)
	}
}
