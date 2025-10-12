// Package components contain UI elements used throughout the app
package components

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type SpanTableModel struct {
	spans []db.Span
	table *table.Model

	width  int
	height int

	db *db.Database
}

func (s SpanTableModel) Init() tea.Cmd {
	return nil
}

func CreateSpanTableModel(db *db.Database) *SpanTableModel {
	cols := []table.Column{
		// StatusCode gutter
		{
			Title: "",
			Width: 1,
		},
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
		table: &tableModel,
		db:    db,
	}
}

type MessageUpdateRootSpanRows struct{}

func (s SpanTableModel) Update(msg tea.Msg) (SpanTableModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	selection := s.SelectedSpanID()

	switch msg.(type) {
	case MessageUpdateRootSpanRows:
		s.spans = s.db.GetSpans()
		rows := make([]table.Row, 0)
		for _, span := range s.spans {
			rows = append(rows, table.Row{statusCodeView(span.StatusCode), span.Name, span.ServiceName, span.StartTime.Format("15:04:05.000"), span.Duration.String()})
		}
		s.table.SetRows(rows)
	}

	*s.table, cmd = s.table.Update(msg)
	cmds = append(cmds, cmd)

	newSelection := s.SelectedSpanID()
	if selection != newSelection {
		cmds = append(cmds, setSelectedSpanCmd(newSelection))
	}

	return s, tea.Batch(cmds...)
}

func statusCodeView(statusCode string) string {
	switch statusCode {
	case "Unset":
		return ""
	case "Ok":
		return "✅"
	case "Error":
		return "⚠️"
	}

	slog.Warn("got status code which shouldn't exist", "code", statusCode)
	return ""
}

func setSelectedSpanCmd(spanID string) tea.Cmd {
	return func() tea.Msg {
		return MessageSetSelectedSpan{SpanID: spanID}
	}
}

func (s SpanTableModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(s.width).Height(s.height-1).Render(s.table.View()),
		s.table.HelpView(),
	)
}

func (s *SpanTableModel) SetHeight(h int) {
	s.height = h
	s.table.SetHeight(s.height - 1)
}

func (s *SpanTableModel) SetWidth(w int) {
	s.width = w
	s.table.SetWidth(s.width)
}

// SelectedSpanID returns the spanID if it exists,
// and "" if not
func (s SpanTableModel) SelectedSpanID() string {
	selectedRow := s.table.SelectedRow()
	if selectedRow == nil {
		return ""
	}

	return s.spans[s.table.Cursor()].ID
}
