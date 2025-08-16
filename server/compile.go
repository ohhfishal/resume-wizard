package server

import (
	"fmt"
	// "github.com/ohhfishal/resume-wizard/components"
	"github.com/ohhfishal/resume-wizard/db"
	// "github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"net/http"
	"strconv"
)

type PostCompileForm struct {
	ResumeID int64
	// PatchID string
}

func (form *PostCompileForm) Parse(r *http.Request) error {
	id, err := strconv.ParseInt(r.FormValue("resume_id"), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid field id: %w", err)
	}
	form.ResumeID = id
	return nil
}

func PostCompileHandler(logger *slog.Logger, database *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request PostCompileForm
		if err := request.Parse(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entry, err := database.GetResumeByID(r.Context(), request.ResumeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Apply patch

		if err := entry.Body.ToHTML(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
