package vld

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
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
			{
				obj.Optional = true

			}
		case "func":
			{
				obj.Func = pair.Val.(func(ctx context.Context, v any) error)
			}
		case "regexp":
			{
				obj.Regexp = pair.Val.(*regexp.Regexp)
			}
		case "minl":
			{
				obj.MinLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "maxl":
			{
				obj.MaxLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "mins":
			{
				obj.MinSize = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "maxs":
			{
				obj.MaxSize = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "minv":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseInt(sv, 10, 64)
				obj.MinInt = sql.Null[int64]{Valid: true, V: iv}
			}
		case "maxv":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseInt(sv, 10, 64)
				obj.MaxInt = sql.Null[int64]{Valid: true, V: iv}
			}
		case "minv.u":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseUint(sv, 10, 64)
				obj.MinUint = sql.Null[uint64]{Valid: true, V: iv}
			}
		case "maxv.u":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseUint(sv, 10, 64)
				obj.MaxUint = sql.Null[uint64]{Valid: true, V: iv}
			}
		case "stringranges":
			{
				obj.StringRanges = pair.Val.([]string)
			}
		case "key":
			{
				obj.Key = pair.Val.(*VldFieldMeta)
			}
		case "ele":
			{
				obj.Ele = pair.Val.(*VldFieldMeta)
			}
		}
	}
	return obj
}

func (builder *commonBuilder[T, S]) With(ctx context.Context) {
	sv := ctx.Value(internal.CtxKeyForVldScheme)
	if sv == nil {
		panic(fmt.Errorf("fv.vld: empty scheme"))
	}
	scheme := sv.(_IScheme)
	obj := builder.Build()
	scheme.updatemeta(builder.ptr, obj)
}
