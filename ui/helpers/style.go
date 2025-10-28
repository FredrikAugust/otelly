package helpers

import "github.com/charmbracelet/lipgloss"

var (
	ColorBlack = lipgloss.Color("#000000")
	ColorBlue  = lipgloss.Color("#0000ff")
	ColorWhite = lipgloss.Color("#ffffff")
	ColorGray  = lipgloss.Color("#aaaaaa")

	ColorBackground       = ColorBlack
	ColorAccentBackground = ColorGray
	ColorBorder           = ColorGray
)

var NavigationPillBaseStyle = lipgloss.NewStyle().Padding(0, 1).MarginRight(1)
