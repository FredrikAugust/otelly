package helpers

import tea "github.com/charmbracelet/bubbletea"

func Cmdize(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
