package ui

import (
	"context"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
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
	return listenForNewSpans(m.bus.TraceBus)
}
