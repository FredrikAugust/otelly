package helpers_test

import (
	"fmt"
	"testing"

	"github.com/fredrikaugust/otelly/ui/helpers"
)

func TestClamp(t *testing.T) {
	tc := []struct {
		min      int
		max      int
		n        int
		expected int
	}{
		{1, 2, 3, 2},
		{0, 2, 1, 1},
		{-1, 10000, 3, 3},
	}

	for _, c := range tc {
		t.Run(fmt.Sprintf("clamps [%v, %v] with n=%v", c.min, c.max, c.n), func(t *testing.T) {
			if helpers.Clamp(c.min, c.n, c.max) != c.expected {
				t.Fatalf("failed to clamp value: %v", c)
			}
		})
	}
}
