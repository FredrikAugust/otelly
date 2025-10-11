package ui

import (
	"time"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type messageResourceSpansArrived struct {
	resourceSpans ptrace.ResourceSpans
}

type messageNewRootSpan struct {
	serviceName string
	name        string
	duration    time.Duration
	startTime   time.Time
}
