package vld

import (
	"fmt"

	"github.com/zzztttkkk/faceless.void/internal"
	"github.com/zzztttkkk/lion"
)

type intBuilder[T lion.IntType] struct {
	commonBuilder[T, intBuilder[T]]
	unsigned bool
}

func (builder *intBuilder[T]) wrapkey(k string) string {
	if builder.unsigned {
		return fmt.Sprintf("%s.u", k)
	}
	return k
}

func (builder *intBuilder[T]) Min(minv T) *intBuilder[T] {
	builder.pairs = append(builder.pairs, internal.PairOf(builder.wrapkey("minv"), minv))
	return builder
}

func (builder *intBuilder[T]) Max(maxv T) *intBuilder[T] {
	builder.pairs = append(builder.pairs, internal.PairOf(builder.wrapkey("maxv"), maxv))
	return builder
}

func IntMeta[T lion.IntType]() *intBuilder[T] {
	return (&intBuilder[T]{unsigned: lion.IsUnsignedInt[T]()})
}

func Int[T lion.IntType](ptr *T) *intBuilder[T] {
	return IntMeta[T]().updateptr(ptr)
}
