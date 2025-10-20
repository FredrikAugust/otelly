// Package ui contains the bubbletea user interface
package ui

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
	"go.uber.org/zap"
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

	bus *bus.TransportBus
}

func NewEntryModel(spans []db.Span, logs []db.Log, bus *bus.TransportBus) tea.Model {
	return EntryModel{
		currentPage: PageSpans,
		spans:       spans,
		logs:        logs,

		spansPageModel: NewSpansPageModel(db.FilterRootSpans(spans)),
		bus:            bus,
	}
}

func (m EntryModel) Init() tea.Cmd {
	return tea.Batch(
		m.spansPageModel.Init(),
		m.listenForLogs(),
		m.listenForSpans(),
	)
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	zap.L().Debug("received tea.Msg", zap.String("type", reflect.TypeOf(msg).Name()))

	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.spansPageModel.SetHeight(msg.Height - 3) // - header
		m.spansPageModel.SetWidth(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			cmds = append(cmds, tea.Quit)
		}
	case MsgNewSpans:
		cmds = append(cmds, m.listenForSpans())
		m.updateSpans(msg.spans)
	case MsgNewLogs:
		cmds = append(cmds, m.listenForLogs())
		m.updateLogs(msg.logs)
	}

	switch m.currentPage {
	case PageSpans:
		m.spansPageModel, cmd = m.spansPageModel.Update(msg)
		cmds = append(cmds, cmd)
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

func (m EntryModel) listenForSpans() tea.Cmd {
	return func() tea.Msg {
		return MsgNewSpans{<-m.bus.SpanBus}
	}
}

func (m EntryModel) listenForLogs() tea.Cmd {
	return func() tea.Msg {
		return MsgNewLogs{<-m.bus.LogBus}
	}
}

func (m *EntryModel) updateSpans(spans []db.Span) {
	m.spans = spans
	m.spansPageModel.SetSpans(db.FilterRootSpans(spans))
}

func (m *EntryModel) updateLogs(logs []db.Log) {
	m.logs = logs
}
