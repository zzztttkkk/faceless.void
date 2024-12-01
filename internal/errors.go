package internal

func Must[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func Must2[A any, B any](a A, b B, e error) (A, B) {
	if e != nil {
		panic(e)
	}
	return a, b
}

func Must3[A any, B any, C any](a A, b B, c C, e error) (A, B, C) {
	if e != nil {
		panic(e)
	}
	return a, b, c
}
