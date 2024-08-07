package fv

import "sync"

// LazyFnc
// to avoid the problem of circular dependency
type LazyFnc[In any, Out any] struct {
	fnc  func(in In) Out
	once sync.Once
}

func NewHole[In any, Out any]() *LazyFnc[In, Out] {
	return &LazyFnc[In, Out]{}
}

func (hole *LazyFnc[In, Out]) Fill(fnc func(in In) Out) {
	hole.once.Do(func() {
		hole.fnc = fnc
	})
}

func (hole *LazyFnc[In, Out]) Exec(in In) Out {
	return hole.fnc(in)
}
