package ui

import (
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type messageResourceSpansArrived struct {
	resourceSpans ptrace.ResourceSpans
}
