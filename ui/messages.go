package ui

import (
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type (
	MessageGoToTrace        struct{ TraceID string }
	MessageReturnToMainPage struct{}

	// MessageResourceSpansArrived signifies that the collector has received new spans
	MessageResourceSpansArrived  struct{ ResourceSpans ptrace.ResourceSpans }
	MessageResourceSpansInserted struct{}

	MessageSetSelectedSpan struct{ Span db.SpanWithResource }
	MessageUpdateRootSpans struct{ NewRootSpans []db.SpanWithResource }

	MessageResourceReceived            struct{ Resource db.Resource }
	MessageResourceAggregationReceived struct {
		Aggregation []db.SpansPerMinuteForServiceModel
	}

	MessageReceivedTraceSpans struct{ Spans []db.SpanWithResource }
)
