package server

import (
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"github.com/ohhfishal/resume-wizard/templates"
	"github.com/ohhfishal/resume-wizard/templates/card"
	"log/slog"
	"net/http"
	"strings"
)

func FormFileResume(r *http.Request, key string) (*resume.Resume, error) {
	file, header, err := r.FormFile(templates.UploadFileKey)
	if err != nil {
		return nil, fmt.Errorf("parsing form: %w", err)
	}

	if header.Size >= MaxFileSize {
		return nil, fmt.Errorf("file too big: %d", header.Size)
	}

	newResume, err := resume.FromContentType(file, header.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}
	return newResume, nil
}

func GetBaseResumeForm(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var baseResume *resume.Resume
		if r.Method == http.MethodPost {
			var err error
			baseResume, err = FormFileResume(r, templates.UploadFileKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		card.BaseResumeReviewForm(baseResume).Render(r.Context(), w)
	}
}

func PostBaseResumeHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("name")
		if title == "" {
			http.Error(w, "missing field: name", http.StatusBadRequest)
			return
		}

		body := r.FormValue("resume")

		var baseResume resume.Resume
		if err := resume.FromYAML(strings.NewReader(body), &baseResume); err != nil {
			http.Error(w,
				fmt.Sprintf("parsing resume: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}

		if _, err := database.InsertBase(r.Context(), db.InsertBaseParams{
			UserID: 0, // TODO: Grab this from somewhere (Probably the context?)
			Name:   title,
			Resume: &baseResume,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Redirect", "/home")
		w.WriteHeader(http.StatusSeeOther)
	}
}
