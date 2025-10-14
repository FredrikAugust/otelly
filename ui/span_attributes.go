package ui

import (
	"reflect"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpanAttributeModel struct {
	attributes map[string]any

	width int
}

func CreateSpanAttributeModel() SpanAttributeModel {
	return SpanAttributeModel{
		attributes: nil,
	}
}

func (m SpanAttributeModel) Init() tea.Cmd {
	return nil
}

func (m SpanAttributeModel) Update(msg tea.Msg) (SpanAttributeModel, tea.Cmd) {
	return m, nil
}

func (m *SpanAttributeModel) SetAttributes(attributes map[string]any) {
	m.attributes = attributes
}

func (m SpanAttributeModel) View() string {
	attributeStrs := make([]string, 0)

	keys := make([]string, 0)
	for k := range m.attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m.attributes[k]

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
		Width(m.width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				TextHeading.Render("Attributes"),
				lipgloss.JoinVertical(lipgloss.Left, attributeStrs...),
			),
		)
}
