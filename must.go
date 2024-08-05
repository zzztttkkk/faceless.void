package fv

func Must[T any](v T, err error) T {
	if err == nil {
		return v
	}
	panic(err)
}

func Must2[T any, A any](v T, a A, err error) (T, A) {
	if err == nil {
		return v, a
	}
	panic(err)
}

func Or[T any](v T, err error, fnc func(error) T) T {
	if err == nil {
		return v
	}
	return fnc(err)
}
