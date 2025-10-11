// Package ui contains code specific to the rendering and
// state management of the user interface.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/ui/components"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
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
func (m *Model) InsertResourceSpans(resourceSpans ptrace.ResourceSpans) tea.Cmd {
	res := resourceSpans.Resource()

	cmds := make([]tea.Cmd, 0)

	for _, scopeSpans := range resourceSpans.ScopeSpans().All() {
		for _, span := range scopeSpans.Spans().All() {
			m.spanIDToResource[span.SpanID()] = res
			m.spanIDToSpan[span.SpanID()] = span

			if span.ParentSpanID().IsEmpty() {
				resName, exists := res.Attributes().Get(string(semconv.ServiceNameKey))
				if !exists {
					resName = pcommon.NewValueStr("unknown")
				}
				cmds = append(cmds, func() tea.Msg {
					return components.MessageNewRootSpan{
						Span:         &span,
						ResourceName: resName.Str(),
					}
				})
			}
		}
		scopeSpans.Spans().MoveAndAppendTo(m.spans)
	}

	return tea.Batch(cmds...)
}

type Model struct {
	currentPage int

	bus *bus.TransportBus

	spans ptrace.SpanSlice

	spanIDToSpan     map[pcommon.SpanID]ptrace.Span
	spanIDToResource map[pcommon.SpanID]pcommon.Resource

	spanTable   *components.SpanTableModel
	spanDetails *components.SpanDetailsModel

	windowSize *windowSize
}

func NewModel(bus *bus.TransportBus) *Model {
	return &Model{
		currentPage: PageMain,
		bus:         bus,
		spanTable:   components.CreateSpanTableModel(),
		spanDetails: components.CreateSpanDetailsModel(),
		windowSize: &windowSize{
			0, 0,
		},
		spans:            ptrace.NewSpanSlice(),
		spanIDToSpan:     make(map[pcommon.SpanID]ptrace.Span),
		spanIDToResource: make(map[pcommon.SpanID]pcommon.Resource),
	}
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return messageResourceSpansArrived{resourceSpans: <-spanChan}
	}
}
