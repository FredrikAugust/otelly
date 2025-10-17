package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
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
	}
}

func (m TracePageModel) Init() tea.Cmd {
	return nil
}

func (m *TracePageModel) SetWidth(w int) {
	m.width = w

	m.spanAttributeModel.width = w
	m.traceTreeViewModel.width = int(math.Floor(float64(w)/float64(2))) - 5
	m.waterfallModel.width = int(math.Ceil(float64(w) / float64(2)))
}

func (m TracePageModel) Update(msg tea.Msg) (TracePageModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			cmds = append(cmds, func() tea.Msg { return MessageReturnToMainPage{} })
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
		selectedSpan := m.spans[m.cursor]
		m.spanAttributeModel.SetAttributes(selectedSpan.Attributes)
	}

	m.waterfallModel.cursor = m.cursor
	m.traceTreeViewModel.cursor = m.cursor

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

	return container.Render(
		lipgloss.JoinVertical(0,
			header(m.tree.Item.Span.TraceID, m.tree.Item.Span.Name),
			"",
			lipgloss.JoinHorizontal(0,
				lipgloss.JoinVertical(0,
					m.traceTreeViewModel.View(),
				),
				lipgloss.NewStyle().Width(5).Render(""),
				lipgloss.JoinVertical(0, m.waterfallModel.View()),
			),
			"",
			m.help.View(m),
			"",
			lipgloss.JoinHorizontal(0,
				"resource", // TODO: create a component here
				m.spanAttributeModel.View(),
			),
		),
	)
}

func header(traceID, name string) string {
	return lipgloss.JoinVertical(lipgloss.Top,
		TextTertiary.Render("#"+traceID),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			TextSecondary.Render("Trace "),
			TextHeading.Render(name),
		),
	)
}

func spanView(item helpers.TraceTreeNodeItem, selected bool) string {
	style := lipgloss.
		NewStyle()

	if selected {
		style = style.Background(ColorAccent)
	}

	secondaryText := item.Span.ServiceName + " • " + item.Span.Duration.Round(time.Millisecond).String()
	if item.DurationOfParent != 1 {
		pctOfParentSpan := item.DurationOfParent * 100
		secondaryText += fmt.Sprintf(" (%.1f%%)", pctOfParentSpan)
	}

	return style.Render(item.Span.Name) + " " + TextTertiary.Render(secondaryText)
}
