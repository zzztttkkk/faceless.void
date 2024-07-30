package fv

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must2[T any, A any](v T, a A, err error) (T, A) {
	if err != nil {
		panic(err)
	}
	return v, a
}

func Must3[T any, A any, B any](v T, a A, b B, err error) (T, A, B) {
	if err != nil {
		panic(err)
	}
	return v, a, b
}
