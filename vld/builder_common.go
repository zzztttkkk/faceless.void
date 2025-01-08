package vld

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
	"github.com/zzztttkkk/lion"
)

type _CommonBuilder[T any, S any] struct {
	ptr   *T
	pairs []internal.Pair[string]
}

func (builder *_CommonBuilder[T, S]) self() *S {
	return (*S)(unsafe.Pointer(builder))
}

func (builder *_CommonBuilder[T, S]) updateptr(ptr *T) *S {
	builder.ptr = ptr
	return builder.self()
}

func (builder *_CommonBuilder[T, S]) set(k string, v any) *S {
	internal.PairsSet(&builder.pairs, k, v)
	return builder.self()
}

func (builder *_CommonBuilder[T, S]) update(k string, v any, update func(prev any) any) *S {
	internal.PairsUpdate(&builder.pairs, k, v, update)
	return builder.self()
}

func (builder *_CommonBuilder[T, S]) optional() *S {
	return builder.set("optional", true)
}

func (builder *_CommonBuilder[T, S]) Func(fnc func(ctx context.Context, v T) error) *S {
	return builder.update("func", fnc, func(prev any) any {
		if prev == nil {
			return func(ctx context.Context, v any) error {
				return fnc(ctx, v.(T))
			}
		}
		prevfn := prev.(func(ctx context.Context, v any) error)
		return func(ctx context.Context, v any) error {
			err := prevfn(ctx, v)
			if err != nil {
				return err
			}
			return fnc(ctx, v.(T))
		}
	})
}

func (builder *_CommonBuilder[T, S]) Build() *VldFieldMeta {
	obj := &VldFieldMeta{
		gotype: lion.Typeof[T](),
	}
	for _, pair := range builder.pairs {
		switch pair.Key {
		case "optional":
			{
				obj.optional = true

			}
		case "func":
			{
				obj._Func = pair.Val.(func(ctx context.Context, v any) error)
			}
		case "regexp":
			{
				obj.regexp = pair.Val.(*regexp.Regexp)
			}
		case "minl":
			{
				obj.minLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "maxl":
			{
				obj.maxLength = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "mins":
			{
				obj.minSize = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "maxs":
			{
				obj.maxSize = sql.Null[int]{V: pair.Val.(int), Valid: true}
			}
		case "minv":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseInt(sv, 10, 64)
				obj.minInt = sql.Null[int64]{Valid: true, V: iv}
			}
		case "maxv":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseInt(sv, 10, 64)
				obj.maxInt = sql.Null[int64]{Valid: true, V: iv}
			}
		case "minv.u":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseUint(sv, 10, 64)
				obj.minUint = sql.Null[uint64]{Valid: true, V: iv}
			}
		case "maxv.u":
			{
				sv := fmt.Sprintf("%v", pair.Val)
				iv, _ := strconv.ParseUint(sv, 10, 64)
				obj.maxUint = sql.Null[uint64]{Valid: true, V: iv}
			}
		case "mintime":
			{
				obj.minTime = sql.Null[time.Time]{V: pair.Val.(time.Time), Valid: true}
			}
		case "maxtime":
			{
				obj.maxTime = sql.Null[time.Time]{V: pair.Val.(time.Time), Valid: true}
			}
		case "stringranges":
			{
				obj.stringRanges = pair.Val.([]string)
			}
		case "key":
			{
				obj.key = pair.Val.(*VldFieldMeta)
			}
		case "ele":
			{
				obj.ele = pair.Val.(*VldFieldMeta)
			}
		case "scheme":
			{
				obj.scheme = pair.Val.(_IScheme)
			}
		}
	}
	return obj
}

func (builder *_CommonBuilder[T, S]) With(ctx context.Context) {
	if ctx == nil {
		return
	}

	sv := ctx.Value(internal.CtxKeyForVldScheme)
	if sv == nil {
		panic(fmt.Errorf("fv.vld: empty scheme"))
	}
	scheme := sv.(_IScheme)
	obj := builder.Build()
	scheme.updatemeta(builder.ptr, obj)
}
