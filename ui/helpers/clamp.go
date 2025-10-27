package helpers

import "cmp"

// Clamp clamps a value to a min and max passed in. Both bounds are inclusive.
// Min must be smaller or equal to max.
func Clamp[T cmp.Ordered](min, n, max T) T {
	if n < min {
		return min
	}

	if n > max {
		return max
	}

	return n
}
