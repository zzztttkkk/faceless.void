package vld

import (
	"fmt"
	"time"
)

type timeVldOptionsKey int

const (
	timeVldOptionsKeyForBegin = timeVldOptionsKey(iota)
	timeVldOptionsKeyForEnd
	timeVldOptionsKeyForCustom
)

type _TimeVldOptionPair = _VldOptionPair[timeVldOptionsKey]

type _TimeVldOptions struct {
	pairs []_TimeVldOptionPair
}

func (opts *_TimeVldOptions) Begin(v time.Time) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, _TimeVldOptionPair{timeVldOptionsKeyForBegin, v})
	return opts
}

func (opts *_TimeVldOptions) End(v time.Time) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, _TimeVldOptionPair{timeVldOptionsKeyForEnd, v})
	return opts
}

func (opts *_TimeVldOptions) Custom(fnc func(*time.Time) error) *_TimeVldOptions {
	opts.pairs = append(opts.pairs, _TimeVldOptionPair{timeVldOptionsKeyForBegin, fnc})
	return opts
}

func (opts *_TimeVldOptions) Func() func(time.Time) error {
	var fncs []func(time.Time) error

	for _, pair := range opts.pairs {
		switch pair.key {
		case timeVldOptionsKeyForBegin:
			{
				begin := pair.val.(time.Time)
				fncs = append(fncs, func(t time.Time) error {
					if t.Sub(begin) < 0 {
						return fmt.Errorf("")
					}
					return nil
				})
				break
			}
		case timeVldOptionsKeyForEnd:
			{
				end := pair.val.(time.Time)
				fncs = append(fncs, func(t time.Time) error {
					if t.Sub(end) > 0 {
						return fmt.Errorf("")
					}
					return nil
				})
				break
			}
		case timeVldOptionsKeyForCustom:
			{
				fncs = append(fncs, pair.val.(func(time.Time) error))
				break
			}
		}
	}
	if len(fncs) < 1 {
		return nil
	}
	return func(t time.Time) error {
		for _, fnc := range fncs {
			err := fnc(t)
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
