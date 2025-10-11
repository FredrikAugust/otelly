package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func Start(ctx context.Context) error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen(), tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
