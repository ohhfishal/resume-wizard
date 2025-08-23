package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ohhfishal/resume-wizard/assets"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/templates"
	"github.com/ohhfishal/resume-wizard/templates/page"
	"github.com/ohhfishal/resume-wizard/wizard"
)

type Config struct {
	DatabaseSource string `short:"s" default:":memory:" env:"DATABASE_SOURCE" help:"Database connection string (sqlite)."`
	Port           string `default:"8080" help:"Port to serve on"`
	Host           string `default:"localhost" help:"Address to serve from"`
}

type Server struct {
	logger   *slog.Logger
	database *db.DB
	config   Config
	wizard   *wizard.Wizard
}

func New(ctx context.Context, config Config, logger *slog.Logger) (*Server, error) {
	if config.DatabaseSource == ":memory:" {
		logger.Warn("using in-memory database")
	}

	database, err := db.Open(ctx, "sqlite3", config.DatabaseSource)
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
	r := chi.NewRouter()

	r.Use(loggingMiddleware(server.logger))
	r.Use(middleware.Recoverer)

	r.Post("/api/dev/application/{session_id}", PostApplicationHandlerNew(server.logger, server.database))

	r.Post("/api/dev/base", PostBaseResumeHandler(server.logger, server.database))
	r.Post("/base/upload", GetBaseResumeForm(server.logger, server.database))
	r.Get("/base/new", GetBaseResumeForm(server.logger, server.database))

	r.Post("/api/dev/generate", GenerateHandler(server.logger, server.database, server.wizard))

	// TODO: Remove old endpoints
	r.Put("/resume/{id}", PutResumeHandler(server.logger, server.database))
	r.Post("/resume", PostResumeHandler(server.logger, server.database))
	r.Post("/compile/resume", PostCompileHandler(server.logger, server.database))
	r.Post("/application", PostApplicationHandler(server.logger, server.database))
	r.Put("/application", PutApplicationHandler(server.logger, server.database))

	r.Mount(
		"/assets",
		http.StripPrefix("/assets", http.FileServer(http.FS(assets.Assets))),
	)

	r.Route("/components", ComponentsHandler(server.logger, server.database))

	r.Get("/base", func(w http.ResponseWriter, r *http.Request) {
		page.BaseResume(page.BaseResumeProps{}).Render(r.Context(), w)
	})
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		page.Login(page.LoginProps{}).Render(r.Context(), w)
	})
	r.Get("/tailor", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("invalid base resume id: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		base, err := server.database.GetBaseResume(r.Context(), db.GetBaseResumeParams{
			UserID: 0, /* TODO: Set to userID */
			ID:     id,
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for base resume: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		page.TailorResume(page.TailorResumeProps{
			Base: base,
		}).Render(r.Context(), w)
	})
	r.Get("/tailor/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		session, err := server.database.GetSession(r.Context(), db.GetSessionParams{
			UserID: 0, /* TODO: Set to userID */
			Uuid:   r.PathValue("uuid"),
		})
		if err != nil {
			http.Error(w,
				fmt.Sprintf("restoring session: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		base, err := server.database.GetBaseResume(r.Context(), db.GetBaseResumeParams{
			UserID: 0, /* TODO: Set to userID */
			ID:     session.BaseResumeID,
		})

		page.TailorResume(page.TailorResumeProps{
			Base:            base,
			Session:         session,
			LockApplication: true,
		}).Render(r.Context(), w)
	})
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		resumes, err := server.database.GetBaseResumes(r.Context(), 0 /* TODO: Set to userID */)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for names: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		applications, err := server.database.GetApplicationsV2(r.Context(), 0 /* TODO: Set to userID */)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("reading database for applications: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		if len(applications) == 0 {
			server.logger.Warn("NO APPLICATIONS")
		}

		page.Home(page.HomeProps{
			Resumes:      resumes,
			Applications: applications,
		}).Render(r.Context(), w)
	})

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
		templates.MainPage(resumes, applications).Render(r.Context(), w)
	})

	r.NotFound(NotFoundHandler)

	s := &http.Server{
		Addr:         net.JoinHostPort(server.config.Host, server.config.Port),
		Handler:      r,
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

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	err := errors.New("Page Not Found")

	w.WriteHeader(http.StatusNotFound)
	switch {
	case strings.Contains(accept, "text/html"):
		w.Header().Set("Content-Type", "text/html")
		page.Error(err).Render(r.Context(), w)
	case strings.Contains(accept, "text/plain"):
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
	case strings.Contains(accept, "application/json"):
		fallthrough
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"error": err})
	}
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
