package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/fredrikaugust/otelly/db"
)

type ResourceModel struct {
	width int

	db *db.Database

	resource *db.Resource

	resourceSpanCountGraphModel ResourceSpanCountGraphModel
}

func CreateResourceModel(db *db.Database) ResourceModel {
	return ResourceModel{
		width:                       10,
		db:                          db,
		resourceSpanCountGraphModel: CreateResourceSpanCountGraphModel(db),
	}
}

func (m ResourceModel) Init() tea.Cmd {
	return m.resourceSpanCountGraphModel.Init()
}

func (m ResourceModel) Update(msg tea.Msg) (ResourceModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.resourceSpanCountGraphModel.width = m.width
	case MessageSetSelectedSpan:
		m.resourceSpanCountGraphModel.resourceID = m.resource.ID
	}

	m.resourceSpanCountGraphModel, cmd = m.resourceSpanCountGraphModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ResourceModel) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(m.width).
		MarginTop(1)

	if m.resource == nil {
		return baseStyle.Render("No resource found")
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Resource"),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top,
					TextSecondary.Render("Service "),
					m.resource.ServiceNamespace,
					".",
					m.resource.ServiceName,
				),
				lipgloss.NewStyle().
					MarginTop(1).
					Render(
						m.resourceSpanCountGraphModel.View(),
					),
			),
		),
	)
}
