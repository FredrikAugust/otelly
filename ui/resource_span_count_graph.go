package ui

import (
	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/db"
	"go.uber.org/zap"
)

type ResourceSpanCountGraphModel struct {
	chart       tslc.Model
	aggregation []db.SpansPerMinuteForServiceModel
	resourceID  string
	db          *db.Database

	width  int
	height int
}

func CreateResourceSpanCountGraphModel(db *db.Database) ResourceSpanCountGraphModel {
	return ResourceSpanCountGraphModel{
		chart: tslc.New(
			0, 10,
			tslc.WithUpdateHandler(tslc.SecondUpdateHandler(60)),
			tslc.WithXLabelFormatter(tslc.HourTimeLabelFormatter()),
		),
		db: db,
	}
}

func (m ResourceSpanCountGraphModel) Init() tea.Cmd {
	return tea.Batch(
		m.chart.Init(),
	)
}

func (m ResourceSpanCountGraphModel) Update(msg tea.Msg) (ResourceSpanCountGraphModel, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.chart.Resize(m.width, m.chart.Height())
		m.chart.DrawBraille()
	case MessageSetSelectedSpan:
		m.updateGraph()
	}

	m.chart, _ = m.chart.Update(msg)

	return m, nil
}

func (m ResourceSpanCountGraphModel) View() string {
	return m.chart.View()
}

func (m *ResourceSpanCountGraphModel) updateGraph() {
	if m.resourceID == "" {
		return
	}
	var err error
	m.aggregation, err = m.db.SpansPerMinuteForService(m.resourceID)
	if err != nil {
		zap.L().Warn("could not get span history for svc", zap.String("id", m.resourceID), zap.Error(err))
	}

	m.chart.ClearAllData()
	for _, history := range m.aggregation {
		m.chart.Push(tslc.TimePoint{Time: history.Timestamp, Value: float64(history.SpanCount)})
	}

	m.chart.Resize(m.width, 10)
	m.chart.DrawBraille()
}
