package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorSecondary  = lipgloss.Color("#afafb2")
	ColorTertiary   = lipgloss.Color("#6b6b6e")
	ColorAccent     = lipgloss.Color("#7D56F4")
	ColorBackground = lipgloss.Color("#0a0a0a")
)

var (
	TextHeading   = lipgloss.NewStyle().Bold(true)
	TextSecondary = lipgloss.NewStyle().Foreground(ColorSecondary)
	TextTertiary  = lipgloss.NewStyle().Foreground(ColorTertiary)
)
