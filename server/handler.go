package server

import (
	"encoding/json"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) Handler

var _ http.Handler = Handler(nil)
var Next Handler = nil

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next := h(w, r)
	if next != nil {
		next.ServeHTTP(w, r)
	}
}

func Chain(handlers ...Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		for _, h := range handlers {
			next := h(w, r)
			if next != nil {
				return next
			}
		}
		return nil
	}
}

func TextStatus(status int) Handler {
	return Text(http.StatusText(status), status)
}

func Text(body string, status int) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return nil
	}
}

func JSON(data any, status int) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(data)
		return nil
	}
}

func HandleHealth() http.Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return TextStatus(http.StatusOK)
	})
}

func HandleNotFound() http.Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) Handler {
		return TextStatus(http.StatusNotFound)
	})
}
