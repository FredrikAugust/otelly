package ui

import (
	"context"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/components"
)

func Start(ctx context.Context, bus *bus.TransportBus, db *db.Database) error {
	slog.Info("initializing and running UI")

	p := tea.NewProgram(
		NewModel(bus, db),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.mainPageModel.Init(),
		listenForNewSpans(m.bus.TraceBus),
		func() tea.Msg {
			// In case we're runnign DB with file we can keep the seed around, and
			// this will make it load into the UI immediately.
			return components.MessageUpdateRootSpanRows{}
		},
	)
}
