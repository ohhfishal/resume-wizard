package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/wizard"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type GenerateForm struct {
	BaseResumeID int64
	Company      string
	Title        string
	Description  string
	UUID         string
}

func GenerateHandler(logger *slog.Logger, database *db.DB, model *wizard.Wizard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, err := ParseGenerateForm(r)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("invalid base resume id: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		session, err := database.CreateSession(r.Context(), db.CreateSessionParams{
			Uuid:         form.UUID,
			BaseResumeID: form.BaseResumeID,
			UserID:       0, /* TODO: Grab from somewhere */
			Company:      form.Company,
			Position:     form.Title,
			Description:  form.Description,
			Resume:       nil,
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("creating session: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		go func(ctx context.Context, logger *slog.Logger, session db.Session) {
			// TODO: Have this be configurable
			ctx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()

			base, err := database.GetBaseResume(ctx, db.GetBaseResumeParams{
				UserID: session.UserID,
				ID:     session.BaseResumeID,
			})
			if err != nil {
				logger.Error("failed to read base resume", slog.Any("error", err))
				return
			}
			annotated, err := model.Annotate(ctx, wizard.AnnotationContext{
				Base:        base,
				Company:     form.Company,
				Position:    form.Title,
				Description: form.Description,
			})
			if err != nil {
				logger.Error("failed to annotate resume", slog.Any("error", err))
				return
			}
			if err := database.AddResumeToSession(ctx, db.AddResumeToSessionParams{
				Resume: annotated,
				UserID: session.UserID,
				Uuid:   session.Uuid,
			}); err != nil {
				logger.Error("failed to update session", slog.Any("error", err))
				return
			}
			// TODO: Mark the session as done
			logger.Debug("tailored resume successfully",
				slog.String("uuid", session.Uuid),
			)
		}(context.Background(), logger, session)

		w.Header().Set("HX-Redirect", fmt.Sprintf("/tailor/%s", session.Uuid))
		w.WriteHeader(http.StatusOK)
	}
}

func ParseGenerateForm(r *http.Request) (*GenerateForm, error) {
	rawID := r.PostFormValue("base_resume_id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid field id: %w", err)
	}

	form := &GenerateForm{
		BaseResumeID: id,
		Company:      r.PostFormValue("company_name"),
		Title:        r.PostFormValue("title"),
		Description:  r.PostFormValue("description"),
		UUID:         uuid.NewString(),
	}

	if form.Company == "" {
		return nil, errors.New(`missing field: "company_name"`)
	} else if form.Title == "" {
		return nil, errors.New(`missing field: "title"`)
	} else if form.Description == "" {
		return nil, errors.New(`missing field: "description"`)
	}
	return form, nil
}
