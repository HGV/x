package x

import (
	"slices"
	"strings"
)

func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Ptr[T any](v T) *T {
	return &v
}

func PtrEqual[T comparable](a, b *T) bool {
	return a == b || (a != nil && b != nil && *a == *b)
}

func IsHGVEmployee(email string) bool {
	lowerEmail := strings.ToLower(email)
	return slices.ContainsFunc([]string{"@hgv.it", "@hgj.it"}, func(domain string) bool {
		return strings.HasSuffix(lowerEmail, domain)
	})
}
