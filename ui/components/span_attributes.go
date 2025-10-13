package components

import (
	"reflect"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/ui/styling"
)

type SpanAttributeModel struct {
	width int

	attributes map[string]any
}

func CreateSpanAttributeModel() *SpanAttributeModel {
	return &SpanAttributeModel{
		width:      0,
		attributes: nil,
	}
}

func (s SpanAttributeModel) Init() tea.Cmd {
	return nil
}

func (s SpanAttributeModel) Update(msg tea.Msg) (SpanAttributeModel, tea.Cmd) {
	return s, nil
}

func (s SpanAttributeModel) View() string {
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
				styling.TextSecondary.Width(24).Render(k),
				" ",
				lipgloss.NewStyle().Render(v.(string)),
			),
		)
	}

	return lipgloss.NewStyle().
		Width(s.width).
		MarginTop(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styling.TextHeading.Render("Attributes"),
				lipgloss.JoinVertical(lipgloss.Left, attributeStrs...),
			),
		)
}
