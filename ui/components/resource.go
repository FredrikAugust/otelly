package components

import (
	"time"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"

	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/styling"
)

type ResourceModel struct {
	width int

	db *db.Database

	resource    *db.Resource
	aggregation []db.SpansPerMinuteForServiceModel

	chart tslc.Model
}

func CreateResourceModel(db *db.Database) *ResourceModel {
	return &ResourceModel{
		width: 10,
		db:    db,
		chart: tslc.New(
			30, 10,
			tslc.WithUpdateHandler(tslc.SecondUpdateHandler(60)),
			tslc.WithXLabelFormatter(tslc.HourTimeLabelFormatter()),
		),
	}
}

func (s ResourceModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return MessageUpdateGraph{}
		},
		s.chart.Init(),
	)
}

type MessageUpdateGraph struct{}

func (s ResourceModel) Update(msg tea.Msg) (ResourceModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg.(type) {
	case MessageSetSelectedSpan:
		s.updateGraph()
	case MessageUpdateGraph:
		cmds = append(cmds, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return MessageUpdateGraph{}
		}))
		s.updateGraph()
	case tea.WindowSizeMsg:
		if s.aggregation != nil {
			s.chart.Resize(s.width, 10)
			s.chart.DrawBraille()
		}
	}

	s.chart, _ = s.chart.Update(msg)

	return s, tea.Batch(cmds...)
}

func (s *ResourceModel) updateGraph() {
	if s.resource == nil {
		return
	}
	var err error
	s.aggregation, err = s.db.SpansPerMinuteForService(s.resource.ID)
	if err != nil {
		zap.L().Warn("could not get span history for svc", zap.String("id", s.resource.ID), zap.Error(err))
	}

	s.chart.ClearAllData()
	for _, history := range s.aggregation {
		s.chart.Push(tslc.TimePoint{Time: history.Timestamp, Value: float64(history.SpanCount)})
	}

	s.chart.Resize(s.width, 10)
	s.chart.SetViewTimeRange(time.Now().Add(-10*time.Minute), time.Now())
	s.chart.DrawBraille()
}

func (s ResourceModel) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(s.width).
		MarginTop(1)

	if s.resource == nil {
		return baseStyle.Render("No resource found")
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Resource"),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top,
					styling.TextSecondary.Render("Service "),
					s.resource.ServiceNamespace,
					".",
					s.resource.ServiceName,
				),
				lipgloss.NewStyle().
					MarginTop(1).
					Render(
						s.chart.View(),
					),
			),
		),
	)
}
