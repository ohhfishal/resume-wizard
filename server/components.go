package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/templates/components"
	"log/slog"
	"net/http"
)

func componentError(w http.ResponseWriter, err error, status int) {
	// TODO: FIX THESE! IE Have its errors be valid HTML
	http.Error(w, err.Error(), status)
}

func ComponentsHandler(logger *slog.Logger, database *db.DB) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/tailoredResumeSection/{uuid}", func(w http.ResponseWriter, r *http.Request) {
			session, err := database.GetSession(r.Context(), db.GetSessionParams{
				Uuid:   r.PathValue("uuid"),
				UserID: 0, /* TODO: SET userID */
			})
			if err != nil {
				componentError(w,
					fmt.Errorf("reading database for base resumes: %s", err),
					http.StatusInternalServerError,
				)
				return
			} else if session.Resume == nil {
				// TODO: Set from the config
				w.Header().Set("Retry-After", "1")
				componentError(w,
					fmt.Errorf("reading database for base resumes: %s", err),
					http.StatusServiceUnavailable,
				)
				return
			}
			components.TailoredResumeSection(components.TailoredResumeSectionProps{
				Session: session,
			}).Render(r.Context(), w)
		})
	}
}
