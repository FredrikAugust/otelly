package ui

import (
	"log/slog"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/ui/components"
	"github.com/fredrikaugust/otelly/ui/pages"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("new msg", "type", reflect.TypeOf(msg).Name())

	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, 0)
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize.height = msg.Height
		m.windowSize.width = msg.Width
	case components.MessageGoToTrace:
		m.currentPage = PageTrace
		cmds = append(cmds, m.tracePageModel.Init())
	case pages.MessageReturnToMainPage:
		m.currentPage = PageMain
		cmds = append(cmds, m.mainPageModel.Init())
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case messageResourceSpansArrived:
		cmds = append(
			cmds,
			listenForNewSpans(m.bus.TraceBus),
			m.InsertResourceSpans(msg.resourceSpans),
		)
	}

	switch m.currentPage {
	case PageMain:
		*m.mainPageModel, cmd = m.mainPageModel.Update(msg)
		cmds = append(cmds, cmd)
	case PageTrace:
		*m.tracePageModel, cmd = m.tracePageModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) InsertResourceSpans(resourceSpans ptrace.ResourceSpans) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	m.db.InsertResourceSpans(resourceSpans)

	*m.mainPageModel, cmd = m.mainPageModel.Update(components.MessageUpdateRootSpanRows{})
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return messageResourceSpansArrived{resourceSpans: <-spanChan}
	}
}
