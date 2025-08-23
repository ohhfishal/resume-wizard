package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

type DB struct {
	*Queries
	conn *sql.DB
}

func Open(ctx context.Context, driverName string, dataSourceName string) (*DB, error) {
	if driverName != "sqlite3" {
		return nil, fmt.Errorf("driverName not supported: %s", driverName)
	}

	db, err := sql.Open("sqlite", dataSourceName)
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
