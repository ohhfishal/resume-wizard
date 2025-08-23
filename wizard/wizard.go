package wizard

import (
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
)

type Wizard struct {
}

func (wizard *Wizard) Annotate(base db.BaseResume) (*resume.Resume, error) {
	// TODO: IMPLEMENT THIS
	return base.Resume, nil
}
