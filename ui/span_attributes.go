package ui

import (
	"fmt"
	"reflect"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpanAttributeModel struct {
	attributes map[string]any

	width  int
	height int
}

func CreateSpanAttributeModel() SpanAttributeModel {
	return SpanAttributeModel{
		attributes: nil,
		height:     5,
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

	// we treat 0 as show everything
	if len(attributeStrs)+1 > m.height {
		attributeStrs = attributeStrs[:m.height-2] // subtract one for help text and one for title
		attributeStrs = append(attributeStrs, TextTertiary.Render(fmt.Sprintf("+ %v attributes", len(keys)-m.height)))
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
