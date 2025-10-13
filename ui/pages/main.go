// Package pages contain the main pages of the UI
package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/components"
)

type MainPageModel struct {
	spanTable   *components.SpanTableModel
	spanDetails *components.SpanDetailsModel
}

func CreateMainPageModel(db *db.Database) *MainPageModel {
	return &MainPageModel{
		spanTable:   components.CreateSpanTableModel(db),
		spanDetails: components.CreateSpanDetailsModel(db),
	}
}

func (m MainPageModel) Init() tea.Cmd {
	return tea.Batch(
		m.spanDetails.Init(),
		m.spanTable.Init(),
	)
}

func (m MainPageModel) View(w, h int) string {
	mainPanelWidth := 76
	sidePanelWidth := w - mainPanelWidth - 4

	m.spanTable.SetHeight(h - 1)
	m.spanTable.SetWidth(mainPanelWidth)

	m.spanDetails.SetHeight(h - 1)
	m.spanDetails.SetWidth(sidePanelWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.spanTable.View(),
		m.spanDetails.View(),
	)
}

func (m MainPageModel) Update(msg tea.Msg) (MainPageModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	*m.spanTable, cmd = m.spanTable.Update(msg)
	cmds = append(cmds, cmd)

	*m.spanDetails, cmd = m.spanDetails.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
