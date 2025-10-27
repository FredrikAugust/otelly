package flamegraph_test

import (
	"iter"
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
	{"test2", time.Millisecond * 500, "test1", now.Add(time.Millisecond * 250)},
	{"test3", time.Millisecond * 250, "test2", now.Add(time.Millisecond * 500)},
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

	t.Run("builds skinny tree and sets offset pct", func(t *testing.T) {
		root, _ := flamegraph.Build(testItemsSkinny, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				StartTime: t.startTime,
				ParentID:  t.parentID,
			}
		})

		assert.InDelta(t, 0, root.OffsetPct, 0.01)
		assert.InDelta(t, 0.25, root.Children[0].OffsetPct, 0.01)
		assert.InDelta(t, 0.5, root.Children[0].Children[0].OffsetPct, 0.01)
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

		assert.EqualValues(t, 1, root.WidthPct)
		assert.EqualValues(t, 0.5, root.Children[0].WidthPct)
		assert.EqualValues(t, 0.25, root.Children[0].Children[0].WidthPct)
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

		assert.Equal(t, root.Children[0].Name, "test2")
		assert.Equal(t, root.Children[1].Name, "test3")
	})

	t.Run("iterate", func(t *testing.T) {
		root, _ := flamegraph.Build(testItemsComplex, func(t testItem) flamegraph.NodeInput {
			return flamegraph.NodeInput{
				ID:        t.name,
				Name:      t.name,
				Duration:  t.duration,
				StartTime: time.Now(),
				ParentID:  t.parentID,
			}
		})

		next, stop := iter.Pull2(root.All())
		defer stop()

		d, n, valid := next()
		assert.True(t, valid)
		assert.Equal(t, 0, d)
		assert.Equal(t, "test1", n.Name)

		d, n, valid = next()
		assert.True(t, valid)
		assert.Equal(t, 1, d)
		assert.Equal(t, "test2", n.Name)

		d, n, valid = next()
		assert.True(t, valid)
		assert.Equal(t, 1, d)
		assert.Equal(t, "test3", n.Name)

		_, _, valid = next()
		assert.False(t, valid)
	})
}
