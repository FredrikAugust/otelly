package ui

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	baseStyles := lipgloss.NewStyle().
		Width(m.windowSize.width).
		Height(m.windowSize.height)

	switch m.currentPage {
	case PageMain:
		mainPanelWidth := int(float32(baseStyles.GetWidth()) * 0.6)
		sidePanelWidth := int(float32(baseStyles.GetWidth()) * 0.4)

		m.spanTable.SetHeight(baseStyles.GetHeight())
		m.spanTable.SetWidth(mainPanelWidth)

		m.spanDetails.SetHeight(baseStyles.GetHeight())
		m.spanDetails.SetWidth(sidePanelWidth)

		return baseStyles.Render(
			lipgloss.JoinHorizontal(lipgloss.Top,
				m.spanTable.View(),
				m.spanDetails.View(),
			),
		)
	default:
		return "unknown view"
	}
}
