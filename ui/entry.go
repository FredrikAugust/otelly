// Package ui contains the bubbletea user interface
package ui

import tea "github.com/charmbracelet/bubbletea"

type EntryModel struct{}

func NewEntryModel() tea.Model {
	return EntryModel{}
}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m EntryModel) View() string {
	return "Entry"
}
