package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type SpanDetailPanelModel struct {
	span *db.Span

	height int
	width  int
}

func NewSpanDetailPanelModel() SpanDetailPanelModel {
	return SpanDetailPanelModel{}
}

func (m SpanDetailPanelModel) Init() tea.Cmd {
	return nil
}

func (m SpanDetailPanelModel) Update(msg tea.Msg) (SpanDetailPanelModel, tea.Cmd) {
	return m, nil
}

func (m SpanDetailPanelModel) View() string {
	container := lipgloss.
		NewStyle().
		Width(m.width).
		MaxWidth(m.width).
		Height(m.height).
		MaxHeight(m.height)

	if m.span == nil {
		return container.Align(lipgloss.Center, lipgloss.Center).Render("No span selected")
	}

	return container.Render(
		helpers.VStack(
			lipgloss.NewStyle().Render("Span", m.span.ID, "•", m.spanKindView()),
			lipgloss.NewStyle().Render(m.span.Name, "•", m.span.Duration.Round(time.Microsecond).String()),
			"", // spacer
			m.traceView(),
			"", // spacer
			m.attributeView(),
			"", // spacer
			m.resourceView(),
		),
	)
}

func (m SpanDetailPanelModel) resourceView() string {
	return "resource"
}

func (m SpanDetailPanelModel) attributeView() string {
	return "attribute"
}

func (m SpanDetailPanelModel) traceView() string {
	return "Trace"
}

func (m SpanDetailPanelModel) spanKindView() string {
	return m.span.Kind
}

func (m *SpanDetailPanelModel) SetHeight(i int) {
	m.height = i
}

func (m *SpanDetailPanelModel) SetWidth(i int) {
	m.width = i
}

func (m *SpanDetailPanelModel) SetSpan(span *db.Span) {
	m.span = span
}
