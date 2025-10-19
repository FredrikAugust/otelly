package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
	"go.uber.org/zap"
)

type SpanDetailsModel struct {
	span db.SpanWithResource

	resourceSpanHistory []db.SpansPerMinuteForServiceModel

	width  int
	height int

	db *db.Database

	waterfallModel     SpanWaterfallModel
	spanAttributeModel SpanAttributeModel
	resourceModel      ResourceModel
}

// Init implements tea.Model.
func (m SpanDetailsModel) Init() tea.Cmd {
	return tea.Batch(
		m.waterfallModel.Init(),
		m.spanAttributeModel.Init(),
		m.resourceModel.Init(),
	)
}

func CreateSpanDetailsModel(db *db.Database) SpanDetailsModel {
	return SpanDetailsModel{
		db: db,

		width:  0,
		height: 0,

		waterfallModel:     CreateSpanWaterfallModel(),
		spanAttributeModel: CreateSpanAttributeModel(),
		resourceModel:      CreateResourceModel(db),
	}
}

func (m SpanDetailsModel) Update(msg tea.Msg) (SpanDetailsModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.waterfallModel.width = m.width - 4
		m.waterfallModel.height = 10
		m.resourceModel.width = m.width - 4
		m.spanAttributeModel.width = m.width - 4
	case MessageSetSelectedSpan:
		m.span = msg.Span
		cmds = append(cmds, m.getTrace(msg.Span.TraceID))
		m.spanAttributeModel.attributes = msg.Span.Attributes
		m.resourceModel, cmd = m.resourceModel.getResourceAndResourceAggregation(msg.Span.ResourceID)
		cmds = append(cmds, cmd)
	case MessageReceivedTraceSpans:
		if len(msg.Spans) == 0 {
			break
		}
		tree, err := helpers.BuildTree(msg.Spans)
		if err != nil {
			zap.L().Warn("could not build tree", zap.Error(err))
			break
		}
		m.waterfallModel.tree = tree
	}

	m.waterfallModel, cmd = m.waterfallModel.Update(msg)
	cmds = append(cmds, cmd)

	m.spanAttributeModel, cmd = m.spanAttributeModel.Update(msg)
	cmds = append(cmds, cmd)

	m.resourceModel, cmd = m.resourceModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SpanDetailsModel) getTrace(traceID string) tea.Cmd {
	return func() tea.Msg {
		spans, err := m.db.GetSpansForTrace(traceID)
		if err != nil {
			zap.L().Warn("could not get trace", zap.String("traceID", traceID))
		}
		return MessageReceivedTraceSpans{
			Spans: spans,
		}
	}
}

func (m SpanDetailsModel) View() string {
	box := lipgloss.NewStyle().
		Width(m.width-2).
		Height(m.height-2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorBorderForeground).
		Padding(0, 1)

	if m.span.ID == "" {
		return box.Foreground(ColorSecondary).Align(lipgloss.Center, lipgloss.Center).Render("No span selected")
	}

	return box.
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				TextTertiary.Render("Span "),
				TextSecondary.Render(m.span.ID),
			),
			lipgloss.NewStyle().Bold(true).Render(m.span.Name),
			"",
			m.waterfallModel.View(),
			"",
			m.spanAttributeModel.View(),
			m.resourceModel.View(),
		))
}
