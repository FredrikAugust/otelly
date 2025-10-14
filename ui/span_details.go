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
func (s SpanDetailsModel) Init() tea.Cmd {
	return tea.Batch(
		s.waterfallModel.Init(),
		s.spanAttributeModel.Init(),
		s.resourceModel.Init(),
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

func (s SpanDetailsModel) Update(msg tea.Msg) (SpanDetailsModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case MessageResetDetail:
		s.span = nil
		s.spanAttributeModel.attributes = nil
		s.resourceModel.resource = nil
	case MessageSetSelectedSpan:
		if s.span != nil && msg.SpanID == s.span.ID {
			// Don't query the span again when we already have it
			break
		}

		span, err := s.db.GetSpan(msg.SpanID)
		if err != nil {
			slog.Warn("could not get span with resource in span details", "spanID", msg.SpanID, "error", err)
		}

		s.span = span
		res, err := s.db.GetResource(span.ResourceID)
		if err != nil {
			slog.Warn("could not get resource", "error", err)
		}

		s.spanAttributeModel.attributes = s.span.Attributes
		s.resourceModel.SetResource(res)
	}

	s.waterfallModel, cmd = s.waterfallModel.Update(msg)
	cmds = append(cmds, cmd)

	s.spanAttributeModel, cmd = s.spanAttributeModel.Update(msg)
	cmds = append(cmds, cmd)

	s.resourceModel, cmd = s.resourceModel.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s SpanDetailsModel) View() string {
	box := lipgloss.NewStyle().
		Width(s.width-2).
		Height(s.height-2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	if s.span == nil {
		return box.Foreground(ColorSecondary).Align(lipgloss.Center, lipgloss.Center).Render("No span selected")
	}

	return box.
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				TextTertiary.Render("Span "),
				TextSecondary.Render(s.span.ID),
			),
			lipgloss.NewStyle().Bold(true).Render(s.span.Name),
			s.waterfallModel.View(),
			"",
			s.spanAttributeModel.View(s.width-4),
			s.resourceModel.View(),
		))
}

func (s *SpanDetailsModel) SetWidth(w int) {
	s.width = w
	s.waterfallModel.width = w - 4
	s.resourceModel.SetWidth(w - 4)
}

func (s *SpanDetailsModel) SetHeight(h int) {
	s.height = h
}
