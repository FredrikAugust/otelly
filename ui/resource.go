package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"

	"github.com/fredrikaugust/otelly/db"
)

type ResourceModel struct {
	width int

	db *db.Database

	resource db.Resource

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resourceSpanCountGraphModel.width = m.width
	case MessageResourceReceived:
		m.resource = msg.Resource
	case MessageResourceAggregationReceived:
		m.resourceSpanCountGraphModel.updateGraph(msg.Aggregation)
	}

	m.resourceSpanCountGraphModel, cmd = m.resourceSpanCountGraphModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ResourceModel) getResourceAndResourceAggregation(resourceID string) (ResourceModel, tea.Cmd) {
	getResourceCmd := func() tea.Msg {
		res, err := m.db.GetResource(resourceID)
		if err != nil {
			zap.L().Warn("could not get resource", zap.String("resourceID", resourceID))
			return nil
		}

		return MessageResourceReceived{Resource: *res}
	}

	getAggregationCmd := func() tea.Msg {
		res, err := m.db.SpansPerMinuteForService(resourceID)
		if err != nil {
			zap.L().Warn("could not get resource aggregation", zap.String("resourceID", resourceID))
			return nil
		}

		return MessageResourceAggregationReceived{Aggregation: res}
	}

	return m, tea.Batch(getResourceCmd, getAggregationCmd)
}

func (m ResourceModel) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(m.width).
		MarginTop(1)

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
