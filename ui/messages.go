package ui

import (
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type (
	MessageGoToTrace    struct{ TraceID string }
	MessageGoToMainPage struct{}
	MessageGoToLogs     struct{}

	// MessageResourceSpansArrived signifies that the collector has received new spans
	MessageResourceSpansArrived struct{ ResourceSpans ptrace.ResourceSpans }
	MessageResourceLogsArrived  struct{ ResourceLogs plog.ResourceLogs }

	MessageSetSelectedSpan struct{ Span db.SpanWithResource }

	MessageUpdateRootSpans struct{ NewRootSpans []db.SpanWithResource }
	MessageUpdateLogs      struct{ NewLogs []db.Log }

	MessageResourceReceived            struct{ Resource db.Resource }
	MessageResourceAggregationReceived struct {
		Aggregation []db.SpansPerMinuteForServiceModel
	}

	MessageReceivedTraceSpans struct{ Spans []db.SpanWithResource }
)
