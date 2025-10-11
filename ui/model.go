// Package ui contains code specific to the rendering and
// state management of the user interface.
package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/ui/components"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	PageMain = iota
)

type windowSize struct {
	height int
	width  int
}
type Service struct {
	Namespace string
	Name      string
}

// InsertResourceSpans takes a resourceSpans object, and returns commands to insert
// the root spans into the table.
func (m *Model) InsertResourceSpans(resourceSpans ptrace.ResourceSpans) []tea.Cmd {
	res := resourceSpans.Resource()

	name, exists := res.Attributes().Get(string(semconv.ServiceNameKey))
	if !exists {
		name = pcommon.NewValueStr("unknown")
	}

	ns, exists := res.Attributes().Get(string(semconv.ServiceNamespaceKey))
	if !exists {
		ns = pcommon.NewValueStr("unknown")
	}

	svc := Service{
		Name:      name.Str(),
		Namespace: ns.Str(),
	}

	cmds := make([]tea.Cmd, 0)

	for _, scopeSpans := range resourceSpans.ScopeSpans().All() {
		for _, span := range scopeSpans.Spans().All() {
			m.spanIDToService[span.SpanID()] = svc
			if span.ParentSpanID().IsEmpty() {
				cmds = append(cmds, func() tea.Msg {
					return messageNewRootSpan{
						serviceName: svc.Name,
						name:        span.Name(),
						startTime:   span.StartTimestamp().AsTime(),
						duration:    span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()),
					}
				})
			}
		}
		scopeSpans.Spans().MoveAndAppendTo(m.spans)
	}

	return cmds
}

type Model struct {
	currentPage int

	bus *bus.TransportBus

	spans           ptrace.SpanSlice
	spanIDToService map[pcommon.SpanID]Service

	spanTable table.Model

	windowSize *windowSize
}

func NewModel(bus *bus.TransportBus) *Model {
	return &Model{
		currentPage: PageMain,
		bus:         bus,
		spanTable:   components.CreateSpanTable(),
		windowSize: &windowSize{
			0, 0,
		},
		spans:           ptrace.NewSpanSlice(),
		spanIDToService: make(map[pcommon.SpanID]Service),
	}
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return messageResourceSpansArrived{resourceSpans: <-spanChan}
	}
}
