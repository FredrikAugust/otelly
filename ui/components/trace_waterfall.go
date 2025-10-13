package components

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/styling"
	"go.uber.org/zap"
)

type SpanWaterfallModel struct {
	traceID string

	spans []db.GetSpansForTraceModel

	width int

	db *db.Database
}

func CreateSpanWaterfallModel(db *db.Database) *SpanWaterfallModel {
	return &SpanWaterfallModel{db: db}
}

func (m SpanWaterfallModel) Update(msg tea.Msg) (SpanWaterfallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case MessageSetSelectedSpan:
		// Load the trace for this span
		span, _ := m.db.GetSpan(msg.SpanID)
		if span != nil {
			m.traceID = span.TraceID
			spans, err := m.db.GetSpansForTrace(span.TraceID)
			if err != nil {
				zap.L().Warn("could not get trace", zap.Error(err))
			}
			m.spans = spans
		}
	}

	return m, nil
}

func (m SpanWaterfallModel) Init() tea.Cmd {
	return nil
}

func (m SpanWaterfallModel) View() string {
	var minTime, maxTime time.Time

	for i, span := range m.spans {
		if i == 0 {
			minTime = span.StartTime
			maxTime = span.StartTime.Add(span.Duration)
			continue
		}

		if span.StartTime.Before(minTime) {
			minTime = span.StartTime
		}
		if span.StartTime.Add(span.Duration).After(maxTime) {
			maxTime = span.StartTime.Add(span.Duration)
		}
	}

	baseStyle := lipgloss.NewStyle().
		Width(m.width).
		MarginTop(1)

	lines := make([]string, 0, len(m.spans))
	for _, span := range m.spans {
		width := int(math.Floor((float64(span.Duration.Nanoseconds()) / float64(maxTime.Sub(minTime).Nanoseconds())) * float64(m.width)))

		// Sometimes really short spans would report as 0
		width = max(width, 1)

		marginLeft := int(math.Floor(float64(span.StartTime.Sub(minTime).Nanoseconds()) / float64(maxTime.Sub(minTime).Nanoseconds()) * float64(m.width)))
		if marginLeft == m.width {
			marginLeft -= 1
		}

		zap.L().Debug(
			"pos",
			zap.Int("width", width),
			zap.Duration("span_duration", span.Duration),
			zap.Duration("trace_duration", maxTime.Sub(minTime)),
			zap.Int("window", m.width),
			zap.Int("margin", marginLeft),
		)

		name := span.Name
		if len(name) > width {
			name = name[:width]
		}

		lines = append(
			lines,
			lipgloss.NewStyle().
				Width(width).
				MarginLeft(marginLeft).
				Background(lipgloss.Color("32")).
				Render(name),
		)
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top, styling.TextHeading.MarginBottom(1).Render("Trace"), " ", styling.TextSecondary.Render(fmt.Sprintf("(%v)", maxTime.Sub(minTime))),
			),
			lipgloss.JoinVertical(lipgloss.Left, lines...),
		),
	)
}
