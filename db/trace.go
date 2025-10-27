package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// InsertResourceSpans inserts the resource and all encompassing spans
// into the database.
func (d *Database) InsertResourceSpans(ctx context.Context, spans ptrace.ResourceSpans) error {
	resID, err := d.InsertResource(ctx, spans.Resource())
	if err != nil {
		return err
	}

	tx, err := d.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, scopeSpans := range spans.ScopeSpans().All() {
		for _, span := range scopeSpans.Spans().All() {
			attrs, err := json.Marshal(span.Attributes().AsRaw())
			if err != nil {
				zap.L().Warn("could not serialize span attributes to JSON", zap.Error(err))
				attrs = []byte("{}")
			}

			zap.L().Debug("inserting new span", zap.Bool("root", span.ParentSpanID().IsEmpty()), zap.String("name", span.Name()))

			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO span VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				span.SpanID().String(),
				span.Name(),
				span.StartTimestamp().AsTime(),
				span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).Nanoseconds(),
				span.TraceID().String(),
				span.Kind().String(),
				sql.NullString{String: span.ParentSpanID().String(), Valid: !span.ParentSpanID().IsEmpty()},
				span.Status().Code().String(),
				sql.NullString{String: span.Status().Message(), Valid: span.Status().Message() != ""},
				attrs,
				resID,
			)
			if err != nil {
				zap.L().Warn("failed to create span", zap.String("name", span.Name()), zap.String("resourceID", resID))
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction trace: %w", err)
	}

	return nil
}

func (d *Database) ClearSpans(ctx context.Context) error {
	_, err := d.sqlDB.ExecContext(ctx, `TRUNCATE TABLE span`)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetSpans(ctx context.Context) ([]Span, error) {
	spans := make([]Span, 0)
	err := d.sqlDB.SelectContext(
		ctx,
		&spans,
		`
		SELECT
			*
		FROM
			span
		ORDER BY
			start_time DESC`,
	)
	if err != nil {
		return spans, err
	}

	return spans, nil
}

func (d *Database) GetSpansForTrace(ctx context.Context, traceID string) ([]Span, error) {
	spans := make([]Span, 0)
	err := d.sqlDB.SelectContext(
		ctx,
		&spans,
		`
		SELECT
			*
		FROM
			span
		WHERE
			trace_id = ?
		ORDER BY
			start_time DESC`,
		traceID,
	)
	if err != nil {
		return spans, err
	}

	return spans, nil
}

func FilterRootSpans(spans []Span) []Span {
	rootSpans := make([]Span, 0)
	for _, span := range spans {
		if !span.ParentSpanID.Valid {
			rootSpans = append(rootSpans, span)
		}
	}

	return rootSpans
}
