package ui_test

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fredrikaugust/otelly/ui"
)

func TestTable(t *testing.T) {
	t.Run("shows one item", func(t *testing.T) {
		table := ui.NewTableModel()
		table.SetHeight(10)
		table.SetWidth(100)

		d := ui.NewDefaultTableItemDelegate()
		d.ContentFn = func() []string { return []string{"my item", "dogs", "xy", "hounds"} }
		table.SetItems([]ui.TableItemDelegate{d})

		view := table.View()

		if !strings.Contains(view, "my item") {
			t.Fatalf("could not find item rendered: %v", view)
		}

		if !strings.Contains(view, "hounds") {
			t.Fatalf("could not find item rendered")
		}
	})

	t.Run("truncates", func(t *testing.T) {
		table := ui.NewTableModel()
		table.SetHeight(10)
		table.SetWidth(8) // should leave 2 width for each cell
		d := ui.NewDefaultTableItemDelegate()
		d.ContentFn = func() []string { return []string{"my item", "dogs", "xy", "hounds"} }
		table.SetItems([]ui.TableItemDelegate{d})

		view := table.View()

		if strings.Contains(view, "my item") {
			t.Fatalf("item should be truncated")
		}

		if strings.Contains(view, "hounds") {
			t.Fatalf("item should be truncated")
		}

		if !strings.Contains(view, "xy") {
			t.Fatalf("item should not be truncated")
		}
	})
}

func TestTable_ItemViewsInViewport(t *testing.T) {
	t.Run("all are shown", func(t *testing.T) {
		table := ui.NewTableModel()
		table.SetHeight(10) // means 8 rows for content
		table.SetWidth(100)
		table.SetRowHeight(2)
		d := ui.NewDefaultTableItemDelegate()
		d.ContentFn = func() []string {
			return []string{
				"my text",
			}
		}
		table.SetItems([]ui.TableItemDelegate{
			d,
			d,
			d,
			d,
		})

		view := table.View()

		re, _ := regexp.Compile("my text")
		matches := re.FindAllString(view, 100)

		if len(matches) != 4 {
			t.Fatalf("only found %v/4 matches in %v", len(matches), view)
		}
	})

	t.Run("don't show if not in viewport", func(t *testing.T) {
		table := ui.NewTableModel()
		table.SetHeight(10) // means 8 rows for content
		table.SetWidth(100)
		table.SetRowHeight(4)
		d := ui.NewDefaultTableItemDelegate()
		d.ContentFn = func() []string {
			return []string{
				"my text",
			}
		}
		table.SetItems([]ui.TableItemDelegate{
			d,
			d,
			d,
			d,
		})

		view := table.View()

		re, _ := regexp.Compile("my text")
		matches := re.FindAllString(view, 100)

		if len(matches) != 2 {
			t.Fatalf("only found %v/2 matches in %v", len(matches), view)
		}
	})
}

func TestTable_Navigation(t *testing.T) {
	t.Run("navigate up and down the table", func(t *testing.T) {
		table := ui.NewTableModel()
		table.SetHeight(4) // means 2 rows for content
		table.SetWidth(20)
		table.SetRowHeight(1)
		d1 := ui.NewDefaultTableItemDelegate()
		d1.ContentFn = func() []string { return []string{"string1"} }
		d2 := ui.NewDefaultTableItemDelegate()
		d2.ContentFn = func() []string { return []string{"string2"} }
		d3 := ui.NewDefaultTableItemDelegate()
		d3.ContentFn = func() []string { return []string{"string3"} }
		d4 := ui.NewDefaultTableItemDelegate()
		d4.ContentFn = func() []string { return []string{"string4"} }
		table.SetItems([]ui.TableItemDelegate{
			d1,
			d2,
			d3,
			d4,
		})

		tc := []struct {
			key            tea.KeyMsg
			expectExist    []string
			expectNotExist []string
		}{
			{tea.KeyMsg{Type: tea.KeyUp}, []string{"string1", "string2"}, []string{"string3", "string4"}},
			{tea.KeyMsg{Type: tea.KeyDown}, []string{"string1", "string2"}, []string{"string3", "string4"}},
			{tea.KeyMsg{Type: tea.KeyDown}, []string{"string2", "string3"}, []string{"string1", "string4"}},
			{tea.KeyMsg{Type: tea.KeyDown}, []string{"string3", "string4"}, []string{"string1", "string2"}},
			{tea.KeyMsg{Type: tea.KeyUp}, []string{"string3", "string4"}, []string{"string1", "string2"}},
			{tea.KeyMsg{Type: tea.KeyUp}, []string{"string2", "string3"}, []string{"string1", "string4"}},
			{tea.KeyMsg{Type: tea.KeyUp}, []string{"string1", "string2"}, []string{"string3", "string4"}},
			{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}, []string{"string3", "string4"}, []string{"string1", "string2"}},
			{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}, []string{"string1", "string2"}, []string{"string3", "string4"}},
		}

		var view string

		for _, c := range tc {
			table, _ = table.Update(c.key)
			view = table.View()

			for _, s := range c.expectExist {
				if !strings.Contains(view, s) {
					t.Fatalf("can't find string %v in %v", s, view)
				}
			}

			for _, s := range c.expectNotExist {
				if strings.Contains(view, s) {
					t.Fatalf("shouldn't find string %v in %v", s, view)
				}
			}
		}
	})
}
