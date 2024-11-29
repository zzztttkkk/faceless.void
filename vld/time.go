package vld

import (
	"context"
	"fmt"
	"time"

	"github.com/zzztttkkk/faceless.void/internal"
)

type timeVldOptionsKey int

const (
	timeVldOptionsKeyForBegin = timeVldOptionsKey(iota)
	timeVldOptionsKeyForEnd
	timeVldOptionsKeyForCustom
)

type _TimeVldOptions struct {
	pairs []internal.Pair[timeVldOptionsKey]
}

func (opts *_TimeVldOptions) Begin(v time.Time) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(timeVldOptionsKeyForBegin, v))
	return opts
}

func (opts *_TimeVldOptions) End(v time.Time) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(timeVldOptionsKeyForEnd, v))
	return opts
}

func (opts *_TimeVldOptions) Custom(fnc func(context.Context, time.Time) error) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, internal.PairOf(timeVldOptionsKeyForBegin, fnc))
	return opts
}

func (opts *_TimeVldOptions) Func() func(context.Context, time.Time) error {
	var fncs []func(context.Context, time.Time) error

	for _, pair := range opts.pairs {
		switch pair.Key {
		case timeVldOptionsKeyForBegin:
			{
				begin := pair.Val.(time.Time)
				fncs = append(fncs, func(ctx context.Context, t time.Time) error {
					if t.Sub(begin) < 0 {
						return fmt.Errorf("")
					}
					return nil
				})
				break
			}
		case timeVldOptionsKeyForEnd:
			{
				end := pair.Val.(time.Time)
				fncs = append(fncs, func(ctx context.Context, t time.Time) error {
					if t.Sub(end) > 0 {
						return fmt.Errorf("")
					}
					return nil
				})
				break
			}
		case timeVldOptionsKeyForCustom:
			{
				fnc := pair.Val.(func(context.Context, time.Time) error)
				fncs = append(fncs, func(ctx context.Context, t time.Time) error {
					err := fnc(ctx, t)
					if err != nil {
						return internal.NewError(ctx, ErrorKindVldCustomFunc, msgForCustomFunc, err)
					}
					return nil
				})
				break
			}
		}
	}
	if len(fncs) < 1 {
		return nil
	}
	return func(ctx context.Context, t time.Time) error {
		for _, fnc := range fncs {
			err := fnc(ctx, t)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Time() *_TimeVldOptions {
	return &_TimeVldOptions{}
}
