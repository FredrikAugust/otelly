// Package styling contains common helpers for typography in the UI
// to avoid repetition
package styling

import "github.com/charmbracelet/lipgloss"

var (
	ColorSecondary = lipgloss.Color("#afafb2")
	ColorTertiary  = lipgloss.Color("#6b6b6e")
)

var (
	TextHeading   = lipgloss.NewStyle().Bold(true)
	TextSecondary = lipgloss.NewStyle().Foreground(ColorSecondary)
	TextTertiary  = lipgloss.NewStyle().Foreground(ColorTertiary)
)
