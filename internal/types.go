package internal

type Empty struct{}

type Set[K comparable] map[K]Empty

type Pair[K comparable] struct {
	Key K
	Val any
}

func PairOf[K comparable](key K, val any) Pair[K] {
	return Pair[K]{key, val}
}
