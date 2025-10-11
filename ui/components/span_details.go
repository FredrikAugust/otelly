package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type SpanDetailsModel struct {
	span     *ptrace.Span
	resource *pcommon.Resource

	width  int
	height int
}

// Init implements tea.Model.
func (s SpanDetailsModel) Init() tea.Cmd {
	return nil
}

func CreateSpanDetailsModel() *SpanDetailsModel {
	return &SpanDetailsModel{
		span:     nil,
		resource: nil,

		width:  0,
		height: 0,
	}
}

type MessageSetSelectedSpan struct {
	Span     *ptrace.Span
	Resource *pcommon.Resource
}

func (s SpanDetailsModel) Update(msg tea.Msg) (SpanDetailsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case MessageSetSelectedSpan:
		s.span = msg.Span
		s.resource = msg.Resource
	}
	return s, nil
}

func (s SpanDetailsModel) View() string {
	box := lipgloss.NewStyle().
		Width(s.width-2).
		Height(s.height-2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	if s.span == nil {
		return box.Foreground(lipgloss.Color("#afafb2")).Align(lipgloss.Center, lipgloss.Center).Render("No span selected")
	}

	return box.
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#6b6b6e")).Render("Span "),
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Render(s.span.SpanID().String()),
			),
			lipgloss.NewStyle().Bold(true).Render(s.span.Name()),
			s.attributeView(),
		))
}

func (s SpanDetailsModel) attributeView() string {
	attributeStrs := make([]string, 0)
	for key, value := range s.span.Attributes().All() {
		attributeStrs = append(
			attributeStrs,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Width(16).Render(key),
				" ",
				lipgloss.NewStyle().Render(value.Str()),
			),
		)
	}

	return lipgloss.NewStyle().
		Width(s.width-6).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render("Attributes"),
				lipgloss.JoinVertical(lipgloss.Left, attributeStrs...),
			),
		)
}

func (s *SpanDetailsModel) SetWidth(w int) {
	s.width = w
}

func (s *SpanDetailsModel) SetHeight(h int) {
	s.height = h
}
