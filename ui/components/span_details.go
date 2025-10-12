package components

import (
	"log/slog"
	"reflect"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"go.uber.org/zap"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
)

type SpanDetailsModel struct {
	span     *db.Span
	resource *db.Resource

	resourceSpanHistory []db.SpansPerMinuteBucket

	width  int
	height int

	chart tslc.Model

	db *db.Database
}

// Init implements tea.Model.
func (s SpanDetailsModel) Init() tea.Cmd {
	return func() tea.Msg {
		return MessageUpdateGraph{}
	}
}

func CreateSpanDetailsModel(db *db.Database) *SpanDetailsModel {
	return &SpanDetailsModel{
		span:     nil,
		resource: nil,

		chart: tslc.New(
			30, 10,
			tslc.WithUpdateHandler(tslc.SecondUpdateHandler(60)),
			tslc.WithXLabelFormatter(tslc.HourTimeLabelFormatter()),
		),

		db: db,

		width:  0,
		height: 0,
	}
}

type MessageSetSelectedSpan struct {
	SpanID string
}

type MessageUpdateGraph struct{}

func (s SpanDetailsModel) Update(msg tea.Msg) (SpanDetailsModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case MessageSetSelectedSpan:
		if s.span != nil && msg.SpanID == s.span.ID {
			// Don't query the span again when we already have it
			break
		}

		span, err := s.db.GetSpan(msg.SpanID)
		if err != nil {
			slog.Warn("could not get span with resource in span details", "spanID", msg.SpanID, "error", err)
			return s, nil
		}

		s.span = span
		res, err := s.db.GetResource(span.ResourceID)
		if err != nil {
			slog.Warn("could not get resource", "error", err)
			return s, nil
		}

		s.resource = res

		s.updateGraph()
	case MessageUpdateGraph:
		cmds = append(cmds, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return MessageUpdateGraph{}
		}))
		s.updateGraph()
	}

	s.chart, _ = s.chart.Update(msg)
	if s.resourceSpanHistory != nil {
		s.chart.Resize(s.width-4, 10)
		s.chart.SetViewTimeRange(time.Now().Add(-10*time.Minute), time.Now())
		s.chart.DrawBraille()
	}

	return s, tea.Batch(cmds...)
}

func (s *SpanDetailsModel) updateGraph() {
	if s.resource == nil {
		return
	}
	var err error
	s.resourceSpanHistory, err = s.db.SpansPerMinuteForService(s.resource.ID)
	if err != nil {
		zap.L().Warn("could not get span history for svc", zap.String("id", s.resource.ID), zap.Error(err))
	}

	s.chart.ClearAllData()
	for _, history := range s.resourceSpanHistory {
		s.chart.Push(tslc.TimePoint{Time: history.Timestamp, Value: float64(history.SpanCount)})
	}
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
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Render(s.span.ID),
			),
			lipgloss.NewStyle().Bold(true).Render(s.span.Name),
			s.attributeView(),
			s.resourceView(),
		))
}

func (s SpanDetailsModel) resourceView() string {
	baseStyle := lipgloss.NewStyle().
		Width(s.width - 4).
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
					lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Render("Service "),
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

func (s SpanDetailsModel) attributeView() string {
	attributeStrs := make([]string, 0)

	keys := make([]string, 0)
	for k := range s.span.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := s.span.Attributes[k]

		// Only show strings
		if reflect.TypeOf(v) != reflect.TypeFor[string]() {
			continue
		}

		attributeStrs = append(
			attributeStrs,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Width(16).Render(k),
				" ",
				lipgloss.NewStyle().Render(v.(string)),
			),
		)
	}

	return lipgloss.NewStyle().
		Width(s.width - 4).
		MarginTop(1).
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
