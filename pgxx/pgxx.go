package pgxx

func LikeBegins(s string) string {
	return s + "%"
}

func LikeEnds(s string) string {
	return "%" + s
}

func LikeContains(s string) string {
	return LikeBegins(LikeEnds(s))
}
