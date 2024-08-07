package internal

func SliceFilter[T any](src []T, fnc func(int, *T) bool) []T {
	var ret = make([]T, 0, len(src))
	for i := 0; i < len(src); i++ {
		if fnc(i, &src[i]) {
			ret = append(ret, src[i])
		}
	}
	return ret
}

func SliceFilterPtr[T any](src []T, fnc func(int, *T) bool) []*T {
	var ret = make([]*T, 0, len(src))
	for i := 0; i < len(src); i++ {
		ptr := &src[i]
		if fnc(i, ptr) {
			ret = append(ret, ptr)
		}
	}
	return ret
}

func SliceMap[T any, R any](src []T, fnc func(int, *T) R) []R {
	var ret = make([]R, 0, len(src))
	for i := 0; i < len(src); i++ {
		ret = append(ret, fnc(i, &src[i]))
	}
	return ret
}
