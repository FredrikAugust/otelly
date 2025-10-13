// Package ui contains code specific to the rendering and
// state management of the user interface.
package ui

import (
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/pages"
)

const (
	PageMain = iota
	PageTrace
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

	mainPageModel  *pages.MainPageModel
	tracePageModel *pages.TracePageModel

	bus *bus.TransportBus
	db  *db.Database

	windowSize *windowSize
}

func NewModel(bus *bus.TransportBus, db *db.Database) *Model {
	return &Model{
		currentPage: PageMain,
		bus:         bus,
		db:          db,
		windowSize: &windowSize{
			0, 0,
		},
		mainPageModel:  pages.CreateMainPageModel(db),
		tracePageModel: pages.CreateTracePageModel(db),
	}
}
