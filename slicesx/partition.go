package slicesx

func Partition[S ~[]E, E any](s S, f func(E) bool) (strue S, sfalse S) {
	for _, e := range s {
		if f(e) {
			strue = append(strue, e)
		} else {
			sfalse = append(sfalse, e)
		}
	}
	return
}
