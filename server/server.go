package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ohhfishal/resume-wizard/db"
)

type Config struct {
	Port           string        `default:"8080" short:"P" help:"Port to serve on"`
	Host           string        `default:"localhost" short:"H" help:"Address to serve from"`
	RequestTimeout time.Duration `default:"30s" help:"How long to keep requests alive"`
	Database       db.Config     `embed:"" prefix:"database-" envprefix:"DATABASE_"`
}

type Server struct {
	logger   *slog.Logger
	database *db.DB
	config   Config
}

func New(ctx context.Context, config Config, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	if config.Database.Source == ":memory:" {
		logger.Warn("using in-memory database")
	}

	database, err := config.Database.Open(context.WithValue(ctx, "logger", logger))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return &Server{
		database: database,
		logger:   logger,
		config:   config,
	}, nil
}

func (server *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()

	handler := loggingMiddleware(server.logger)(
		timeoutMiddleware(server.config.RequestTimeout)(
			recoveryMiddleware()(mux),
		),
	)

	mux.Handle("GET /health", HandleHealth())
	mux.Handle("POST /api/match", Chain(
		WithBearerAuth(nil),
		HandlePostMatch(),
	))

	mux.Handle("/", HandleNotFound())

	s := &http.Server{
		Addr:         net.JoinHostPort(server.config.Host, server.config.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		server.logger.Info("shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			server.logger.Error("closing server",
				slog.Any("error", err),
			)
		}
	}()

	server.logger.Info(
		"starting server",
		slog.String("port", server.config.Port),
		slog.String("host", server.config.Host),
	)
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

func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "timeout")
	}
}

func recoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}
			}()
			next.ServeHTTP(w, r)
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
