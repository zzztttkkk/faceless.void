package fv

func Must[T any](v T, err error) T {
	if err == nil {
		return v
	}
	panic(err)
}

func Or[T any](v T, err error, fnc func(error) T) T {
	if err == nil {
		return v
	}
	return fnc(err)
}
