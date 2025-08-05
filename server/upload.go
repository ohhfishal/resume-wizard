package server

import (
	"errors"
	"fmt"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const MaxFileSize = 12_000 // 12KB

const (
	EventResumeUploaded     = "resumeUploaded"
	EventApplicationsUpdate = "applicationsUpdate"
)

var UploadFileTypes = []string{
	"application/yaml",
	"application/json",
}

type PostApplicationInput struct {
	ResumeID int64
	Company  string
	Position string
	Status   string
}

func Parse(r *http.Request) (*PostApplicationInput, error) {
	id := r.FormValue("resume_id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid field id: %w", err)
	}

	company := r.FormValue("company")
	position := r.FormValue("position")
	status := r.FormValue("status")

	if company == "" {
		return nil, errors.New("missing field: company")
	} else if position == "" {
		return nil, errors.New("missing field: position")
	} else if status == "" {
		return nil, errors.New("missing field: position")
	} else if status != "pending" {
		return nil, errors.New("non-pending status not implemented")
	}

	return &PostApplicationInput{
		ResumeID: intID,
		Company:  company,
		Position: position,
		Status:   status,
	}, nil
}

func PostApplicationHandler(logger *slog.Logger, database *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inputs, err := Parse(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := database.InsertLog(r.Context(), db.InsertLogParams{
			ResumeID: inputs.ResumeID,
			Company:  inputs.Company,
			Position: inputs.Position,
		}); err != nil {
			http.Error(w,
				fmt.Sprintf("could not insert into database: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("HX-Trigger", EventApplicationsUpdate)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func PutResumeHandler(logger *slog.Logger, database *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, fmt.Errorf("invalid id: %w", err).Error(), http.StatusBadRequest)
			return
		}

		body := r.FormValue("resume")

		var updatedResume resume.Resume
		if err := resume.FromYAML(strings.NewReader(body), &updatedResume); err != nil {
			http.Error(w,
				fmt.Sprintf("parsing resume: %s", err.Error()),
				http.StatusBadRequest,
			)
		}

		if err := database.UpdateResume(r.Context(), db.UpdateResumeParams{
			ID:   intID,
			Body: &updatedResume,
		}); err != nil {
			http.Error(w,
				fmt.Sprintf("updating resume: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func PostResumeHandler(logger *slog.Logger, database *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Handle this smarter. Don't allow repeats
		name := r.FormValue(NameKey)
		if name == "" {
			http.Error(w,
				"missing field: name",
				http.StatusBadRequest,
			)
			return
		}

		file, header, err := r.FormFile(UploadFileKey)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading file: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}

		if header.Size >= MaxFileSize {
			http.Error(w,
				fmt.Sprintf("file too big: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}

		var contentType string
		for _, userType := range header.Header.Values("Content-Type") {
			for _, allowedType := range UploadFileTypes {
				if allowedType == userType {
					contentType = allowedType
				}
			}
		}

		var newResume resume.Resume
		switch contentType {
		case "application/json":
			if err := resume.FromJSON(file, &newResume); err != nil {
				http.Error(w,
					fmt.Sprintf("parsing json: %s", err.Error()),
					http.StatusBadRequest,
				)
			}
		case "application/yaml":
			if err := resume.FromYAML(file, &newResume); err != nil {
				http.Error(w,
					fmt.Sprintf("parsing yaml: %s", err.Error()),
					http.StatusBadRequest,
				)
			}
		default:
			http.Error(w,
				fmt.Sprintf(
					"invalid content-type for file must be one of %v",
					UploadFileTypes,
				),
				http.StatusBadRequest,
			)
			return
		}

		if _, err := database.InsertResume(r.Context(), db.InsertResumeParams{
			Name: name,
			Body: &newResume,
		}); err != nil {
			http.Error(w,
				fmt.Sprintf("could not insert into database: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		logger.Debug("got a new resume",
			slog.Any("resume", newResume),
		)

		w.Header().Set("HX-Trigger", EventResumeUploaded)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
