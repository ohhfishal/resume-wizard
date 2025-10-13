package server

import (
	"net/http"
)

func HandlePostMatch(store ResumeStore) Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return TextStatus(http.StatusOK)
	})
}
