package ui

import (
	"reflect"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpanAttributeModel struct {
	attributes map[string]any
}

func CreateSpanAttributeModel() SpanAttributeModel {
	return SpanAttributeModel{
		attributes: nil,
	}
}

func (s SpanAttributeModel) Init() tea.Cmd {
	return nil
}

func (s SpanAttributeModel) Update(msg tea.Msg) (SpanAttributeModel, tea.Cmd) {
	return s, nil
}

func (s *SpanAttributeModel) SetAttributes(attributes map[string]any) {
	s.attributes = attributes
}

func (s SpanAttributeModel) View(w int) string {
	attributeStrs := make([]string, 0)

	keys := make([]string, 0)
	for k := range s.attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := s.attributes[k]

		// Only show strings
		if reflect.TypeOf(v) != reflect.TypeFor[string]() {
			continue
		}

		attributeStrs = append(
			attributeStrs,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				TextSecondary.Width(30).Render(k),
				" ",
				lipgloss.NewStyle().Render(v.(string)),
			),
		)
	}

	return lipgloss.NewStyle().
		Width(w).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				TextHeading.Render("Attributes"),
				lipgloss.JoinVertical(lipgloss.Left, attributeStrs...),
			),
		)
}
