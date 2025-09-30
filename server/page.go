package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/templates/page"
)

func MainPage(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resumes, err := database.GetResumes(r.Context(), 0 /* TODO: Set to userID */)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for names: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		applications, err := database.GetApplicationHistory(r.Context(), 0 /* TODO: Set to userID */)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for applications: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		if len(applications) == 0 {
			logger.Warn("NO APPLICATIONS")
		}

		page.Home(page.HomeProps{
			Resumes:      resumes,
			Applications: applications,
		}).Render(r.Context(), w)
	}
}

func ResumeFormPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page.ResumeForm(page.ResumeFormProps{}).Render(r.Context(), w)
	}
}

func ApplyPage(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("invalid base resume id: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		base, err := database.GetResume(r.Context(), db.GetResumeParams{
			UserID: 0, /* TODO: Set to userID */
			ID:     id,
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for base resume: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		page.Apply(page.ApplyProps{
			Base: base,
		}).Render(r.Context(), w)
	}
}
