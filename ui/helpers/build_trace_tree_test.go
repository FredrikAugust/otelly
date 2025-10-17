package helpers_test

import (
	"database/sql"
	"iter"
	"testing"
	"time"

	"github.com/fredrikaugust/otelly/db"
	"github.com/fredrikaugust/otelly/ui/helpers"

	"github.com/stretchr/testify/assert"
)

var testSpans = []db.SpanWithResource{
	{Span: db.Span{ID: "1", ParentSpanID: sql.NullString{Valid: false}, StartTime: time.Now(), Duration: time.Second * 1}},
	{Span: db.Span{ID: "2", ParentSpanID: sql.NullString{String: "1", Valid: true}, StartTime: time.Now().Add(-1 * time.Hour), Duration: time.Hour * 2}},
	{Span: db.Span{ID: "3", ParentSpanID: sql.NullString{String: "2", Valid: true}, StartTime: time.Now(), Duration: time.Minute * 1}},
	{Span: db.Span{ID: "4", ParentSpanID: sql.NullString{String: "1", Valid: true}, StartTime: time.Now(), Duration: time.Minute * 1}},
}

func TestBuildTree(t *testing.T) {
	t.Run("valid tree", func(t *testing.T) {
		tree, err := helpers.BuildTree(testSpans)
		assert.Nil(t, err, "could not build tree")
		assert.Equal(t, tree.Item.Span.ID, testSpans[0].ID)

		assert.Len(t, tree.Children, 2)
		assert.Len(t, tree.Children[0].Children, 1)

		assert.Equal(t, tree.Children[0].Item.Span.ID, testSpans[1].ID)
		assert.Equal(t, tree.Children[0].Children[0].Item.Span.ID, testSpans[2].ID)
		assert.Equal(t, tree.Children[1].Item.Span.ID, testSpans[3].ID)
	})

	t.Run("no root", func(t *testing.T) {
		_, err := helpers.BuildTree([]db.SpanWithResource{
			{Span: db.Span{ID: "3", ParentSpanID: sql.NullString{String: "2", Valid: true}}},
		})
		assert.Error(t, err)
	})
}

func TestGetTimeRange(t *testing.T) {
	t.Run("calculates time range", func(t *testing.T) {
		tree, _ := helpers.BuildTree(testSpans)
		start, end := tree.GetTimeRange()
		assert.Equal(t, testSpans[1].StartTime, start)
		assert.Equal(t, testSpans[1].StartTime.Add(testSpans[1].Duration), end)
	})
}

func TestIterator(t *testing.T) {
	t.Run("iterates over all", func(t *testing.T) {
		tree, _ := helpers.BuildTree(testSpans)

		count := 0
		for range tree.All() {
			count += 1
		}

		assert.Equal(t, 4, count)
	})

	t.Run("iterates in right order", func(t *testing.T) {
		tree, _ := helpers.BuildTree(testSpans)
		iterator := tree.All()
		next, _ := iter.Pull2(iterator)
		depth, item, _ := next()
		assert.Equal(t, testSpans[0].ID, item.Span.ID)
		assert.Equal(t, 0, depth)
		depth, item, _ = next()
		assert.Equal(t, testSpans[1].ID, item.Span.ID)
		assert.Equal(t, 1, depth)
		depth, item, _ = next()
		assert.Equal(t, testSpans[3].ID, item.Span.ID)
		assert.Equal(t, 1, depth)
		depth, item, _ = next()
		assert.Equal(t, testSpans[2].ID, item.Span.ID)
		assert.Equal(t, 2, depth)
		_, _, done := next()
		assert.False(t, done)
		// TODO: check depth is right
	})
}

func TestDurationPct(t *testing.T) {
	t.Run("set duration", func(t *testing.T) {
		tree, _ := helpers.BuildTree(testSpans)
		assert.Equal(t, float64(1), tree.Item.DurationOfParent)
		assert.InEpsilon(t, float64(0.0083333333), tree.Children[0].Children[0].Item.DurationOfParent, 0.0001)
	})
}

func TestParentStartTime(t *testing.T) {
	t.Run("gets right parent start time", func(t *testing.T) {
		tree, _ := helpers.BuildTree(testSpans)
		tree.Item.ParentStartTime = tree.Item.Span.StartTime
		tree.Children[0].Item.ParentStartTime = testSpans[0].StartTime
	})
}
