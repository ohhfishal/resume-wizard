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

func Open(ctx context.Context, driverName string, dataSourceName string) (*Queries, error) {
	if driverName != "sqlite3" {
		return nil, fmt.Errorf("driverName not supported: %s", driverName)
	}

	// TODO: Fix hack
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("running migration: %w", err)
	}
	return New(db), nil
}
