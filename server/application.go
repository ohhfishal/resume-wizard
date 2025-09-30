package server

import (
	"database/sql"
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	// "github.com/ohhfishal/resume-wizard/resume"
	// "github.com/ohhfishal/resume-wizard/templates/card"
	"log/slog"
	"net/http"
	"strconv"
	// "strings"
	// "time"
)

func PostApplicationHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user int64 // TODO: FIXME: Add user

		resumeID, err := strconv.ParseInt(r.FormValue("resume_id"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid resume_id: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if _, err := database.GetResume(r.Context(), db.GetResumeParams{
			UserID: user,
			ID:     resumeID,
		}); err != nil {
			http.Error(w, fmt.Sprintf("could not get resume: %d", resumeID), http.StatusBadRequest)
			return
		}

		status := r.FormValue("status")
		if status == "" {
			status = "pending"
		}

		descriptionString := r.FormValue("description")
		description := sql.NullString{
			String: descriptionString,
			Valid:  len(descriptionString) > 0,
		}

		params := db.InsertApplicationActivityParams{
			UserID:      user,
			ResumeID:    resumeID,
			Company:     r.FormValue("company_name"),
			Position:    r.FormValue("title"),
			Status:      status,
			Description: description,
		}
		if _, err := database.InsertApplicationActivity(r.Context(), params); err != nil {
			http.Error(w, fmt.Sprintf("inserting application: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
	}
}

func PutApplicationHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusInternalServerError)
		return
		// 	user, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf("invalid user_id: %s", err.Error()), http.StatusBadRequest)
		// 		return
		// 	}
		// 	// TODO: Validate userID is who is asking
		//
		// 	app, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf("invalid application id: %s", err.Error()), http.StatusBadRequest)
		// 		return
		// 	}
		//
		// 	applied, err := time.Parse(time.DateOnly, r.FormValue("applied_at"))
		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf("invalid applied_at: %s", err.Error()), http.StatusBadRequest)
		// 		return
		// 	}
		//
		// 	status := r.FormValue("status")
		// 	row, err := database.UpdateApplication(r.Context(),
		// 		db.UpdateApplicationParams{
		// 			Status:    status,
		// 			AppliedAt: applied,
		// 			UserID:    user,
		// 			ID:        app,
		// 		},
		// 	)
		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf("updating row: %s", err.Error()), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	card.ApplicationsRow(row).Render(r.Context(), w)
	}
}
