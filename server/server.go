package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ohhfishal/resume-wizard/assets"
)

type Server struct {
	logger *slog.Logger
	port   string
}

func New(logger *slog.Logger) (*Server, error) {
	return &Server{
		logger: logger,
		// TODO: Fix hardcoding
		port: "8080",
	}, nil
}

func (server *Server) Run(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(loggingMiddleware(server.logger))
	r.Use(middleware.Recoverer)

	r.Mount(
		"/assets",
		http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement
		MainPage().Render(r.Context(), w)
	})

	r.Post("/resume", NewUploadHandler(server.logger))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	})

	s := &http.Server{
		Addr:         ":" + server.port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		server.logger.Info("shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			server.logger.Error("closing server",
				"error", err,
			)
		}
	}()

	server.logger.Info("starting server", "port", server.port)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)

			logger.Info("replied to request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.statusCode,
				"duration", time.Since(start).String(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
