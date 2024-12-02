package fv

import "sync"

// _LazyFnc
// to avoid the problem of circular dependency
type _LazyFnc[In any, Out any] struct {
	fnc  func(In) Out
	once sync.Once
}

func LazyFunc[In any, Out any]() *_LazyFnc[In, Out] {
	return &_LazyFnc[In, Out]{}
}

func (hole *_LazyFnc[In, Out]) Fill(fnc func(In) Out) {
	hole.once.Do(func() {
		hole.fnc = fnc
	})
}

func (hole *_LazyFnc[In, Out]) Exec(in In) Out {
	return hole.fnc(in)
}
