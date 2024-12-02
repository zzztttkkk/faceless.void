package fv

import (
	"reflect"
)

type Dalgate[T any] struct {
	fnc T
	ok  bool
}

func (delgate *Dalgate[T]) Set(fnc T) {
	delgate.fnc = fnc
	fv := reflect.ValueOf(fnc)
	if fv.IsValid() && !fv.IsZero() && !fv.IsNil() {
		delgate.ok = true
	}
}

func (delgate *Dalgate[T]) Func() T {
	if !delgate.ok {
		panic("")
	}
	return delgate.fnc
}
