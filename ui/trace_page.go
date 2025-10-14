package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type TracePageModel struct {
	db *db.Database

	spinner spinner.Model

	spans []db.GetSpansForTraceModel

	keyMap []key.Binding
	help   help.Model

	cursor int

	spanAttributeModel SpanAttributeModel
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
		db:      db,
		spinner: spinner.New(spinner.WithSpinner(spinner.Points)),
		cursor:  0,
		keyMap: []key.Binding{
			table.DefaultKeyMap().LineUp,
			table.DefaultKeyMap().LineDown,
			table.DefaultKeyMap().GotoTop,
			table.DefaultKeyMap().GotoBottom,
			key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q", "back")),
		},
		help:               help.New(),
		spanAttributeModel: CreateSpanAttributeModel(),
	}
}

func (m TracePageModel) Init() tea.Cmd {
	return m.spinner.Tick
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
		res, err := m.db.GetSpansForTrace(msg.TraceID)
		if err == nil {
			m.spans = res
		}
	}

	selectedSpan := m.spans[m.cursor]
	m.spanAttributeModel.SetAttributes(selectedSpan.Attributes)

	return m, tea.Batch(cmds...)
}

func (m TracePageModel) View(w, h int) string {
	container := lipgloss.NewStyle().Width(w).Height(h)

	var rootSpan db.GetSpansForTraceModel
	for _, span := range m.spans {
		if !span.ParentSpanID.Valid {
			rootSpan = span
			break
		}
	}

	tree := buildTree(rootSpan, m.spans)
	row := 0

	hierarchicalView := treeView(tree, &row, &m.cursor, nil)
	startTime, endTime, waterfallView := WaterfallLinesForSpans(
		container.GetWidth()-lipgloss.Width(hierarchicalView)-5,
		m.spans,
		func(span *db.GetSpansForTraceModel) string { return "" },
	)

	return container.Render(
		lipgloss.JoinVertical(0,
			header(tree.span.TraceID, tree.span.Name, startTime, endTime),
			"",
			lipgloss.JoinHorizontal(0,
				lipgloss.JoinVertical(0,
					lipgloss.JoinVertical(0,
						hierarchicalView,
					),
				),
				lipgloss.NewStyle().Width(5).Render(""),
				lipgloss.JoinVertical(0.0, waterfallView...),
			),
			m.help.View(m),
			"",
			lipgloss.JoinHorizontal(0,
				"resource",
				m.spanAttributeModel.View(w),
			),
		),
	)
}

func header(traceID, name string, start, end time.Time) string {
	return lipgloss.JoinVertical(lipgloss.Top,
		TextTertiary.Render("#"+traceID),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			TextSecondary.Render("Trace "),
			TextHeading.Render(name),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			TextSecondary.Render(start.Format("2006-01-02 15:04:05")),
			TextTertiary.Render(" — "),
			TextSecondary.Render(end.Format("2006-01-02 15:04:05")),
			TextSecondary.Render(" ("+end.Sub(start).Round(time.Millisecond).String()+")"),
		),
	)
}

type traceNode struct {
	span     db.GetSpansForTraceModel
	children []traceNode
}

func spanView(span db.GetSpansForTraceModel, selected bool, parentSpan *db.GetSpansForTraceModel) string {
	style := lipgloss.
		NewStyle()

	if selected {
		style = style.Background(ColorAccent)
	}

	secondaryText := span.ServiceName + " • " + span.Duration.Round(time.Millisecond).String()
	if parentSpan != nil {
		pctOfParentSpan := (span.Duration.Seconds() / parentSpan.Duration.Seconds()) * 100
		secondaryText += fmt.Sprintf(" (%.1f%%)", pctOfParentSpan)
	}

	return style.Render(span.Name) + " " + TextTertiary.Render(secondaryText)
}

func treeView(tree traceNode, row, cursor *int, parentNode *db.GetSpansForTraceModel) string {
	currentRow := *row
	*row += 1

	strs := make([]string, 0)

	strs = append(strs, spanView(tree.span, currentRow == *cursor, parentNode))

	for _, child := range tree.children {
		strs = append(strs, lipgloss.NewStyle().PaddingLeft(2).Render(treeView(child, row, cursor, &tree.span)))
	}

	return lipgloss.JoinVertical(0,
		strs...,
	)
}

func buildTree(rootSpan db.GetSpansForTraceModel, spans []db.GetSpansForTraceModel) traceNode {
	node := traceNode{}

	children := make([]traceNode, 0)
	for _, span := range spans {
		if span.ParentSpanID.String == rootSpan.ID {
			children = append(children, buildTree(span, spans))
		}
	}

	node.span = rootSpan
	node.children = children

	return node
}
