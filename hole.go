package fv

import "sync"

// to avoid the problem of circular dependency
type Hole[In any, Out any] struct {
	fnc  func(in In) Out
	once sync.Once
}

func NewHole[In any, Out any]() *Hole[In, Out] {
	return &Hole[In, Out]{}
}

func (hole *Hole[In, Out]) Fill(fnc func(in In) Out) {
	hole.once.Do(func() {
		hole.fnc = fnc
	})
}

func (hole *Hole[In, Out]) Exec(in In) Out {
	return hole.fnc(in)
}
