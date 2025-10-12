// Package components contain UI elements used throughout the app
package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type SpanTableModel struct {
	spans      []db.Span
	tableModel *table.Model

	width  int
	height int

	db *db.Database
}

func (s SpanTableModel) Init() tea.Cmd {
	return nil
}

func CreateSpanTableModel(db *db.Database) *SpanTableModel {
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
		tableModel: &tableModel,
		db:         db,
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
			rows = append(rows, table.Row{span.Name, span.ServiceName, span.StartTime.Format("15:04:05.000"), span.Duration.String()})
		}
		s.tableModel.SetRows(rows)
	}

	*s.tableModel, cmd = s.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	newSelection := s.SelectedSpanID()
	if selection != newSelection {
		cmds = append(cmds, setSelectedSpanCmd(newSelection))
	}

	return s, tea.Batch(cmds...)
}

func setSelectedSpanCmd(spanID string) tea.Cmd {
	return func() tea.Msg {
		return MessageSetSelectedSpan{SpanID: spanID}
	}
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
	s.tableModel.SetHeight(s.height - 1)
}

func (s *SpanTableModel) SetWidth(w int) {
	s.width = w
	s.tableModel.SetWidth(s.width)
}

// SelectedSpanID returns the spanID if it exists,
// and "" if not
func (s SpanTableModel) SelectedSpanID() string {
	selectedRow := s.tableModel.SelectedRow()
	if selectedRow == nil {
		return ""
	}

	return s.spans[s.tableModel.Cursor()].ID
}
