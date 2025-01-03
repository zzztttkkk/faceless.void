package vld

import (
	"context"
	"database/sql"
	"regexp"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type commonBuilder[T any, S any] struct {
	ptr   *T
	pairs []internal.Pair[string]
}

func (builder *commonBuilder[T, S]) self() *S {
	return (*S)(unsafe.Pointer(builder))
}

func (builder *commonBuilder[T, S]) updateptr(ptr *T) *S {
	builder.ptr = ptr
	return builder.self()
}

func (builder *commonBuilder[T, S]) Optional() *S {
	builder.pairs = append(builder.pairs, internal.PairOf("optional", true))
	return builder.self()
}

func (builder *commonBuilder[T, S]) Func(fnc func(ctx context.Context, v T) error) *S {
	builder.pairs = append(builder.pairs, internal.PairOf("func", func(ctx context.Context, val any) error { return fnc(ctx, val.(T)) }))
	return builder.self()
}

func (builder *commonBuilder[T, S]) Build() *VldFieldMeta {
	obj := &VldFieldMeta{}
	for _, pair := range builder.pairs {
		switch pair.Key {
		case "optional":
			obj.Optional = true
		case "func":
			obj.Func = pair.Val.(func(ctx context.Context, v any) error)
		case "regexp":
			obj.Regexp = pair.Val.(*regexp.Regexp)
		case "minl":
			obj.MinLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
		case "maxl":
			obj.MaxLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
		case "stringranges":
			obj.StringRanges = pair.Val.([]string)
		}
	}
	return obj
}

func (builder *commonBuilder[T, S]) Finish(scheme *_Scheme) {
	obj := builder.Build()
	fv := scheme.typeinfo.FieldByUnsafePtr(unsafe.Pointer(builder.ptr))
	fv.UpdateMetainfo(obj)
}
