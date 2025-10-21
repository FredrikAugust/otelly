// Package helpers contains helpful utils for bubbletea
package helpers

import (
	"github.com/charmbracelet/lipgloss"
)

func HStack(strs ...string) string {
	return lipgloss.JoinHorizontal(0, strs...)
}

func VStack(strs ...string) string {
	return lipgloss.JoinVertical(0, strs...)
}
