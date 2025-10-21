// Package bus contains the message bus used to pass new data
// from the telemetry part of the app to the UI.
package bus

import (
	"github.com/fredrikaugust/otelly/db"
)

type TransportBus struct {
	SpanBus chan []db.Span
	LogBus  chan []db.Log
}

func NewTransportBus() *TransportBus {
	return &TransportBus{
		SpanBus: make(chan []db.Span),
		LogBus:  make(chan []db.Log),
	}
}
