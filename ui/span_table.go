package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type SpanTableModel struct {
	spans []db.Span
	table table.Model

	width  int
	height int

	db *db.Database
}

func (s SpanTableModel) Init() tea.Cmd {
	return nil
}

func CreateSpanTableModel(db *db.Database) SpanTableModel {
	cols := []table.Column{
		{
			Title: "Name",
			Width: 34,
		},
		{
			Title: "Service",
			Width: 16,
		},
		{
			Title: "Start time",
			Width: 14,
		},
		{
			Title: "Duration",
			Width: 12,
		},
	}

	tableModel := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(table.Styles{
			Selected: lipgloss.NewStyle().Bold(true).Background(ColorAccent),
			Header:   lipgloss.NewStyle().Bold(true),
			Cell:     lipgloss.NewStyle(),
		}),
	)

	return SpanTableModel{
		table: tableModel,
		db:    db,
	}
}

func (s SpanTableModel) Update(msg tea.Msg) (SpanTableModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	selection := s.GetSelectedSpanID()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.table.SetWidth(s.width)
		s.table.SetHeight(s.height - 1) // subtract help height
	case MessageUpdateRootSpans:
		s.spans = msg.NewRootSpans
		rows := make([]table.Row, 0)
		for _, span := range s.spans {
			rows = append(rows, table.Row{span.Name, span.ServiceName, span.StartTime.Format("15:04:05.000"), span.Duration.Round(time.Millisecond).String()})
		}
		s.table.SetRows(rows)
		// TODO: restore Cursor to selection
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			s.db.Clear()
			s.table.SetCursor(0)
			s.spans = make([]db.Span, 0)
			s.table.SetRows(make([]table.Row, 0))
			cmds = append(cmds, func() tea.Msg { return MessageResetDetail{} })
		case "enter":
			if s.GetSelectedSpanID() != "" {
				cmds = append(cmds, func() tea.Msg {
					return MessageGoToTrace{
						TraceID: s.spans[s.table.Cursor()].TraceID,
					}
				})
			}
		}
	}

	s.table, cmd = s.table.Update(msg)
	cmds = append(cmds, cmd)

	if len(s.spans) > 0 {
		newSelection := s.GetSelectedSpanID()
		if selection != newSelection {
			cmds = append(cmds, setSelectedSpanCmd(newSelection))
		}
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
		lipgloss.NewStyle().Height(s.height-1).Render(s.table.View()),
		lipgloss.NewStyle().Render(
			s.table.Help.ShortHelpView(
				append(
					s.table.KeyMap.ShortHelp(),
					key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "clear data")),
				),
			),
		),
	)
}

// GetSelectedSpanID returns the spanID if it exists,
// and "" if not
func (s SpanTableModel) GetSelectedSpanID() string {
	selectedRow := s.table.SelectedRow()
	if selectedRow == nil {
		return ""
	}

	return s.spans[s.table.Cursor()].ID
}
