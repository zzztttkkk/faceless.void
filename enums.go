package fv

import (
	"fmt"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _IEnum interface {
	internal.IntType
	fmt.Stringer
}

type EnumPair struct {
	Name  string
	Value int64
}

func EnumPairs[T _IEnum](min T, max T) []EnumPair {
	var pairs []EnumPair
	for i := min + 1; i < max; i++ {
		pairs = append(pairs, EnumPair{i.String(), int64(i)})
	}
	return pairs
}

func EnumNames[T _IEnum](min T, max T) []string {
	var names []string
	for i := min + 1; i < max; i++ {
		names = append(names, i.String())
	}
	return names
}
