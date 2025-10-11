package ui

import "github.com/charmbracelet/lipgloss"

var baseStyles = lipgloss.
	NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() string {
	baseStyles = baseStyles.Width(m.windowSize.width - baseStyles.GetHorizontalFrameSize()).Height(m.windowSize.height - baseStyles.GetVerticalFrameSize())

	switch m.currentPage {
	case PageMain:
		m.spanTable.SetHeight(baseStyles.GetHeight() - 1)
		m.spanTable.SetWidth(baseStyles.GetWidth())

		return baseStyles.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.spanTable.View(),
				m.spanTable.HelpView(),
			),
		)
	default:
		return "unknown view"
	}
}
