// Package components contain UI elements used throughout the app
package components

import (
	"github.com/charmbracelet/bubbles/table"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func CreateSpanTable() table.Model {
	cols := []table.Column{
		{
			Title: "Name",
			Width: 16,
		},
		{
			Title: "Service",
			Width: 16,
		},
		{
			Title: "Start time",
			Width: 12,
		},
		{
			Title: "Duration",
			Width: 8,
		},
	}

	return table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithStyles(table.DefaultStyles()),
	)
}

func SpanToRow(span ptrace.Span) table.Row {
	return []string{span.SpanID().String(), span.Name(), span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime()).String()}
}
