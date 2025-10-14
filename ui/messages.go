package ui

import (
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type (
	MessageGoToTrace        struct{ TraceID string }
	MessageReturnToMainPage struct{}

	// MessageResourceSpansArrived signifies that the collector has received new spans
	MessageResourceSpansArrived struct{ resourceSpans ptrace.ResourceSpans }

	MessageSetSelectedSpan    struct{ SpanID string }
	MessageResetDetail        struct{}
	MessageUpdateRootSpanRows struct{}
)
