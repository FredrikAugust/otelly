package ui

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type SpansPageModel struct {
	spans []db.Span

	width  int
	height int
}

func NewSpansPageModel(spans []db.Span) SpansPageModel {
	return SpansPageModel{spans: spans}
}

func (m SpansPageModel) Init() tea.Cmd {
	return nil
}

func (m SpansPageModel) Update(msg tea.Msg) (SpansPageModel, tea.Cmd) {
	return m, nil
}

func (m SpansPageModel) View() string {
	return helpers.HStack(m.TableView(), m.DetailView())
}

func (m SpansPageModel) TableView() string {
	// width  = 2/3 of window width
	// height = viewport height -  border
	container := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(helpers.ColorBorder).Width(int(math.Floor(float64(m.width)*2.0/3.0 - 2.0))).Height(m.height - 2)

	return container.Render("table")
}

func (m SpansPageModel) DetailView() string {
	// width  = 1/3 of window width
	// height = viewport height -  border
	container := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(helpers.ColorBorder).Width(int(math.Ceil(float64(m.width)*1.0/3.0 - 2.0))).Height(m.height - 2)

	return container.Render("detail")
}
