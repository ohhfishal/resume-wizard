package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ohhfishal/resume-wizard/components"
	"github.com/ohhfishal/resume-wizard/db"
	"log/slog"
	"net/http"
	"strconv"
)

func componentError(w http.ResponseWriter, err error, status int) {
	// TODO: FIX THESE! IE Have its errors be valid HTML
	http.Error(w, err.Error(), status)
}

func ComponentsHandler(logger *slog.Logger, database *db.Queries) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/applicationsTable", func(w http.ResponseWriter, r *http.Request) {
			resumes, err := database.GetResumes(r.Context())
			if err != nil {
				componentError(w,
					fmt.Errorf("reading database for names: %s", err),
					http.StatusInternalServerError,
				)
				return
			}

			applications, err := database.GetApplications(r.Context())
			if err != nil {
				componentError(w,
					fmt.Errorf("reading database for applications: %s", err),
					http.StatusInternalServerError,
				)
				return
			}
			components.ApplicationsTable(resumes, applications).Render(r.Context(), w)
		})
		r.Get("/resumeDropdown", func(w http.ResponseWriter, r *http.Request) {
			resumes, err := database.GetResumes(r.Context())
			if err != nil {
				componentError(w,
					fmt.Errorf("reading database for names: %s", err),
					http.StatusInternalServerError,
				)
				return
			}
			listener := r.URL.Query().Get("listener")
			components.ResumeDropdown(resumes, listener).Render(r.Context(), w)
		})
		r.Get("/resumeEditor", func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("resume_id")
			if id == "" {
				components.ResumeEditor(db.Resume{ID: -1}).Render(r.Context(), w)
				return
			}
			intID, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				componentError(w,
					fmt.Errorf("could not conver id to int: %s", err),
					http.StatusBadRequest,
				)
				return
			}

			resume, err := database.GetResumeByID(r.Context(), intID)
			if err != nil {
				componentError(w,
					fmt.Errorf("reading database for resume: %s", err),
					http.StatusInternalServerError,
				)
				return
			}
			components.ResumeEditor(resume).Render(r.Context(), w)
		})
	}
}
