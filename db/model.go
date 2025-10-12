package db

import "time"

type Span struct {
	TraceID   string    `db:"trace_id"`
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	StartTime time.Time `db:"start_time"`
	// Duration in milliseconds
	Duration time.Duration `db:"duration_ms"`

	ServiceName string `db:"service_name"`
}

type Resource struct {
	ServiceName      string `db:"name"`
	ServiceNamespace string `db:"namespace"`
}
