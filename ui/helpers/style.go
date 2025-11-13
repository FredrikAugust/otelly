package helpers

import "github.com/charmbracelet/lipgloss"

// Based on Dracula
var (
	ColorPrimary           = lipgloss.Color("#bd93f9")
	ColorPrimaryForeground = lipgloss.Color("#282a36")

	ColorSecondary           = lipgloss.Color("#ff79c6")
	ColorSecondaryForeground = lipgloss.Color("#282a36")

	ColorAccent           = lipgloss.Color("#50fa7b")
	ColorAccentForeground = lipgloss.Color("#282a36")

	ColorBackground = lipgloss.Color("#282a36")
	ColorForeground = lipgloss.Color("#f8f8f2")

	ColorMuted           = lipgloss.Color("#44475a")
	ColorMutedForeground = lipgloss.Color("#6272a4")

	ColorCard           = lipgloss.Color("#44475a")
	ColorCardForeground = lipgloss.Color("#f8f8f2")

	ColorPopover           = lipgloss.Color("#282a36")
	ColorPopoverForeground = lipgloss.Color("#f8f8f2")

	ColorBorder = lipgloss.Color("#6272a4")
	ColorInput  = lipgloss.Color("#44475a")
	ColorRing   = lipgloss.Color("#bd93f9")

	ColorDestructive           = lipgloss.Color("#ff5555")
	ColorDestructiveForeground = lipgloss.Color("#f8f8f2")

	// Legacy/utility colors
	ColorBlack = lipgloss.Color("#000000")
	ColorWhite = lipgloss.Color("#ffffff")
	ColorGray  = lipgloss.Color("#aaaaaa")
)

var NavigationPillBaseStyle = lipgloss.NewStyle().Background(ColorBackground).Padding(0, 1)
