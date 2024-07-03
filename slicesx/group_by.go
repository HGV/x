package slicesx

func GroupByFunc[K comparable, S ~[]E, E any](s S, key func(E) K) map[K]S {
	m := make(map[K]S)
	for _, e := range s {
		key := key(e)
		m[key] = append(m[key], e)
	}
	return m
}
