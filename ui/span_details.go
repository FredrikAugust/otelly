package ui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type SpanDetailsModel struct {
	span *db.Span

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
		span: nil,

		db: db,

		width:  0,
		height: 0,

		waterfallModel:     CreateSpanWaterfallModel(db),
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
		m.resourceModel.width = m.width - 4
	case MessageResetDetail:
		m.span = nil
		m.spanAttributeModel.attributes = nil
		m.resourceModel.resource = nil
	case MessageSetSelectedSpan:
		if m.span != nil && msg.SpanID == m.span.ID {
			// Don't query the span again when we already have it
			break
		}

		span, err := m.db.GetSpan(msg.SpanID)
		if err != nil {
			slog.Warn("could not get span with resource in span details", "spanID", msg.SpanID, "error", err)
		}

		m.span = span
		res, err := m.db.GetResource(span.ResourceID)
		if err != nil {
			slog.Warn("could not get resource", "error", err)
		}

		m.spanAttributeModel.attributes = m.span.Attributes
		m.resourceModel.resource = res
	}

	m.waterfallModel, cmd = m.waterfallModel.Update(msg)
	cmds = append(cmds, cmd)

	m.spanAttributeModel, cmd = m.spanAttributeModel.Update(msg)
	cmds = append(cmds, cmd)

	m.resourceModel, cmd = m.resourceModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SpanDetailsModel) View() string {
	box := lipgloss.NewStyle().
		Width(m.width-2).
		Height(m.height-2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	if m.span == nil {
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
			m.waterfallModel.View(),
			"",
			m.spanAttributeModel.View(m.width-4),
			m.resourceModel.View(),
		))
}
