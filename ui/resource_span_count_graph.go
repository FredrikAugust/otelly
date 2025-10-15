package ui

import (
	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/db"
)

type ResourceSpanCountGraphModel struct {
	chart tslc.Model
	db    *db.Database

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
	}

	m.chart, _ = m.chart.Update(msg)

	return m, nil
}

func (m ResourceSpanCountGraphModel) View() string {
	return m.chart.View()
}

func (m *ResourceSpanCountGraphModel) updateGraph(aggregation []db.SpansPerMinuteForServiceModel) {
	m.chart.ClearAllData()
	for _, history := range aggregation {
		m.chart.Push(tslc.TimePoint{Time: history.Timestamp, Value: float64(history.SpanCount)})
	}

	m.chart.Resize(m.width, 10)
	m.chart.DrawBraille()
}
