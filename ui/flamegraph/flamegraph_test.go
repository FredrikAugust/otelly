package flamegraph_test

import (
	"testing"
	"time"

	"github.com/fredrikaugust/otelly/ui/flamegraph"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	name      string
	duration  time.Duration
	parentID  string
	startTime time.Time
}

var now = time.Now()

var testItemsSkinny = []testItem{
	{"test1", time.Second, "", now},
	{"test2", time.Millisecond * 500, "test1", now},
	{"test3", time.Millisecond * 250, "test2", now},
}

var testItemsComplex = []testItem{
	{"test1", time.Second, "", now},
	{"test2", time.Millisecond * 500, "test1", now.Add(-1 * time.Minute)},
	{"test3", time.Millisecond * 250, "test1", now.Add(5 * time.Minute)},
}

func TestFlamegraph_Build(t *testing.T) {
	t.Run("builds a flamegraph", func(t *testing.T) {
		root, err := flamegraph.Build(testItemsSkinny, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				StartTime: time.Now(),
				ParentID:  t.parentID,
			}
		})

		assert.NotNil(t, root)
		assert.Nil(t, err)
	})

	t.Run("fails to build empty graph", func(t *testing.T) {
		_, err := flamegraph.Build([]testItem{}, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{}
		})

		assert.ErrorContains(t, err, "can't build flamegraph from empty items")
	})

	t.Run("fails to build graph with no root", func(t *testing.T) {
		_, err := flamegraph.Build([]testItem{
			{
				name:      "test",
				duration:  10 * time.Second,
				parentID:  "dog",
				startTime: time.Now(),
			},
		}, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				ParentID:  t.parentID,
				StartTime: t.startTime,
			}
		})

		assert.ErrorContains(t, err, "len(roots)=0")
	})

	t.Run("builds skinny tree and sets correct pct", func(t *testing.T) {
		root, err := flamegraph.Build(testItemsSkinny, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				StartTime: time.Now(),
				ParentID:  t.parentID,
			}
		})

		assert.Nil(t, err)

		assert.Len(t, root.Children, 1)
		assert.Len(t, root.Children[0].Children, 1)
		assert.Len(t, root.Children[0].Children[0].Children, 0)

		assert.EqualValues(t, root.WidthPct, 1)
		assert.EqualValues(t, root.Children[0].WidthPct, 0.5)
		assert.EqualValues(t, root.Children[0].Children[0].WidthPct, 0.25)
	})

	t.Run("sorts the entries correctly", func(t *testing.T) {
		root, err := flamegraph.Build(testItemsComplex, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				StartTime: time.Now(),
				ParentID:  t.parentID,
			}
		})

		assert.Nil(t, err)

		assert.Equal(t, root.Children[0].Name, "test3")
		assert.Equal(t, root.Children[1].Name, "test2")
	})
}
