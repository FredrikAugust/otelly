package ui

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type SpanWaterfallModel struct {
	tree helpers.TraceTreeNode

	width  int
	height int

	renderLine func(span *db.SpanWithResource, width int) string

	cursor int
}

func CreateSpanWaterfallModel() SpanWaterfallModel {
	return SpanWaterfallModel{
		cursor: -1,
		height: 10,
		renderLine: func(span *db.SpanWithResource, width int) string {
			if len(span.Name) > width {
				return span.Name[:width]
			}
			return span.Name
		},
	}
}

func (m SpanWaterfallModel) Update(msg tea.Msg) (SpanWaterfallModel, tea.Cmd) {
	return m, nil
}

func (m SpanWaterfallModel) Init() tea.Cmd {
	return nil
}

func (m SpanWaterfallModel) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(m.width)

	var (
		lineIdx          int
		minTime, maxTime = m.tree.GetTimeRange()
	)

	lines := make([]string, 0)
	for _, item := range m.tree.All() {
		width := int(math.Round((float64(item.Span.Duration) / float64(maxTime.Sub(minTime))) * float64(m.width)))

		// Sometimes really short spans would report as 0
		width = max(width, 1)

		// The duration from the last startTime (parent) to this one. Duration between starts.
		delayStr := TextTertiary.Render("↪" + item.Span.StartTime.Sub(item.ParentStartTime).Round(time.Millisecond).String())
		marginLeft := int(math.Round(float64(item.Span.StartTime.Sub(minTime)) / float64(maxTime.Sub(minTime)) * float64(m.width)))

		var backgroundColor, foregroundColor lipgloss.Color
		if lineIdx == m.cursor {
			backgroundColor = ColorForeground
			foregroundColor = ColorAccent
		} else {
			backgroundColor = ColorAccent
			foregroundColor = ColorForeground
		}

		// Compensate for the max(width, 1) which is for very short spans
		if marginLeft == m.width {
			marginLeft -= 1
		}

		// The textual content
		body := m.renderLine(&item.Span, width)

		// If we have room to add the delayStr before the span "block", we add it
		if lipgloss.Width(delayStr) <= marginLeft {
			marginLeft -= lipgloss.Width(delayStr)
			body = lipgloss.NewStyle().MarginLeft(marginLeft).Render(
				lipgloss.JoinHorizontal(0, delayStr, lipgloss.NewStyle().Foreground(foregroundColor).Background(backgroundColor).Width(width).Render(body)),
			)
		} else {
			body = lipgloss.NewStyle().MarginLeft(marginLeft).Render(
				lipgloss.NewStyle().Foreground(foregroundColor).Background(backgroundColor).Width(width).Render(body),
			)
		}

		lines = append(
			lines,
			body,
		)

		lineIdx++
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	if lipgloss.Height(content) > m.height {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinVertical(0, lines[:m.height-1]...),
			TextTertiary.Render(fmt.Sprintf("+ %v spans", len(lines)-m.height+1)),
		)
	}

	return baseStyle.Render(
		content,
	)
}
