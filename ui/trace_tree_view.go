package ui

import (
	"fmt"
	"time"

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
		lines = append(lines, lipgloss.NewStyle().PaddingLeft(depth).MaxWidth(m.width).Render(spanView(item, row == m.cursor)))
		row++
	}

	return lipgloss.NewStyle().Width(m.width).Render(lipgloss.JoinVertical(0, lines...))
}

func spanView(item helpers.TraceTreeNodeItem, selected bool) string {
	style := lipgloss.
		NewStyle()

	if selected {
		style = style.Background(ColorAccent)
	}

	secondaryText := item.Span.ServiceName + " • " + item.Span.Duration.Round(time.Millisecond).String()
	if item.DurationOfParent != 1 {
		pctOfParentSpan := item.DurationOfParent * 100
		secondaryText += fmt.Sprintf(" (%.1f%%)", pctOfParentSpan)
	}

	return style.Render(item.Span.Name) + " " + TextTertiary.Render(secondaryText)
}
