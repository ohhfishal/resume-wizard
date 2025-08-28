package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

//go:embed schema.sql
var schema string

type Config struct {
	Source      string `short:"s" default:":memory:" env:"SOURCE" help:"Database connection string."`
	Driver      string `short:"d" enum:"sqlite" default:"sqlite" help:"SQL driver to use."`
	UseTempFile bool   `help:"Enable to use a copy of the source file (sqlite)."`
}

type DB struct {
	*Queries
	conn *sql.DB
}

func (config *Config) Open(ctx context.Context) (*DB, error) {
	switch config.Driver {
	case "sqlite":
		return config.openSQLite(ctx)
	default:
		return nil, fmt.Errorf("driver not supported: %s", config.Driver)
	}
}

func (q *DB) BeginTx(ctx context.Context) (*sql.Tx, *DB, error) {
	tx, err := q.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("beginning transaction: %w", err)
	}

	// Keep the original connection reference in the transaction Queries too
	txQueries := &DB{
		Queries: q.WithTx(tx),
		conn:    q.conn,
	}

	return tx, txQueries, nil
}

func (config *Config) openSQLite(ctx context.Context) (*DB, error) {
	if config.UseTempFile {
		if err := config.SetSourceToTemp(ctx); err != nil {
			return nil, fmt.Errorf("switching to temp file: %w", err)
		}
	}

	db, err := sql.Open("sqlite", config.Source)
	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("running migration: %w", err)
	}

	return &DB{
		Queries: New(db),
		conn:    db,
	}, nil
}

func (config *Config) SetSourceToTemp(ctx context.Context) error {
	original := strings.Split(config.Source, ";")

	source, err := os.Open(original[0])
	if err != nil {
		return fmt.Errorf("opening original file: %w", err)
	}
	defer source.Close()

	temp, err := os.CreateTemp("", "temp-wizard-*.db")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer temp.Close()

	written, err := io.Copy(temp, source)
	if err != nil {
		return fmt.Errorf("copying content to temp file: %w", err)
	}

	config.Source = strings.Join(append(
		[]string{temp.Name()}, original[1:]...,
	), ";")
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		logger.Info("using temp file for sqlite",
			slog.String("file", temp.Name()),
			slog.String("new_source", config.Source),
			slog.Int64("size", written),
		)
	}

	return nil

}
