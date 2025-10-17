package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

func (d *Database) GetResource(id string) (*Resource, error) {
	var res Resource

	err := d.sqlDB.Get(
		&res,
		`
		SELECT
			*
		FROM
			resource
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *Database) GetRootSpans() []SpanWithResource {
	spans := make([]SpanWithResource, 0)
	err := d.sqlDB.Select(
		&spans,
		`
		SELECT
			s.trace_id,
			s.id,
			s.name,
			s.start_time,
			s.duration_ns,
			s.status_code,
			s.attributes,
			s.resource_id,
			r.service_name
		FROM
			span s
		LEFT JOIN resource r ON s.resource_id = r.id
		WHERE
			s.parent_span_id IS NULL
		ORDER BY
			s.start_time DESC`,
	)
	if err != nil {
		return spans
	}

	return spans
}

type SpansPerMinuteForServiceModel struct {
	Timestamp time.Time `db:"bucket_start"`
	SpanCount int       `db:"span_count"`
}

func (d *Database) SpansPerMinuteForService(resourceID string) ([]SpansPerMinuteForServiceModel, error) {
	query := `
	SELECT
		date_trunc('minute', start_time) as bucket_start,
		COUNT(*) as span_count
	FROM
		span s
	LEFT JOIN resource r ON
		s.resource_id = r.id
	WHERE
		r.id = $1
	GROUP BY
		bucket_start
	ORDER BY
		bucket_start DESC
	`

	res := make([]SpansPerMinuteForServiceModel, 0)

	err := d.sqlDB.Select(&res, query, resourceID)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("got span history", zap.Int("num", len(res)))

	return res, nil
}

func (d *Database) GetSpansForTrace(traceID string) ([]SpanWithResource, error) {
	query := `
	SELECT
		s.id,
		s.start_time,
		s.trace_id,
		s.name,
		s.parent_span_id,
		s.duration_ns,
		s.attributes,
		r.service_name
	FROM
		span s
	LEFT JOIN resource r ON
		s.resource_id = r.id
	WHERE
		s.trace_id = $1
	ORDER BY
		s.start_time`

	res := make([]SpanWithResource, 0)

	err := d.sqlDB.Select(&res, query, traceID)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("got spans for trace", zap.String("traceID", traceID), zap.Int("spans", len(res)))

	return res, nil
}
