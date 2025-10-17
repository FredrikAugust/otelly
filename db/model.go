package db

import (
	"database/sql"
	"time"
)

type Span struct {
	TraceID      string         `db:"trace_id"`
	ID           string         `db:"id"`
	Name         string         `db:"name"`
	StartTime    time.Time      `db:"start_time"`
	Duration     time.Duration  `db:"duration_ns"`
	ParentSpanID sql.NullString `db:"parent_span_id"`

	StatusCode    string         `db:"status_code"`
	StatusMessage sql.NullString `db:"status_message"`

	Attributes map[string]any `db:"attributes"`

	ResourceID string `db:"resource_id"`
}

type SpanWithResource struct {
	Span

	ServiceName      string `db:"service_name"`
	ServiceNamespace string `db:"service_namespace"`
}

type Log struct {
	SpanID         sql.NullString `db:"span_id"`
	Body           string         `db:"body"`
	Timestamp      time.Time      `db:"timestamp"`
	SeverityNumber int            `db:"severity_number"`
	SeverityText   string         `db:"severity_text"`
	ResourceID     string         `db:"resource_id"`
	Attributes     map[string]any `db:"attributes"`
}

type Resource struct {
	ID               string `db:"id"`
	ServiceName      string `db:"service_name"`
	ServiceNamespace string `db:"service_namespace"`
}
