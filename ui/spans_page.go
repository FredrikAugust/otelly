package ui

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type SpansPageModel struct {
	spans []db.Span

	width  int
	height int

	tableModel TableModel
}

func NewSpansPageModel(spans []db.Span) SpansPageModel {
	return SpansPageModel{spans: spans, tableModel: NewTableModel()}
}

func (m SpansPageModel) Init() tea.Cmd {
	return helpers.Cmdize(MsgSpanPageUpdateTable{})
}

func (m SpansPageModel) Update(msg tea.Msg) (SpansPageModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg.(type) {
	case MsgSpanPageUpdateTable:
		m.updateTable()
	}

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SpansPageModel) View() string {
	return helpers.HStack(m.tableView(), m.detailView())
}

func (m SpansPageModel) tableView() string {
	// width  = 2/3 of window width
	// height = viewport height -  border
	container := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(helpers.ColorBorder)

	return container.Render(m.tableModel.View())
}

func (m SpansPageModel) detailView() string {
	// width  = 1/3 of window width
	// height = viewport height -  border
	container := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(helpers.ColorBorder).Width(int(math.Ceil(float64(m.width)*1.0/3.0 - 2.0))).Height(m.height - 2)

	return container.Render("detail")
}

func (m *SpansPageModel) SetSpans(spans []db.Span) {
	m.spans = spans
	m.updateTable()
}

func (m *SpansPageModel) updateTable() {
	items := make([]TableItemDelegate, len(m.spans))
	for i, span := range m.spans {
		d := NewDefaultTableItemDelegate()
		d.ContentFn = func() []string {
			return []string{
				span.Name,
				span.StartTime.Format("15:04:05"),
				span.Duration.Round(time.Millisecond).String(),
			}
		}
		items[i] = d
	}
	m.tableModel.SetItems(items)
}

func (m *SpansPageModel) SetWidth(w int) {
	m.width = w
	m.tableModel.SetWidth(int(math.Floor(float64(w)*2.0/3.0 - 2.0)))
}

func (m *SpansPageModel) SetHeight(h int) {
	m.height = h
	m.tableModel.SetHeight(h - 2)
}
