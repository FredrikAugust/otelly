package ui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type TableItemDelegate interface {
	Content() []string
}

type DefaultTableItemDelegate struct {
	ContentFn func() []string
}

func (d DefaultTableItemDelegate) Content() []string {
	return d.ContentFn()
}

func NewDefaultTableItemDelegate() DefaultTableItemDelegate {
	return DefaultTableItemDelegate{
		func() []string { return []string{"empty"} },
	}
}

type TableModel struct {
	items     []TableItemDelegate
	itemViews [][]string
	rowHeight int

	cursorRow    int
	cursorColumn int

	yOffset int

	width  int
	height int
}

func (m *TableModel) SetRowHeight(i int) {
	m.rowHeight = i
}

func NewTableModel() TableModel {
	return TableModel{
		itemViews: make([][]string, 0),
		rowHeight: 1,
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.cursorRow += 1
		case "k", "up":
			m.cursorRow -= 1
		case "l", "right":
			m.cursorColumn += 1
		case "h", "left":
			m.cursorColumn -= 1
		case "g":
			m.cursorRow = 0
		case "G":
			m.cursorRow = len(m.itemViews) - 1
		}

		m.cursorRow = helpers.Clamp(0, m.cursorRow, len(m.items)-1)
		if len(m.itemViews) > 0 {
			m.cursorColumn = helpers.Clamp(0, m.cursorColumn, len(m.itemViews[0])-1)
		} else {
			m.cursorColumn = 0
		}
	}

	m.updateYOffset()

	return m, nil
}

// updateYOffset calculates and sets the yOffset which is how far up/down the viewport is scrolled.
func (m *TableModel) updateYOffset() {
	selectedItemYOffset := m.cursorRow * m.rowHeight
	if selectedItemYOffset >= m.yOffset+m.contentHeight() {
		m.yOffset += selectedItemYOffset - (m.yOffset + m.contentHeight()) + 1
	} else if selectedItemYOffset < m.yOffset {
		m.yOffset -= m.yOffset - selectedItemYOffset
	}
}

func (m TableModel) contentHeight() int {
	return m.height - 2
}

func (m TableModel) View() string {
	if len(m.itemViews) == 0 {
		return "no items"
	}

	numCols := len(m.itemViews[0])

	rows := make([]string, len(m.itemViews))
	for i, cols := range m.itemViews {
		row := ""
		for j, col := range cols {
			style := lipgloss.NewStyle().Width(m.width / numCols).MaxWidth(m.width / numCols).Height(m.rowHeight).MaxHeight(m.rowHeight)
			if m.cursorRow == i && m.cursorColumn == j {
				style = style.Background(helpers.ColorAccentBackground).Foreground(helpers.ColorBlack)
			}
			row += style.Render(col)
		}
		rows[i] = row
	}

	rowStack := helpers.VStack(rows...)
	rowStack = strings.Join(strings.Split(rowStack, "\n")[m.yOffset:], "\n")

	return lipgloss.NewStyle().Height(m.height).Width(m.width).Render(
		helpers.VStack(
			m.HeaderView(),
			lipgloss.NewStyle().Height(m.height-2).MaxHeight(m.height-2).Render(
				rowStack,
			),
			m.HelpView(),
		),
	)
}

func (m TableModel) HelpView() string {
	return helpers.HStack(
		strconv.Itoa(m.cursorRow+1), " / ", strconv.Itoa(len(m.items)),
	)
}

func (m TableModel) HeaderView() string {
	return "header"
}

func (m *TableModel) SetItems(items []TableItemDelegate) {
	m.items = items
	m.itemViews = make([][]string, len(items))

	for i, item := range items {
		m.itemViews[i] = item.Content()
	}
}

func (m *TableModel) SetWidth(i int) {
	m.width = i
}

func (m *TableModel) SetHeight(i int) {
	m.height = i
}
