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

	tableModel           TableModel
	spanDetailPanelModel SpanDetailPanelModel
}

func NewSpansPageModel(spans []db.Span, db *db.Database) SpansPageModel {
	tm := NewTableModel()
	tm.SetColumnDefinitions([]ColumnDefinition{
		{3, "Name"},
		{1, "Start time"},
		{1, "Duration"},
	})
	return SpansPageModel{
		spans:                spans,
		tableModel:           tm,
		spanDetailPanelModel: NewSpanDetailPanelModel(db),
	}
}

func (m SpansPageModel) Init() tea.Cmd {
	return tea.Batch(
		helpers.Cmdize(MsgSpanPageUpdateTable{}),
		m.tableModel.Init(),
		m.spanDetailPanelModel.Init(),
	)
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

	item, ok := m.tableModel.SelectedItem().(*spanTableItemDelegate)
	if ok {
		m.spanDetailPanelModel, cmd = m.spanDetailPanelModel.UpdateSpan(item.span)
	} else {
		m.spanDetailPanelModel, cmd = m.spanDetailPanelModel.UpdateSpan(nil)
	}
	cmds = append(cmds, cmd)

	m.spanDetailPanelModel, cmd = m.spanDetailPanelModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SpansPageModel) View() string {
	return helpers.HStack(m.tableView(), m.detailView())
}

func (m SpansPageModel) tableView() string {
	container := lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(helpers.ColorBorder)

	return container.Render(m.tableModel.View())
}

func (m SpansPageModel) detailView() string {
	container := lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(helpers.ColorBorder)

	return container.Render(m.spanDetailPanelModel.View())
}

func (m *SpansPageModel) SetSpans(spans []db.Span) {
	m.spans = spans
	m.updateTable()
}

type spanTableItemDelegate struct {
	span *db.Span
}

func (d spanTableItemDelegate) Content() []string {
	return []string{
		d.span.Name,
		d.span.StartTime.Format("15:04:05"),
		d.span.Duration.Round(time.Microsecond).String(),
	}
}

func (m *SpansPageModel) updateTable() {
	items := make([]TableItemDelegate, len(m.spans))
	for i, span := range m.spans {
		d := &spanTableItemDelegate{span: &span}
		items[i] = d
	}
	m.tableModel.SetItems(items)
}

func (m *SpansPageModel) SetWidth(w int) {
	m.width = w
	m.tableModel.SetWidth(int(math.Floor(float64(w)*2.0/3.0)) - 2)
	m.spanDetailPanelModel.SetWidth(int(math.Ceil(float64(w)*(1.0/3.0))) - 2)
}

func (m *SpansPageModel) SetHeight(h int) {
	m.height = h
	m.tableModel.SetHeight(h - 2)
	m.spanDetailPanelModel.SetHeight(h - 2)
}
