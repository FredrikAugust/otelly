// Package helpers contains helper functions for the UI
package helpers

import (
	"errors"
	"iter"
	"time"

	"github.com/fredrikaugust/otelly/db"
)

type TraceTreeNodeItem struct {
	Span db.SpanWithResource
	// DurationOfParent is a float 0..1 which is how much of the parent duration this span constitutes
	DurationOfParent float64
	ParentStartTime  time.Time
}

type TraceTreeNode struct {
	Item     TraceTreeNodeItem
	Children []TraceTreeNode
}

func (t TraceTreeNode) All() iter.Seq2[int, TraceTreeNodeItem] {
	type queueItem struct {
		item TraceTreeNode

		depth int
	}

	return func(yield func(int, TraceTreeNodeItem) bool) {
		queue := []queueItem{{item: t, depth: 0}}

		for len(queue) != 0 {
			head := queue[0]
			if !yield(head.depth, head.item.Item) {
				return
			}
			queue = queue[1:] // Pop
			for _, child := range head.item.Children {
				queue = append(queue, queueItem{item: child, depth: head.depth + 1})
			}
		}
	}
}

func (t *TraceTreeNode) GetTimeRange() (time.Time, time.Time) {
	var (
		startTime = t.Item.Span.StartTime
		endTime   = t.Item.Span.StartTime.Add(t.Item.Span.Duration)
	)

	for _, child := range t.Children {
		childStart, childEnd := child.GetTimeRange()
		if childStart.Before(startTime) {
			startTime = childStart
		}
		if childEnd.After(endTime) {
			endTime = childEnd
		}
	}

	return startTime, endTime
}

func BuildTree(spans []db.SpanWithResource) (TraceTreeNode, error) {
	var rootSpan db.SpanWithResource

	for _, span := range spans {
		if !span.ParentSpanID.Valid {
			rootSpan = span
			break
		}
	}

	if rootSpan.ID == "" {
		return TraceTreeNode{}, errors.New("could not find root span in spans passed to build tree")
	}

	return populateRoot(db.SpanWithResource{}, rootSpan, spans), nil
}

func populateRoot(parent, root db.SpanWithResource, spans []db.SpanWithResource) TraceTreeNode {
	node := TraceTreeNode{}
	if parent.ID == "" {
		node.Item.DurationOfParent = 1
		node.Item.ParentStartTime = root.StartTime
	} else {
		node.Item.DurationOfParent = float64(root.Duration) / float64(parent.Duration)
		node.Item.ParentStartTime = parent.StartTime
	}

	children := make([]TraceTreeNode, 0)
	for _, span := range spans {
		if span.ParentSpanID.String == root.ID {
			children = append(children, populateRoot(root, span, spans))
		}
	}

	node.Item.Span = root
	node.Children = children

	return node
}
