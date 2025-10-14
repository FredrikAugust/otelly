package ui

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"go.uber.org/zap"
)

type SpanWaterfallModel struct {
	traceID string

	spans []db.GetSpansForTraceModel

	width int

	db *db.Database
}

func CreateSpanWaterfallModel(db *db.Database) SpanWaterfallModel {
	return SpanWaterfallModel{db: db}
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
	baseStyle := lipgloss.NewStyle().
		Width(m.width).
		MarginTop(1)

	minTime, maxTime, lines := WaterfallLinesForSpans(m.width, m.spans, func(span *db.GetSpansForTraceModel) string { return span.Name })

	numLines := len(lines)

	// Genius trick to avoid having to deal with singular/plural span(s) string
	if numLines > 7 {
		lines = lines[:5]
		lines = append(lines, TextSecondary.Render(fmt.Sprintf("+ %v more spans", numLines-5)))
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top, TextHeading.MarginBottom(1).Render("Trace"), " ", TextSecondary.Render(fmt.Sprintf("(%v)", maxTime.Sub(minTime))),
			),
			lipgloss.JoinVertical(lipgloss.Left, lines...),
		),
	)
}

func WaterfallLinesForSpans(w int, spans []db.GetSpansForTraceModel, lineContent func(span *db.GetSpansForTraceModel) string) (time.Time, time.Time, []string) {
	var minTime, maxTime time.Time

	for i, span := range spans {
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

	lines := make([]string, 0, len(spans))
	for _, span := range spans {
		width := int(math.Round((float64(span.Duration.Nanoseconds()) / float64(maxTime.Sub(minTime).Nanoseconds())) * float64(w)))

		// Sometimes really short spans would report as 0
		width = max(width, 1)

		marginLeft := int(math.Round(float64(span.StartTime.Sub(minTime).Nanoseconds()) / float64(maxTime.Sub(minTime).Nanoseconds()) * float64(w)))
		if marginLeft == w {
			marginLeft -= 1
		}

		body := lineContent(&span)
		if len(body) > width {
			body = body[:width]
		}

		lines = append(
			lines,
			lipgloss.NewStyle().
				Width(width).
				MarginLeft(marginLeft).
				Background(ColorAccent).
				Render(body),
		)
	}

	return minTime, maxTime, lines
}
