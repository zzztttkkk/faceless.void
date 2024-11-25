package vld

import (
	"fmt"

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

func (opts *_IntVldOptions[T]) Func() func(T) error {
	var fncs []func(T) error

	for _, pair := range opts.pairs {
		switch pair.key {
		case intVldOptionsKeyForMinV:
			{
				minv := pair.val.(T)
				fncs = append(fncs, func(t T) error {
					if t < minv {
						return fmt.Errorf("< minv(%d)", minv)
					}
					return nil
				})
				break
			}
		case intVldOptionsKeyForMaxV:
			{
				maxv := pair.val.(T)
				fncs = append(fncs, func(t T) error {
					if t > maxv {
						return fmt.Errorf("> maxv(%d)", maxv)
					}
					return nil
				})
			}
		case intVldOptionsKeyForCustom:
			{
				fnc := pair.val.(func(T) error)
				if fnc != nil {
					fncs = append(fncs, fnc)
				}
			}
		}
	}

	if len(fncs) < 1 {
		return nil
	}
	return func(t T) error {
		for _, fnc := range fncs {
			err := fnc(t)
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
