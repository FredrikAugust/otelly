package components

import (
	"log/slog"
	"reflect"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fredrikaugust/otelly/db"
)

type SpanDetailsModel struct {
	span     *db.Span
	resource *db.Resource

	width  int
	height int

	db *db.Database
}

// Init implements tea.Model.
func (s SpanDetailsModel) Init() tea.Cmd {
	return nil
}

func CreateSpanDetailsModel(db *db.Database) *SpanDetailsModel {
	return &SpanDetailsModel{
		span:     nil,
		resource: nil,

		db: db,

		width:  0,
		height: 0,
	}
}

type MessageSetSelectedSpan struct {
	SpanID string
}

func (s SpanDetailsModel) Update(msg tea.Msg) (SpanDetailsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case MessageSetSelectedSpan:
		if s.span != nil && msg.SpanID == s.span.ID {
			// Don't query the span again when we already have it
			break
		}

		span, err := s.db.GetSpan(msg.SpanID)
		if err != nil {
			slog.Warn("could not get span with resource in span details", "spanID", msg.SpanID, "error", err)
		} else {
			s.span = span
			res, err := s.db.GetResource(span.ResourceID)
			if err != nil {
				slog.Warn("could not get resource", "error", err)
			} else {
				s.resource = res
			}
		}
	}
	return s, nil
}

func (s SpanDetailsModel) View() string {
	box := lipgloss.NewStyle().
		Width(s.width-2).
		Height(s.height-2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	if s.span == nil {
		return box.Foreground(lipgloss.Color("#afafb2")).Align(lipgloss.Center, lipgloss.Center).Render("No span selected")
	}

	return box.
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#6b6b6e")).Render("Span "),
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Render(s.span.ID),
			),
			lipgloss.NewStyle().Bold(true).Render(s.span.Name),
			s.attributeView(),
			s.resourceView(),
		))
}

func (s SpanDetailsModel) resourceView() string {
	baseStyle := lipgloss.NewStyle().
		Width(s.width-6).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	if s.resource == nil {
		return baseStyle.Render("No resource found")
	}

	return baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Resource"),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Render("Service "),
					s.resource.ServiceNamespace,
					".",
					s.resource.ServiceName,
				),
			),
		),
	)
}

func (s SpanDetailsModel) attributeView() string {
	attributeStrs := make([]string, 0)

	keys := make([]string, 0)
	for k := range s.span.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := s.span.Attributes[k]

		// Only show strings
		if reflect.TypeOf(v) != reflect.TypeFor[string]() {
			continue
		}

		attributeStrs = append(
			attributeStrs,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.NewStyle().Foreground(lipgloss.Color("#afafb2")).Width(16).Render(k),
				" ",
				lipgloss.NewStyle().Render(v.(string)),
			),
		)
	}

	return lipgloss.NewStyle().
		Width(s.width-6).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render("Attributes"),
				lipgloss.JoinVertical(lipgloss.Left, attributeStrs...),
			),
		)
}

func (s *SpanDetailsModel) SetWidth(w int) {
	s.width = w
}

func (s *SpanDetailsModel) SetHeight(h int) {
	s.height = h
}
