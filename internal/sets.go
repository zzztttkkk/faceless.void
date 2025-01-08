package internal

func MakeSet[T comparable](vals []T) map[T]struct{} {
	mv := map[T]struct{}{}
	for _, v := range vals {
		mv[v] = struct{}{}
	}
	return mv
}
