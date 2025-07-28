package server

import (
	"fmt"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"net/http"
)

const MaxFileSize = 12_000 // 12KB

var UploadFileTypes = []string{
	"application/yaml",
	"application/json",
}

func NewUploadHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile(UploadFileKey)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading file: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}

		if header.Size >= MaxFileSize {
			http.Error( w,
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

		logger.Debug("got a new resume",
			slog.Any("resume", newResume),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
