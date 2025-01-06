package vld

import "time"

type _TimeBuilder struct {
	_CommonBuilder[time.Time, _TimeBuilder]
}

func (builder *_TimeBuilder) Earliest(tv time.Time) *_TimeBuilder {
	return builder.set("mintime", tv)
}

func (builder *_TimeBuilder) Latest(tv time.Time) *_TimeBuilder {
	return builder.set("maxtime", tv)
}

func TimeMeta() *_TimeBuilder {
	obj := &_TimeBuilder{}
	return obj
}

func Time(fptr *time.Time) *_TimeBuilder {
	return TimeMeta().updateptr(fptr)
}
