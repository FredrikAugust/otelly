// Package bus contains the message bus used to pass new data
// from the telemetry part of the app to the UI.
package bus

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type TransportBus struct {
	TraceBus chan ptrace.ResourceSpans
	LogBus   chan plog.ResourceLogs
}

func NewTransportBus() *TransportBus {
	return &TransportBus{
		TraceBus: make(chan ptrace.ResourceSpans),
		LogBus:   make(chan plog.ResourceLogs),
	}
}
