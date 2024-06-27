package cmpx

import "time"

func Time(a, b time.Time) int {
	return int(a.Sub(b))
}

func Bool[T ~bool](a, b T) int {
	if !a && b {
		return -1
	}
	if a && !b {
		return +1
	}
	return 0
}
