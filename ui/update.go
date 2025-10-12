package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/ui/components"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, 0)
	)

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
		cmds = append(
			cmds,
			listenForNewSpans(m.bus.TraceBus),
			m.InsertResourceSpans(msg.resourceSpans),
		)
	}

	*m.spanTable, cmd = m.spanTable.Update(msg)
	cmds = append(cmds, cmd)

	*m.spanDetails, cmd = m.spanDetails.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) InsertResourceSpans(resourceSpans ptrace.ResourceSpans) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	m.db.InsertResourceSpans(resourceSpans)

	*m.spanTable, cmd = m.spanTable.Update(components.MessageUpdateRootSpanRows{})
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}
