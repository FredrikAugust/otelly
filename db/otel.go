package db

import (
	"fmt"
	"log/slog"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// InsertResourceSpans inserts the resource and all encompassing spans
// into the database.
//
// Needs to be optimised to use batch insert. Now it runs N queries where
// N is |spans|.
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

	_, err := d.sqlDB.Exec(`INSERT OR IGNORE INTO resource VALUES ($1, $2, $3)`,
		resID,
		resName.Str(),
		resNamespace.Str(),
	)
	if err != nil {
		slog.Warn("couldn't insert resource", "resourceID", resID)
		return err
	}

	for _, scopeSpans := range spans.ScopeSpans().All() {
		for _, span := range scopeSpans.Spans().All() {
			_, err := d.sqlDB.Exec(
				`INSERT INTO span VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				span.SpanID().String(),
				span.Name(),
				span.StartTimestamp().AsTime(),
				span.EndTimestamp().AsTime(),
				span.TraceID().String(),
				span.ParentSpanID().String(),
				resID,
			)
			if err != nil {
				slog.Warn("failed to insert span", "spanID", span.SpanID().String(), "error", err)
				return err
			}
		}
	}

	return nil
}

func (d *Database) GetSpan(id string) (*Span, error) {
	var span Span

	err := d.sqlDB.Get(&span, `SELECT id, name FROM span WHERE id = $1`, id)
	if err != nil {
		slog.Warn("failed to get span", "spanID", id)
		return nil, err
	}

	slog.Debug("got span", "id", id)

	return &span, nil
}

func (d *Database) GetSpans() []Span {
	spans := make([]Span, 0)
	err := d.sqlDB.Select(&spans, `SELECT id, name FROM span`)
	if err != nil {
		return spans
	}

	slog.Debug("got spans", "len", len(spans))

	return spans
}
