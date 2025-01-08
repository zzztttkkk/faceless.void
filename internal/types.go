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

func PairsSet[K comparable](pairs *[]Pair[K], key K, val any) {
	found := false
	for idx := range *pairs {
		pp := &((*pairs)[idx])
		if pp.Key == key {
			pp.Val = val
			found = true
			break
		}
	}
	if !found {
		*pairs = append(*pairs, Pair[K]{Key: key, Val: val})
	}
}

func PairsUpdate[K comparable](pairs *[]Pair[K], key K, val any, update func(prev any) any) {
	found := false
	for idx := range *pairs {
		pp := &((*pairs)[idx])
		if pp.Key == key {
			pp.Val = update(pp.Val)
			found = true
			break
		}
	}
	if !found {
		*pairs = append(*pairs, Pair[K]{Key: key, Val: update(nil)})
	}
}
