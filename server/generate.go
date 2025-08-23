package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/templates/page"
	"github.com/ohhfishal/resume-wizard/wizard"
	"log/slog"
	"net/http"
	"strconv"
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
		base, err := database.GetBaseResume(r.Context(), db.GetBaseResumeParams{
			UserID: 0, /* TODO: Set to userID */
			ID:     form.BaseResumeID,
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for base resume: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		annotated, err := model.Annotate(base)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("annotating resume: %s", err.Error()),
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
			Resume:       annotated,
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("creating session: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("HX-Redirect", fmt.Sprintf("/tailor/%s", session.Uuid))
		w.WriteHeader(http.StatusOK)
		return

		page.TailorResume(page.TailorResumeProps{
			Base:            base,
			Session:         session,
			LockApplication: true,
		}).Render(r.Context(), w)
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
