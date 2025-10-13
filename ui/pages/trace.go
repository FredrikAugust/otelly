package pages

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/components"
)

type TracePageModel struct {
	db *db.Database

	traceID string

	spinner spinner.Model

	spans []db.GetSpansForTraceModel
}

func CreateTracePageModel(db *db.Database) *TracePageModel {
	return &TracePageModel{
		db:      db,
		spinner: spinner.New(spinner.WithSpinner(spinner.Points)),
	}
}

func (m TracePageModel) Init() tea.Cmd {
	return m.spinner.Tick
}

type MessageReturnToMainPage struct{}

func (m TracePageModel) Update(msg tea.Msg) (TracePageModel, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmds = append(cmds, func() tea.Msg { return MessageReturnToMainPage{} })
		}
	case components.MessageGoToTrace:
		res, err := m.db.GetSpansForTrace(msg.TraceID)
		if err == nil {
			m.spans = res
		}
	}

	return m, tea.Batch(cmds...)
}

func (m TracePageModel) View(w, h int) string {
	container := lipgloss.NewStyle().Width(w-2).Padding(0, 1)

	// var rootSpan db.GetSpansForTraceModel
	strs := make([]string, 0)
	for _, span := range m.spans {
		strs = append(strs, span.Name)
		// TODO: get in model and check for no parent here
	}

	return container.Render(
		lipgloss.JoinVertical(
			lipgloss.Left, strs...,
		),
	)
}
