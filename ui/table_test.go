package ui_test

import (
	"regexp"
	"strings"
	"testing"

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

		if len(table.ItemViewsInViewport()) != 2 {
			t.Fatalf("table calculates wrong number of visible items: %v", len(table.ItemViewsInViewport()))
		}

		if len(matches) != 2 {
			t.Fatalf("only found %v/2 matches in %v", len(matches), view)
		}
	})
}
