package db

import "time"

type Span struct {
	TraceID   string    `db:"trace_id"`
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	StartTime time.Time `db:"start_time"`
	// Duration in milliseconds
	Duration time.Duration `db:"duration_ms"`

	StatusCode    string `db:"status_code"`
	StatusMessage string `db:"status_message"`

	Attributes map[string]any `db:"attributes"`

	ResourceID string `db:"resource_id"`

	ServiceName string `db:"service_name"`
}

type Resource struct {
	ID               string `db:"id"`
	ServiceName      string `db:"service_name"`
	ServiceNamespace string `db:"service_namespace"`
}
