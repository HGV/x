package slicesx

import "slices"

func Filter[S ~[]E, E any](s S, f func(E) bool) S {
	return slices.DeleteFunc(slices.Clone(s), func(e E) bool {
		return !f(e)
	})
}
