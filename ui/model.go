// Package ui contains code specific to the rendering and
// state management of the user interface.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	PageMain = iota
)

type Model struct {
	CurrentPage int
}

func NewModel() *Model {
	return &Model{
		CurrentPage: PageMain,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
