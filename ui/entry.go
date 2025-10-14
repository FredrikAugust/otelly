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
	"go.opentelemetry.io/collector/pdata/ptrace"
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
		listenForNewSpans(m.bus.TraceBus),
		func() tea.Msg {
			// In case we're runnign DB with file we can keep the seed around, and
			// this will make it load into the UI immediately.
			return MessageUpdateRootSpanRows{}
		},
	)
}

const (
	PageMain = iota
	PageTrace
)

type windowSize struct {
	height int
	width  int
}
type Service struct {
	Namespace string
	Name      string
}

type Model struct {
	currentPage int

	mainPageModel  MainPageModel
	tracePageModel TracePageModel

	bus *bus.TransportBus
	db  *db.Database

	windowSize *windowSize
}

func NewModel(bus *bus.TransportBus, db *db.Database) *Model {
	return &Model{
		currentPage: PageMain,
		bus:         bus,
		db:          db,
		windowSize: &windowSize{
			0, 0,
		},
		mainPageModel:  CreateMainPageModel(db),
		tracePageModel: CreateTracePageModel(db),
	}
}

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
	case MessageGoToTrace:
		m.currentPage = PageTrace
		cmds = append(cmds, m.tracePageModel.Init())
	case MessageReturnToMainPage:
		m.currentPage = PageMain
		cmds = append(cmds, m.mainPageModel.Init())
	case MessageResourceSpansArrived:
		cmds = append(
			cmds,
			listenForNewSpans(m.bus.TraceBus),
			m.InsertResourceSpans(msg.resourceSpans),
		)
	}

	switch m.currentPage {
	case PageMain:
		m.mainPageModel, cmd = m.mainPageModel.Update(msg)
		cmds = append(cmds, cmd)
	case PageTrace:
		m.tracePageModel, cmd = m.tracePageModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) InsertResourceSpans(resourceSpans ptrace.ResourceSpans) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	m.db.InsertResourceSpans(resourceSpans)

	m.mainPageModel, cmd = m.mainPageModel.Update(MessageUpdateRootSpanRows{})
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func listenForNewSpans(spanChan chan ptrace.ResourceSpans) tea.Cmd {
	return func() tea.Msg {
		return MessageResourceSpansArrived{resourceSpans: <-spanChan}
	}
}

func (m Model) View() string {
	baseStyles := lipgloss.NewStyle().
		Width(m.windowSize.width).
		Height(m.windowSize.height)

	return baseStyles.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Padding(0, 1).Width(baseStyles.GetWidth()).Background(ColorAccent).Bold(true).Render("Otelly"),
			m.GetPageView(baseStyles.GetWidth(), baseStyles.GetHeight()-1),
		),
	)
}

func (m *Model) GetPageView(w, h int) string {
	switch m.currentPage {
	case PageMain:
		return m.mainPageModel.View(w, h)
	case PageTrace:
		return m.tracePageModel.View(w, h)
	}

	return "unknown page"
}
