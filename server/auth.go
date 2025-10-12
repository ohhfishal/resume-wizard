package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ohhfishal/resume-wizard/db"
	"net/http"
)

type User db.User

type UserStore interface {
	GetUser(ctx context.Context, id string) (db.User, error)
	CreateUser(ctx context.Context, id string) (db.User, error)
}

func WithBearerAuth(store UserStore) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		auth := r.Header.Get("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			return Text("Missing Bearer Auth", http.StatusUnauthorized)
		}
		token := auth[7:]
		if token == "" {
			return Text("Empty Auth", http.StatusUnauthorized)
		}

		user, err := store.GetUser(r.Context(), token)
		if err != nil {
			return TextStatus(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "user", user)
		*r = *r.WithContext(ctx)
		return Next
	})
}

func RegisterUser(store UserStore) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		user, err := store.CreateUser(r.Context(), uuid.NewString())
		if err != nil {
			return Text(
				fmt.Errorf("could not create user: %w", err).Error(),
				http.StatusInternalServerError,
			)
		}
		return JSON(user, http.StatusCreated)
	})
}
