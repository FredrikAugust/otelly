package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/ui/styling"
)

func (m Model) View() string {
	baseStyles := lipgloss.NewStyle().
		Width(m.windowSize.width).
		Height(m.windowSize.height)

	return baseStyles.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Padding(0, 1).Width(baseStyles.GetWidth()).Background(styling.ColorAccent).Bold(true).Render("Otelly"),
			m.GetPageView(baseStyles.GetWidth(), baseStyles.GetHeight()-1),
		),
	)
}

func (m *Model) GetPageView(w, h int) string {
	switch m.currentPage {
	case PageMain:
		return m.mainPageModel.View(w, h)
	case PageTrace:
		return m.tracePageModel.View(w, h)
	}

	return "unknown page"
}
