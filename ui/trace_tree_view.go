package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type TraceTreeViewModel struct {
	tree helpers.TraceTreeNode

	cursor int
	width  int
}

func CreateTraceTreeViewModel() TraceTreeViewModel {
	return TraceTreeViewModel{}
}

func (m TraceTreeViewModel) Init() tea.Cmd {
	return nil
}

func (m TraceTreeViewModel) Update() (TraceTreeViewModel, tea.Cmd) {
	return m, nil
}

func (m TraceTreeViewModel) View() string {
	row := 0

	lines := make([]string, 0)
	for depth, item := range m.tree.All() {
		lines = append(lines, lipgloss.NewStyle().PaddingLeft(depth).Render(spanView(item, row == m.cursor)))
		row++
	}

	return lipgloss.NewStyle().Width(m.width).Render(lipgloss.JoinVertical(0, lines...))
}
