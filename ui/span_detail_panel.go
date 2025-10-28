package ui

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/flamegraph"
	"github.com/fredrikaugust/otelly/ui/helpers"
	"go.uber.org/zap"
)

type SpanDetailPanelModel struct {
	span *db.Span

	tree flamegraph.Node

	height int
	width  int

	db *db.Database
}

func NewSpanDetailPanelModel(db *db.Database) SpanDetailPanelModel {
	return SpanDetailPanelModel{db: db}
}

func (m SpanDetailPanelModel) Init() tea.Cmd {
	return nil
}

func (m SpanDetailPanelModel) Update(msg tea.Msg) (SpanDetailPanelModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case MsgLoadTrace:
		cmds = append(
			cmds,
			func() tea.Msg {
				spans, err := m.db.GetSpansForTrace(context.Background(), msg.traceID)
				if err != nil {
					zap.L().Warn("could not get spans for trace", zap.String("traceID", msg.traceID), zap.Error(err))
					return nil
				}
				node, err := flamegraph.Build(spans, func(s db.Span) flamegraph.NodeInput {
					return flamegraph.NodeInput{
						ID:        s.ID,
						Name:      s.Name,
						Duration:  s.Duration,
						ParentID:  s.ParentSpanID.String,
						StartTime: s.StartTime,
					}
				})
				if err != nil {
					zap.L().Warn("could not create flamegraph for trace", zap.String("traceID", msg.traceID), zap.Error(err))
				}
				return MsgTreeUpdated{tree: node}
			},
		)
	case MsgTreeUpdated:
		m.tree = msg.tree
	}

	return m, tea.Batch(cmds...)
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
	if m.tree.Name == "" {
		return "Trace not set"
	}

	spans := make([]string, 0)

	for _, c := range m.tree.All() {
		width := helpers.Clamp(1, int(c.WidthPct*float64(m.width)), m.width)
		offset := int(c.OffsetPct * float64(m.width))

		name := lipgloss.NewStyle().Render(
			lipgloss.NewStyle().Render(c.Name),
			lipgloss.NewStyle().Faint(true).Render(c.Duration.Round(time.Microsecond).String()),
		)

		spans = append(
			spans,
			helpers.HStack(
				strings.Repeat(" ", offset),
				lipgloss.
					NewStyle().
					Width(width).
					MaxWidth(width).
					Inline(true).
					Background(helpers.ColorBlue).
					Render(name),
			),
		)
	}

	return helpers.VStack(
		"Trace",
		helpers.VStack(spans...),
	)
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

func (m SpanDetailPanelModel) UpdateSpan(span *db.Span) (SpanDetailPanelModel, tea.Cmd) {
	if span == nil {
		m.span = nil
		m.tree = flamegraph.Node{}

		return m, nil
	}

	if m.span != nil && m.span.ID == span.ID {
		return m, nil
	}

	m.span = span

	return m, helpers.Cmdize(MsgLoadTrace{traceID: span.TraceID})
}
