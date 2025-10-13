package server

import (
	"context"
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"net/http"
)

type ResumeStore interface {
	CreateResume(ctx context.Context, arg db.CreateResumeParams) (db.Resume, error)
	GetResume(ctx context.Context, arg db.GetResumeParams) (db.Resume, error)
	GetResumes(ctx context.Context, userID string) ([]db.Resume, error)
}

func HandlePostResume(store ResumeStore) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		user, ok := r.Context().Value("user").(User)
		if !ok {
			return Text("Missing user field", http.StatusBadRequest)
		}

		resumeInput, err := Decode[resume.Resume](r)
		if err != nil {
			return Text(fmt.Errorf("invalid resume: %w", err).Error(), http.StatusBadRequest)
		}

		resumeEntry, err := store.CreateResume(r.Context(), db.CreateResumeParams{
			UserID: user.ID,
			Name:   resumeInput.Title,
			Resume: &resumeInput,
		})

		if err != nil {
			return InternalServerError(fmt.Errorf("could not create resume: %w", err))
		}
		return JSON(resumeEntry, http.StatusCreated)
	})
}
