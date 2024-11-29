package internal

type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Empty struct{}

type Set[K comparable] map[K]Empty

type Pair[K comparable] struct {
	Key K
	Val any
}

func PairOf[K comparable](key K, val any) Pair[K] {
	return Pair[K]{key, val}
}
