package db

import "time"

type Span struct {
	TraceID       string        `db:"trace_id"`
	ID            string        `db:"id"`
	Name          string        `db:"name"`
	StartDuration time.Time     `db:"start_duration"`
	Duration      time.Duration `db:"duration"`
}

type Resource struct {
	Name      string `db:"name"`
	Namespace string `db:"namespace"`
}
