package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize.height = msg.Height
		m.windowSize.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case messageResourceSpansArrived:
		return m, tea.Batch(
			listenForNewSpans(m.bus.TraceBus),
			tea.Batch(m.InsertResourceSpans(msg.resourceSpans)...),
		)
	case messageNewRootSpan:
		m.spanTable.SetRows(append(m.spanTable.Rows(), table.Row{
			msg.name,
			msg.serviceName,
			msg.startTime.Format("15:04:05.000"),
			msg.duration.String(),
		}))
		return m, nil
	}

	m.spanTable, cmd = m.spanTable.Update(msg)

	return m, cmd
}
