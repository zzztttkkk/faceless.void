package vld

import (
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

func (builder *commonBuilder[T, S]) Func(fnc func(v any) error) *S {
	builder.pairs = append(builder.pairs, internal.PairOf("func", fnc))
	return builder.self()
}

func (builder *commonBuilder[T, S]) Build() *VldFieldMeta {
	obj := &VldFieldMeta{}
	for _, pair := range builder.pairs {
		switch pair.Key {
		case "optional":
			obj.Optional = true
		case "func":
			obj.Func = pair.Val.(func(v any) error)

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
