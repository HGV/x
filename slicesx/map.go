package slicesx

func Map[S ~[]E, E, T any](s S, f func(E) T) []T {
	t := make([]T, len(s))
	for i, e := range s {
		t[i] = f(e)
	}
	return t
}
