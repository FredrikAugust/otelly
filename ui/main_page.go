package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type MainPageModel struct {
	spanTable   SpanTableModel
	spanDetails SpanDetailsModel
}

func CreateMainPageModel(db *db.Database) MainPageModel {
	return MainPageModel{
		spanTable:   CreateSpanTableModel(db),
		spanDetails: CreateSpanDetailsModel(db),
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
	sidePanelWidth := w - mainPanelWidth

	m.spanTable.SetHeight(h)
	m.spanTable.SetWidth(mainPanelWidth)

	m.spanDetails.SetHeight(h)
	m.spanDetails.SetWidth(sidePanelWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.spanTable.View(),
		m.spanDetails.View(),
	)
}

func (m MainPageModel) Update(msg tea.Msg) (MainPageModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.spanTable, cmd = m.spanTable.Update(msg)
	cmds = append(cmds, cmd)

	m.spanDetails, cmd = m.spanDetails.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
