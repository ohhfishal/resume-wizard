package server

import (
	"encoding/csv"
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	"log/slog"
	"net/http"
	"time"
)

func GetExportHandler(logger *slog.Logger, database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := r.PathValue("format")
		switch format {
		case "csv":
			applications, err := database.GetApplicationHistory(r.Context(), 0 /* TODO: Replace with userID */)
			if err != nil {
				http.Error(w, fmt.Sprintf("getting applications: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/csv")
			writer := csv.NewWriter(w)
			defer writer.Flush()

			if err := writer.Write([]string{
				"Resume ID",
				"Company",
				"Position",
				"Description",
				"Status",
				"Created At",
			}); err != nil {
				http.Error(w, fmt.Sprintf("writing application: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			for _, app := range applications {
				var description string
				// TODO: Set the description
				description = "TODO"
				if err := writer.Write([]string{
					fmt.Sprintf("%d", app.ResumeID),
					app.Company,
					app.Position,
					description,
					app.Status,
					app.CreatedAt.Format(time.DateOnly),
				}); err != nil {
					http.Error(w, fmt.Sprintf("writing application: %s", err.Error()), http.StatusInternalServerError)
					return
				}
			}
		default:
			http.Error(w, fmt.Sprintf("unknown format: %s", format), http.StatusInternalServerError)
			return
		}
	}
}
