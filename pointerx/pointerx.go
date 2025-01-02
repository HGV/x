package pointerx

// Deprecated: pointerx.Ptr is deprecated. Use x.Ptr instead.
func Ptr[T any](v T) *T {
	return &v
}
