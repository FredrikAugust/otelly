// Package ui contains code specific to the rendering and
// state management of the user interface.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/components"
	"go.opentelemetry.io/collector/pdata/ptrace"
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

type Model struct {
	currentPage int

	bus *bus.TransportBus

	db *db.Database

	spanTable   *components.SpanTableModel
	spanDetails *components.SpanDetailsModel

	windowSize *windowSize
}

func NewModel(bus *bus.TransportBus, db *db.Database) *Model {
	return &Model{
		currentPage: PageMain,
		bus:         bus,
		spanTable:   components.CreateSpanTableModel(db),
		spanDetails: components.CreateSpanDetailsModel(db),
		db:          db,
		windowSize: &windowSize{
			0, 0,
		},
	}
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return messageResourceSpansArrived{resourceSpans: <-spanChan}
	}
}
