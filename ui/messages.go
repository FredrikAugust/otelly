package ui

import (
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/flamegraph"
)

type (
	MsgSpanPageUpdateTable struct{}
	MsgNewSpans            struct{ spans []db.Span }
	MsgNewLogs             struct{ logs []db.Log }

	MsgLoadTrace   struct{ traceID string }
	MsgTreeUpdated struct{ tree flamegraph.Node }
)
