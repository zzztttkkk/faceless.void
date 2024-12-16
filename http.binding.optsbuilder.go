package fv

import (
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type _CommonBindingOptionsBuilder[T any, Sub any] struct {
	ptr   unsafe.Pointer
	pairs []internal.Pair[string]
}

func (builder *_CommonBindingOptionsBuilder[T, Sub]) self() *Sub {
	return (*Sub)(unsafe.Pointer(builder))
}

func (builder *_CommonBindingOptionsBuilder[T, Sub]) Where(src BindingSrcKind) *Sub {
	builder.pairs = append(builder.pairs, internal.PairOf("where", src))
	return builder.self()
}

func (builder *_CommonBindingOptionsBuilder[T, Sub]) Alias(alias ...string) *Sub {
	builder.pairs = append(builder.pairs, internal.PairOf("alias", alias))
	return builder.self()
}

func (builder *_CommonBindingOptionsBuilder[T, Sub]) Optional() *Sub {
	builder.pairs = append(builder.pairs, internal.PairOf("optional", true))
	return builder.self()
}

func (builder *_CommonBindingOptionsBuilder[T, Sub]) Default(dv T) *Sub {
	builder.pairs = append(builder.pairs, internal.PairOf("default", dv))
	return builder.self()
}

type _StringBindingOptionsBuilder struct {
	_CommonBindingOptionsBuilder[string, _StringBindingOptionsBuilder]
}

func BindingOpts(ptr *string) *_StringBindingOptionsBuilder {
	obj := &_StringBindingOptionsBuilder{}
	obj.ptr = unsafe.Pointer(ptr)
	return obj
}
