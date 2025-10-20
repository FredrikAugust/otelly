package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
)

// InsertResourceSpans inserts the resource and all encompassing spans
// into the database.
func (d *Database) InsertResourceSpans(spans ptrace.ResourceSpans) error {
	resName, exists := spans.Resource().Attributes().Get(string(semconv.ServiceNameKey))
	if !exists {
		resName = pcommon.NewValueStr("unknown")
	}
	resNamespace, exists := spans.Resource().Attributes().Get(string(semconv.ServiceNamespaceKey))
	if !exists {
		resNamespace = pcommon.NewValueStr("unknown")
	}
	resID := fmt.Sprintf("%s:%s", resName.Str(), resNamespace.Str())

	tx, err := d.BeginTx(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT OR IGNORE INTO resource VALUES ($1, $2, $3)`,
		resID,
		resName.Str(),
		resNamespace.Str(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, scopeSpans := range spans.ScopeSpans().All() {
		for _, span := range scopeSpans.Spans().All() {
			attrs, err := json.Marshal(span.Attributes().AsRaw())
			if err != nil {
				attrs = []byte("{}")
			}

			zap.L().Debug("inserting new span", zap.Bool("root", span.ParentSpanID().IsEmpty()), zap.String("name", span.Name()))

			_, err = tx.Exec(
				`INSERT INTO span VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				span.SpanID().String(),
				span.Name(),
				span.StartTimestamp().AsTime(),
				span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Nanoseconds(),
				span.TraceID().String(),
				sql.NullString{String: span.ParentSpanID().String(), Valid: !span.ParentSpanID().IsEmpty()},
				span.Status().Code().String(),
				sql.NullString{String: span.Status().Code().String(), Valid: span.Status().Message() != ""},
				attrs,
				resID,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *Database) ClearSpans() error {
	_, err := d.sqlDB.Exec(`TRUNCATE TABLE span`)
	if err != nil {
		return err
	}

	return nil
}
