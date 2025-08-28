package server

import (
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"github.com/ohhfishal/resume-wizard/templates/card"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PostApplicationHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("resume")

		var applicationResume resume.Resume
		if err := resume.FromYAML(strings.NewReader(body), &applicationResume); err != nil {
			http.Error(w,
				fmt.Sprintf("parsing resume: %s", err.Error()),
				http.StatusBadRequest,
			)
		}

		tx, dbtx, err := database.BeginTx(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("starting transaction: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		session, err := dbtx.GetSession(r.Context(), db.GetSessionParams{
			UserID: 0, /* TODO: Set to userID */
			Uuid:   r.PathValue("session_id"),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("getting session: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		_, err = dbtx.CreateApplication(r.Context(), db.CreateApplicationParams{
			UserID:       session.UserID,
			BaseResumeID: session.BaseResumeID,
			Company:      session.Company,
			Position:     session.Position,
			Description:  session.Description,
			Resume:       &applicationResume,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("creating application: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if err := dbtx.SoftDeleteSession(r.Context(), db.SoftDeleteSessionParams{
			UserID: 0, /* TODO: Set to userID */
			Uuid:   session.Uuid,
		}); err != nil {
			http.Error(w, fmt.Sprintf("deleting session: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, fmt.Sprintf("commit transaction: %s", err.Error()), http.StatusInternalServerError)
		}

		w.Header().Set("HX-Redirect", "/")
	}
}

func PutApplicationHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid user_id: %s", err.Error()), http.StatusBadRequest)
			return
		}
		// TODO: Validate userID is who is asking

		app, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid application id: %s", err.Error()), http.StatusBadRequest)
			return
		}

		applied, err := time.Parse(time.DateOnly, r.FormValue("applied_at"))
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid applied_at: %s", err.Error()), http.StatusBadRequest)
			return
		}

		status := r.FormValue("status")
		row, err := database.UpdateApplication(r.Context(),
			db.UpdateApplicationParams{
				Status:    status,
				AppliedAt: applied,
				UserID:    user,
				ID:        app,
			},
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("updating row: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		card.ApplicationsRow(row).Render(r.Context(), w)
	}
}
