package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"
	"go.uber.org/zap"
)

type TracePageModel struct {
	db *db.Database

	spans []db.SpanWithResource
	tree  helpers.TraceTreeNode

	keyMap []key.Binding
	help   help.Model

	cursor int

	spanAttributeModel SpanAttributeModel
	waterfallModel     SpanWaterfallModel
	traceTreeViewModel TraceTreeViewModel

	viewportModel viewport.Model

	width  int
	height int
}

// FullHelp implements help.KeyMap.
func (m TracePageModel) FullHelp() [][]key.Binding {
	return nil
}

// ShortHelp implements help.KeyMap.
func (m TracePageModel) ShortHelp() []key.Binding {
	return m.keyMap
}

func CreateTracePageModel(db *db.Database) TracePageModel {
	return TracePageModel{
		db:     db,
		cursor: 0,
		keyMap: []key.Binding{
			table.DefaultKeyMap().LineUp,
			table.DefaultKeyMap().LineDown,
			table.DefaultKeyMap().GotoTop,
			table.DefaultKeyMap().GotoBottom,
			key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q", "back")),
		},
		help:               help.New(),
		spanAttributeModel: CreateSpanAttributeModel(),
		waterfallModel:     CreateSpanWaterfallModel(),
		traceTreeViewModel: CreateTraceTreeViewModel(),
		viewportModel:      viewport.New(0, 0),
	}
}

func (m TracePageModel) Init() tea.Cmd {
	return nil
}

func (m *TracePageModel) SetWidth(w int) {
	m.width = w

	m.spanAttributeModel.width = w - 2

	m.viewportModel.Width = w - 2 // borders

	m.traceTreeViewModel.width = int(math.Floor(float64(m.viewportModel.Width)/float64(2))) - 5
	m.waterfallModel.width = int(math.Ceil(float64(m.viewportModel.Width) / float64(2)))
}

func (m *TracePageModel) SetHeight(h int) {
	m.height = h

	m.viewportModel.Height = h - 2 - 10 - 2 - 1 // viewport - header - spanattributes - borders - help view
	m.waterfallModel.height = 999999            // we don't want to limit it

	m.spanAttributeModel.height = 10 - 2 // borders
}

func (m TracePageModel) Update(msg tea.Msg) (TracePageModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			cmds = append(cmds, func() tea.Msg { return MessageGoToMainPage{} })
		case "j", "down":
			if m.cursor != len(m.spans)-1 {
				m.cursor += 1
			}
		case "k", "up":
			if m.cursor != 0 {
				m.cursor -= 1
			}
		case "g":
			m.cursor = 0
		case "G":
			m.cursor = len(m.spans) - 1
		}
	case MessageGoToTrace:
		m.cursor = 0
		m.viewportModel.YOffset = 0
		cmds = append(cmds, m.getTrace(msg.TraceID))
	case MessageReceivedTraceSpans:
		m.spans = msg.Spans
		tree, err := helpers.BuildTree(msg.Spans)
		if err != nil {
			zap.L().Warn("could not build tree", zap.Error(err))
		} else {
			m.tree = tree
			m.waterfallModel.tree = tree
			m.traceTreeViewModel.tree = tree
		}
	}

	m.updateSpanAttributes()

	// We purposefully don't call out to viewportModel.Update here
	// as we want to have a scroll where you don't scroll down unless
	// you're at the bottom of the viewport.
	scrollPadding := 1
	if m.cursor >= m.viewportModel.YOffset+m.viewportModel.Height-1-scrollPadding {
		m.viewportModel.SetYOffset(m.cursor - m.viewportModel.Height + 1 + scrollPadding)
	}
	if m.cursor < m.viewportModel.YOffset+scrollPadding {
		m.viewportModel.SetYOffset(m.cursor - scrollPadding)
	}

	m.waterfallModel.cursor = m.cursor
	m.traceTreeViewModel.cursor = m.cursor

	m.viewportModel.SetContent(
		lipgloss.JoinHorizontal(0,
			lipgloss.JoinVertical(0,
				m.traceTreeViewModel.View(),
			),
			lipgloss.NewStyle().Width(5).Render(""),
			lipgloss.JoinVertical(0, m.waterfallModel.View()),
		),
	)

	return m, tea.Batch(cmds...)
}

func (m TracePageModel) getTrace(traceID string) tea.Cmd {
	return func() tea.Msg {
		res, err := m.db.GetSpansForTrace(traceID)
		if err != nil {
			zap.L().Warn("failed to get trace", zap.Error(err), zap.String("traceID", traceID))
			return nil
		}

		return MessageReceivedTraceSpans{Spans: res}
	}
}

func (m TracePageModel) View() string {
	if len(m.spans) == 0 {
		// TODO: replace with spinner
		return "no spans loaded for trace"
	}

	container := lipgloss.NewStyle().Width(m.width).Height(m.height)

	helpView := m.help.View(m)
	scrollPosView := TextSecondary.Render(fmt.Sprintf("%v / %v", m.cursor+1, m.viewportModel.TotalLineCount()))

	return container.Render(
		lipgloss.JoinVertical(0,
			m.headerView(),
			lipgloss.NewStyle().Width(m.viewportModel.Width).Border(lipgloss.NormalBorder()).BorderForeground(ColorBorderForeground).Render(
				m.viewportModel.View(),
				lipgloss.JoinHorizontal(0,
					helpView,
					strings.Repeat(" ", m.viewportModel.Width-lipgloss.Width(helpView)-lipgloss.Width(scrollPosView)),
					scrollPosView,
				),
			),
			lipgloss.NewStyle().Width(m.spanAttributeModel.width).Height(m.spanAttributeModel.height).Border(lipgloss.NormalBorder()).BorderForeground(ColorBorderForeground).Render(m.spanAttributeModel.View()),
		),
	)
}

func (m *TracePageModel) headerView() string {
	traceID := m.tree.Item.Span.TraceID
	name := m.tree.Item.Span.Name

	return lipgloss.JoinVertical(lipgloss.Top,
		TextTertiary.Render("#"+traceID),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			TextSecondary.Render("Trace "),
			TextHeading.Render(name),
		),
	)
}

// updateSpanAttributes gets the span under the cursor and
// sets the span attributes model to that span's attributes.
//
// If no span is under cursor it no-ops.
func (m *TracePageModel) updateSpanAttributes() {
	if len(m.spans) == 0 {
		return
	}

	selectedSpan := m.spans[m.cursor]
	m.spanAttributeModel.SetAttributes(selectedSpan.Attributes)
}
