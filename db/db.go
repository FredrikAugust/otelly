// Package db contains implementation of a database interface using DuckDB.
package db

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type Database struct {
	sqlDB *sqlx.DB

	// resourceLock will be held by the resource writing to the resource
	// table as it will be contested when the collector flushes new spans+logs
	// with the same _new_ resource. sadly the on conflict ignore isn't fast
	// enough to catch the conflict.
	resourceLock sync.Mutex
}

func NewDB() (*Database, error) {
	// We always want to start with a clean slate
	// dbFile, err := os.Stat("./local.db")
	// if err == nil {
	// 	os.Remove(dbFile.Name())
	// }

	db, err := sqlx.Open("duckdb", "./local.db")
	if err != nil {
		return nil, err
	}

	return &Database{
		sqlDB: db,
	}, err
}

func (d *Database) Close() error {
	slog.Info("closing db")
	return d.sqlDB.Close()
}

func (d *Database) Migrate(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS resource (id VARCHAR PRIMARY KEY, service_name VARCHAR, service_namespace VARCHAR)`,
		`CREATE TABLE IF NOT EXISTS span (
			id VARCHAR PRIMARY KEY,
			name VARCHAR,
			start_time TIMESTAMP,
			duration_ns INTEGER,
			trace_id VARCHAR,
			parent_span_id VARCHAR,
			status_code VARCHAR,
			status_message VARCHAR,
			attributes JSON,
			resource_id VARCHAR,
			FOREIGN KEY (resource_id) REFERENCES resource (id)
		)`,
		`CREATE TABLE IF NOT EXISTS log (
			span_id VARCHAR,
			body VARCHAR,
			timestamp TIMESTAMP,
			severity_number INTEGER,
			severity_text VARCHAR,
			resource_id VARCHAR,
			attributes JSON,
			FOREIGN KEY (resource_id) REFERENCES resource (id)
		)`,
	}

	for _, migration := range migrations {
		_, err := d.sqlDB.ExecContext(ctx, migration)
		if err != nil {
			slog.Error("failed to apply migration", "migration", migration)
			return err
		}
	}
	slog.Info("finished migrating DB", "numMigrations", len(migrations))

	return nil
}

func (d *Database) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.sqlDB.BeginTx(ctx, &sql.TxOptions{})
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.sqlDB.ExecContext(ctx, query, args...)
}
