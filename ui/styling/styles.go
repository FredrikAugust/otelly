// Package styling contains common helpers for ui styling
package styling

import "github.com/charmbracelet/lipgloss"

var (
	ColorSecondary = lipgloss.Color("#afafb2")
	ColorTertiary  = lipgloss.Color("#6b6b6e")
	ColorAccent    = lipgloss.Color("#2e6f40")
)

var (
	TextHeading   = lipgloss.NewStyle().Bold(true)
	TextSecondary = lipgloss.NewStyle().Foreground(ColorSecondary)
	TextTertiary  = lipgloss.NewStyle().Foreground(ColorTertiary)
)
