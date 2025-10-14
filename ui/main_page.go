package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type MainPageModel struct {
	spanTable   SpanTableModel
	spanDetails SpanDetailsModel

	width  int
	height int
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

func (m MainPageModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.spanTable.View(),
		m.spanDetails.View(),
	)
}

func (m MainPageModel) Update(msg tea.Msg) (MainPageModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mainPanelWidth := 76
		sidePanelWidth := m.width - mainPanelWidth

		m.spanTable.height = m.height
		m.spanTable.width = mainPanelWidth

		m.spanDetails.height = m.height
		m.spanDetails.width = sidePanelWidth
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
