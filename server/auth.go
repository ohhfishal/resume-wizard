package server

import (
	"context"
	"net/http"
)

type User struct {
}

type UserStore interface {
	LookupUser(ctx context.Context, user string) (*User, error)
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

		user, err := store.LookupUser(r.Context(), token)
		if err != nil {
			return TextStatus(http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "user", user)
		*r = *r.WithContext(ctx)
		return Next
	})
}
