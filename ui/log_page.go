package ui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/db"
)

type LogPageModel struct {
	logs []db.Log
}

func CreateLogPageModel() LogPageModel {
	return LogPageModel{
		logs: make([]db.Log, 0),
	}
}

func (m LogPageModel) Init() tea.Cmd {
	return nil
}

func (m LogPageModel) View() string {
	return "log " + strconv.Itoa(len(m.logs))
}

func (m LogPageModel) Update(msg tea.Msg) (LogPageModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			cmds = append(cmds, func() tea.Msg { return MessageGoToMainPage{} })
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}
