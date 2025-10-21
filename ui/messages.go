package ui

import "github.com/fredrikaugust/otelly/db"

type (
	MsgSpanPageUpdateTable struct{}
	MsgNewSpans            struct{ spans []db.Span }
	MsgNewLogs             struct{ logs []db.Log }
)
