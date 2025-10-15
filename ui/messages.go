package ui

import (
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type (
	MessageGoToTrace        struct{ TraceID string }
	MessageReturnToMainPage struct{}

	// MessageResourceSpansArrived signifies that the collector has received new spans
	MessageResourceSpansArrived struct{ ResourceSpans ptrace.ResourceSpans }

	MessageSetSelectedSpan struct{ SpanID string }
	MessageResetDetail     struct{}
	MessageUpdateRootSpans struct{ NewRootSpans []db.Span }
)
