// Package ui contains the bubbletea user interface
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
)

type Page uint

const (
	PageSpans Page = iota
	PageLogs
)

type EntryModel struct {
	currentPage Page

	width  int
	height int

	spans []db.Span
	logs  []db.Log

	spansPageModel SpansPageModel
}

func NewEntryModel(spans []db.Span, logs []db.Log) tea.Model {
	return EntryModel{
		currentPage: PageSpans,
		spans:       spans,
		logs:        logs,

		spansPageModel: NewSpansPageModel(spans),
	}
}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.spansPageModel.height = msg.Height - 3 // - header
		m.spansPageModel.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m EntryModel) View() string {
	var page string

	switch m.currentPage {
	case PageSpans:
		page = m.spansPageModel.View()
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(
		helpers.VStack(
			m.HeaderView(),
			page,
		),
	)
}

func (m EntryModel) HeaderView() string {
	container := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(helpers.ColorBorder).Width(m.width - 2).Height(1)

	spans := helpers.NavigationPillBaseStyle
	logs := helpers.NavigationPillBaseStyle

	switch m.currentPage {
	case PageSpans:
		spans = spans.Background(helpers.ColorAccentBackground).Foreground(helpers.ColorBlack)
	case PageLogs:
		logs = logs.Background(helpers.ColorAccentBackground).Foreground(helpers.ColorBlack)
	}

	return container.Render(
		helpers.HStack(
			spans.Render("(1) Spans"),
			logs.Render("(2) Logs"),
		),
	)
}
