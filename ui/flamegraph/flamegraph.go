// Package flamegraph is meant to help build tree like structures
// for showing flamegraphs
package flamegraph

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type Node struct {
	Name      string
	StartTime time.Time
	Duration  time.Duration
	// WidthPct is a number 0 to 1 which is the width of the total time frame it should take up
	WidthPct float64
	Children []Node
}

type NodeInput struct {
	ID        string
	Name      string
	Duration  time.Duration
	ParentID  string
	StartTime time.Time
}

func Build[T any](items []T, retriever func(T) NodeInput) (Node, error) {
	if len(items) == 0 {
		return Node{}, errors.New("can't build flamegraph from empty items")
	}

	var start, end time.Time

	nis := make([]NodeInput, len(items))
	for i, item := range items {
		nis[i] = retriever(item)

		niStart, niEnd := nis[i].StartTime, nis[i].StartTime.Add(nis[i].Duration)

		if start.IsZero() {
			start, end = niStart, niEnd
		} else {
			if niStart.Before(start) {
				start = niStart
			}

			if niEnd.After(end) {
				end = niEnd
			}
		}
	}

	slices.SortFunc(nis, func(a, b NodeInput) int {
		return b.StartTime.Compare(a.StartTime)
	})

	roots := findNodesBelongingToParent(nil, nis, start, end)

	if len(roots) != 1 {
		return Node{}, fmt.Errorf("when building node either 0 or more than 1 root was found: len(roots)=%v", len(roots))
	}

	return roots[0], nil
}

// findNodesBelongingToParent takes in a reference to the parent you want
// to find the children for, and all the nodes you want to look through, and
// returns a slice of the parent's children.
func findNodesBelongingToParent(parent *NodeInput, inputs []NodeInput, start, end time.Time) []Node {
	results := make([]Node, 0)

	for _, input := range inputs {
		if parent == nil {
			// This means we're looking for the root node
			if input.ParentID == "" {
				return []Node{
					{
						Name:      input.Name,
						Duration:  input.Duration,
						StartTime: input.StartTime,
						WidthPct:  float64(input.Duration) / float64(end.Sub(start)),
						Children:  findNodesBelongingToParent(&input, inputs, start, end),
					},
				}
			} else {
				continue
			}
		}

		if input.ParentID == parent.ID {
			results = append(results, Node{
				Name:      input.Name,
				Duration:  input.Duration,
				StartTime: input.StartTime,
				WidthPct:  float64(input.Duration) / float64(end.Sub(start)),
				Children:  findNodesBelongingToParent(&input, inputs, start, end),
			})
		}
	}

	return results
}
