package sqltypes

import "github.com/zzztttkkk/faceless.void/internal"

type intType[T internal.IntType] struct {
	typecommon[intType[T], T]
}

func (it *intType[T]) AutoIncr() *intType[T] {
	it.pairs = append(it.pairs, internal.PairOf("autoincr", true))
	return it
}

func TinyInt(name string) *intType[int8] {
	ins := &intType[int8]{}
	ins.sqltype(name, "int", 8)
	return ins
}

func SmallInt(name string) *intType[int16] {
	ins := &intType[int16]{}
	ins.sqltype(name, "int", 16)
	return ins
}

func Int(name string) *intType[int32] {
	ins := &intType[int32]{}
	ins.sqltype(name, "int", 32)
	return ins
}

func BigInt(name string) *intType[int64] {
	ins := &intType[int64]{}
	ins.sqltype(name, "int", 64)
	return ins
}
