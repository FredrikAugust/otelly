// Package components contain UI elements used throughout the app
package components

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type SpanTableModel struct {
	tableModel *table.Model

	width  int
	height int

	rootSpans []*ptrace.Span
}

func (s SpanTableModel) Init() tea.Cmd {
	return nil
}

func CreateSpanTableModel() *SpanTableModel {
	cols := []table.Column{
		{
			Title: "Name",
			Width: 16,
		},
		{
			Title: "Service",
			Width: 16,
		},
		{
			Title: "Start time",
			Width: 12,
		},
		{
			Title: "Duration",
			Width: 8,
		},
	}

	tableModel := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(table.DefaultStyles()),
	)

	return &SpanTableModel{
		rootSpans:  make([]*ptrace.Span, 0),
		tableModel: &tableModel,
	}
}

func SpanToRow(span ptrace.Span) table.Row {
	return []string{span.SpanID().String(), span.Name(), span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).String()}
}

type MessageNewRootSpan struct {
	Span         *ptrace.Span
	ResourceName string
}

func (s SpanTableModel) Update(msg tea.Msg) (SpanTableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case MessageNewRootSpan:
		s.rootSpans = append(s.rootSpans, msg.Span)
		s.tableModel.SetRows(
			append(
				s.tableModel.Rows(),
				table.Row{
					msg.Span.Name(),
					msg.ResourceName,
					msg.Span.StartTimestamp().AsTime().Format("15:04:05.000"),
					msg.Span.EndTimestamp().AsTime().Sub(msg.Span.StartTimestamp().AsTime()).String(),
				},
			),
		)
	}

	*s.tableModel, cmd = s.tableModel.Update(msg)

	return s, cmd
}

func (s SpanTableModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(s.width).Height(s.height-1).Render(s.tableModel.View()),
		s.tableModel.HelpView(),
	)
}

func (s *SpanTableModel) SetHeight(h int) {
	s.height = h
}

func (s *SpanTableModel) SetWidth(w int) {
	s.width = w
}

// SelectedSpan returns the span, true if a span is selected
// and nil, false if not.
func (s SpanTableModel) SelectedSpan() (*ptrace.Span, bool) {
	idx := s.tableModel.Cursor()
	slog.Info("cursor", "pos", idx, "rootspanlen", len(s.rootSpans))
	if idx >= len(s.rootSpans) || idx < 0 {
		return nil, false
	}
	return s.rootSpans[idx], true
}
