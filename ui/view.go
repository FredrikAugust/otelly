package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	baseStyles := lipgloss.NewStyle().
		Width(m.windowSize.width).
		Height(m.windowSize.height)

	switch m.currentPage {
	case PageMain:
		sidePanelWidth := 80
		mainPanelWidth := baseStyles.GetWidth() - sidePanelWidth

		m.spanTable.SetHeight(baseStyles.GetHeight() - 1)
		m.spanTable.SetWidth(mainPanelWidth)

		m.spanDetails.SetHeight(baseStyles.GetHeight() - 1)
		m.spanDetails.SetWidth(sidePanelWidth)

		return baseStyles.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Padding(0, 1).Width(baseStyles.GetWidth()).Background(lipgloss.Color("32")).Bold(true).Render("Otelly"),
				lipgloss.JoinHorizontal(lipgloss.Top,
					m.spanTable.View(),
					m.spanDetails.View(),
				),
			),
		)
	default:
		return "unknown view"
	}
}
