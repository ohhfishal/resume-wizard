package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ohhfishal/resume-wizard/assets"
	"github.com/ohhfishal/resume-wizard/db"
)

type Server struct {
	logger   *slog.Logger
	database *db.Queries
	host     string
	port     string
}

func New(logger *slog.Logger, database *db.Queries) (*Server, error) {
	return &Server{
		database: database,
		logger:   logger,
		host:     "0.0.0.0", // TODO: Fix hardcoding
		port:     "8080",    // TODO: Fix hardcoding
	}, nil
}

func (server *Server) Run(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(loggingMiddleware(server.logger))
	r.Use(middleware.Recoverer)

	r.Put("/resume/{id}", PutResumeHandler(server.logger, server.database))
	r.Post("/resume", PostResumeHandler(server.logger, server.database))
	r.Post("/application", PostApplicationHandler(server.logger, server.database))

	r.Mount(
		"/assets",
		http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))),
	)

	r.Route("/components", ComponentsHandler(server.logger, server.database))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resumes, err := server.database.GetResumes(r.Context())
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for names: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		applications, err := server.database.GetApplications(r.Context())
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for applications: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		server.logger.Debug(
			"got from db",
			"names", resumes,
			"apps", applications,
		)

		// TODO: Implement
		MainPage(resumes, applications).Render(r.Context(), w)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
	})

	s := &http.Server{
		Addr:         net.JoinHostPort(server.host, server.port),
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

	server.logger.Info("starting server", "port", server.port, "host", server.host)
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
