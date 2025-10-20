package helpers

import "cmp"

func Clamp[T cmp.Ordered](min, n, max T) T {
	if n < min {
		return min
	}

	if n > max {
		return max
	}

	return n
}
