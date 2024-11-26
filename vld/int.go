package vld

import (
	"context"
	"fmt"

	"github.com/zzztttkkk/faceless.void/i18n"
	"github.com/zzztttkkk/faceless.void/internal"
)

type intVldOptionsKey int

const (
	intVldOptionsKeyForMinV = intVldOptionsKey(iota)
	intVldOptionsKeyForMaxV
	intVldOptionsKeyForCustom
)

type _IntVldOptionPair = _VldOptionPair[intVldOptionsKey]

type _IntVldOptions[T internal.IntType] struct {
	pairs []_IntVldOptionPair
}

func (opts *_IntVldOptions[T]) MinValue(v T) *_IntVldOptions[T] {
	opts.pairs = append(opts.pairs, _IntVldOptionPair{intVldOptionsKeyForMinV, v})
	return opts
}

func (opts *_IntVldOptions[T]) MaxValue(v T) *_IntVldOptions[T] {
	opts.pairs = append(opts.pairs, _IntVldOptionPair{intVldOptionsKeyForMaxV, v})
	return opts
}

var (
	msgForIntValueLTMin = i18n.New("fv.vld: less than min(%d), %d")
	msgForIntValueGTMax = i18n.New("fv.vld: greater than max(%d), %d")
	msgForCustomFunc    = i18n.New("fv.vld: custom function error(%s), %v")
)

func (opts *_IntVldOptions[T]) Func() func(context.Context, T) error {
	var fncs []func(context.Context, T) error

	var min, max T
	var minok, maxok bool
	for _, pair := range opts.pairs {
		switch pair.key {
		case intVldOptionsKeyForMinV:
			{
				minv := pair.val.(T)

				min = minv
				minok = true

				fncs = append(fncs, func(ctx context.Context, t T) error {
					if t < minv {
						return newerror(ctx, ErrorKindIntLtMin, msgForIntValueLTMin, minv, t)
					}
					return nil
				})
				break
			}
		case intVldOptionsKeyForMaxV:
			{
				maxv := pair.val.(T)

				max = maxv
				maxok = true

				fncs = append(fncs, func(ctx context.Context, t T) error {
					if t > maxv {
						return newerror(ctx, ErrorKindIntGtMax, msgForIntValueGTMax, maxv, t)
					}
					return nil
				})
			}
		case intVldOptionsKeyForCustom:
			{
				fnc := pair.val.(func(context.Context, T) error)
				if fnc != nil {
					fncs = append(fncs, func(ctx context.Context, t T) error {
						err := fnc(ctx, t)
						if err != nil {
							return newerror(ctx, ErrorKindCustomFunc, msgForCustomFunc, err, t)
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

func Integer[T internal.IntType]() *_IntVldOptions[T] {
	return &_IntVldOptions[T]{}
}
