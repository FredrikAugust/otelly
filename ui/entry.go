// Package ui contains the UI for the application
package ui

import (
	"context"
	"log/slog"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/bus"
	"github.com/fredrikaugust/otelly/db"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

func Start(ctx context.Context, bus *bus.TransportBus, db *db.Database) error {
	slog.Info("initializing and running UI")

	p := tea.NewProgram(
		NewModel(bus, db),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.mainPageModel.Init(),
		m.logPageModel.Init(),
		listenForNewSpans(m.bus.TraceBus),
		listenForNewLogs(m.bus.LogBus),
		m.fetchRootSpans(),
		m.fetchLogs(),
	)
}

const (
	PageMain = iota
	PageTrace
	PageLogs
)

type Model struct {
	currentPage int

	mainPageModel  MainPageModel
	tracePageModel TracePageModel
	logPageModel   LogPageModel

	bus *bus.TransportBus
	db  *db.Database

	height int
	width  int
}

func NewModel(bus *bus.TransportBus, db *db.Database) *Model {
	return &Model{
		currentPage:    PageMain,
		bus:            bus,
		db:             db,
		mainPageModel:  CreateMainPageModel(db),
		tracePageModel: CreateTracePageModel(db),
		logPageModel:   CreateLogPageModel(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("entrypoint update", "type", reflect.TypeOf(msg).Name())

	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, 0)
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.mainPageModel.width = msg.Width
		m.mainPageModel.height = msg.Height - 1

		m.tracePageModel.SetWidth(msg.Width)
		m.tracePageModel.SetHeight(msg.Height - 1)
	case MessageGoToTrace:
		m.currentPage = PageTrace
		cmds = append(cmds, m.tracePageModel.Init())
	case MessageGoToMainPage:
		m.currentPage = PageMain
		cmds = append(cmds, m.mainPageModel.Init())
	case MessageGoToLogs:
		m.currentPage = PageLogs
	case MessageResourceSpansArrived:
		cmds = append(
			cmds,
			listenForNewSpans(m.bus.TraceBus),
			m.insertResourceSpans(msg.ResourceSpans),
		)
	case MessageResourceLogsArrived:
		cmds = append(
			cmds,
			listenForNewLogs(m.bus.LogBus),
			m.insertResourceLogs(msg.ResourceLogs),
		)
	case MessageUpdateRootSpans:
		if len(msg.NewRootSpans) > 0 && len(m.mainPageModel.spanTable.spans) == 0 {
			cmds = append(cmds, func() tea.Msg { return MessageSetSelectedSpan{Span: msg.NewRootSpans[0]} })
		}
		m.mainPageModel.SetSpans(msg.NewRootSpans)
	case MessageUpdateLogs:
		m.logPageModel.logs = msg.NewLogs
	}

	switch m.currentPage {
	case PageMain:
		m.mainPageModel, cmd = m.mainPageModel.Update(msg)
		cmds = append(cmds, cmd)
	case PageTrace:
		m.tracePageModel, cmd = m.tracePageModel.Update(msg)
		cmds = append(cmds, cmd)
	case PageLogs:
		m.logPageModel, cmd = m.logPageModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	pageContent := ""

	switch m.currentPage {
	case PageMain:
		pageContent = m.mainPageModel.View()
	case PageTrace:
		pageContent = m.tracePageModel.View()
	case PageLogs:
		pageContent = m.logPageModel.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(m.width).Background(ColorAccent).Bold(true).Render("Otelly"),
		pageContent,
	)
}

func (m Model) fetchRootSpans() tea.Cmd {
	return func() tea.Msg {
		newRootSpans := m.db.GetRootSpans()

		return MessageUpdateRootSpans{
			NewRootSpans: newRootSpans,
		}
	}
}

func (m Model) fetchLogs() tea.Cmd {
	return func() tea.Msg {
		newLogs, err := m.db.GetLogs()
		if err != nil {
			zap.L().Warn("could not fetch logs", zap.Error(err))
			return nil
		}

		return MessageUpdateLogs{
			NewLogs: newLogs,
		}
	}
}

func (m Model) insertResourceSpans(resourceSpans ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		err := m.db.InsertResourceSpans(resourceSpans)
		if err != nil {
			zap.L().Warn("could not insert resource spans", zap.Error(err))
			return nil
		}
		return m.fetchRootSpans()()
	}
}

func (m Model) insertResourceLogs(resourceLogs plog.ResourceLogs) tea.Cmd {
	return func() tea.Msg {
		err := m.db.InsertResourceLogs(resourceLogs)
		if err != nil {
			zap.L().Warn("could not insert resource logs", zap.Error(err))
			return nil
		}
		return m.fetchLogs()()
	}
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return MessageResourceSpansArrived{ResourceSpans: <-spanChan}
	}
}

func listenForNewLogs(logChan chan plog.ResourceLogs) tea.Cmd {
	return func() tea.Msg {
		return MessageResourceLogsArrived{ResourceLogs: <-logChan}
	}
}
