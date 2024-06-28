package cmpx

func Bool[T ~bool](a, b T) int {
	if !a && b {
		return -1
	}
	if a && !b {
		return +1
	}
	return 0
}
