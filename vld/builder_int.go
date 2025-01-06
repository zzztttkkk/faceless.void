package vld

import (
	"fmt"

	"github.com/zzztttkkk/lion"
)

type _IntBuilder[T lion.IntType] struct {
	_CommonBuilder[T, _IntBuilder[T]]
	unsigned bool
}

func (builder *_IntBuilder[T]) wrapkey(k string) string {
	if builder.unsigned {
		return fmt.Sprintf("%s.u", k)
	}
	return k
}

func (builder *_IntBuilder[T]) Min(minv T) *_IntBuilder[T] {
	return builder.set(builder.wrapkey("minv"), minv)
}

func (builder *_IntBuilder[T]) Max(maxv T) *_IntBuilder[T] {
	return builder.set(builder.wrapkey("maxv"), maxv)
}

func IntMeta[T lion.IntType]() *_IntBuilder[T] {
	return (&_IntBuilder[T]{unsigned: lion.IsUnsignedInt[T]()})
}

func Int[T lion.IntType](ptr *T) *_IntBuilder[T] {
	return IntMeta[T]().updateptr(ptr)
}
