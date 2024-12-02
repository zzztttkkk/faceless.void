package vld

import (
	"context"
	"fmt"

	"github.com/zzztttkkk/faceless.void/internal"
)

type intVldOptionsKey int

const (
	intVldOptionsKeyForMinV = intVldOptionsKey(iota)
	intVldOptionsKeyForMaxV
	intVldOptionsKeyForCustom
)

type _IntVldBuilder[T internal.IntType] struct {
	pairs []internal.Pair[intVldOptionsKey]
}

func (opts *_IntVldBuilder[T]) MinValue(v T) *_IntVldBuilder[T] {
	opts.pairs = append(opts.pairs, internal.PairOf(intVldOptionsKeyForMinV, v))
	return opts
}

func (opts *_IntVldBuilder[T]) MaxValue(v T) *_IntVldBuilder[T] {
	opts.pairs = append(opts.pairs, internal.PairOf(intVldOptionsKeyForMaxV, v))
	return opts
}

func (opts *_IntVldBuilder[T]) Custom(fnc func(context.Context, T) error) *_IntVldBuilder[T] {
	opts.pairs = append(opts.pairs, internal.PairOf(intVldOptionsKeyForCustom, fnc))
	return opts
}

var (
	msgForIntValueLTMin = internal.NewI18nString("fv.vld: less than min(%d), %d")
	msgForIntValueGTMax = internal.NewI18nString("fv.vld: greater than max(%d), %d")
	msgForCustomFunc    = internal.NewI18nString("fv.vld: custom function error(%s), %v")
)

func (opts *_IntVldBuilder[T]) Func() func(context.Context, T) error {
	var fncs []func(context.Context, T) error

	var min, max T
	var minok, maxok bool
	for _, pair := range opts.pairs {
		switch pair.Key {
		case intVldOptionsKeyForMinV:
			{
				minv := pair.Val.(T)

				min = minv
				minok = true

				fncs = append(fncs, func(ctx context.Context, t T) error {
					if t < minv {
						return internal.NewError(ctx, ErrorKindVldIntLtMin, msgForIntValueLTMin, minv, t)
					}
					return nil
				})
				break
			}
		case intVldOptionsKeyForMaxV:
			{
				maxv := pair.Val.(T)

				max = maxv
				maxok = true

				fncs = append(fncs, func(ctx context.Context, t T) error {
					if t > maxv {
						return internal.NewError(ctx, ErrorKindVldIntGtMax, msgForIntValueGTMax, maxv, t)
					}
					return nil
				})
			}
		case intVldOptionsKeyForCustom:
			{
				fnc := pair.Val.(func(context.Context, T) error)
				if fnc != nil {
					fncs = append(fncs, func(ctx context.Context, t T) error {
						err := fnc(ctx, t)
						if err != nil {
							return internal.NewError(ctx, ErrorKindVldCustomFunc, msgForCustomFunc, err, t)
						}
						return nil
					})
				}
			}
		}
	}

	if minok && maxok && min > max {
		panic(fmt.Errorf("fv.vld: min(%d) > max(%d)", min, max))
	}

	if len(fncs) < 1 {
		return nil
	}
	return func(ctx context.Context, t T) error {
		for _, fnc := range fncs {
			err := fnc(ctx, t)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Integer[T internal.IntType]() *_IntVldBuilder[T] {
	return &_IntVldBuilder[T]{}
}
