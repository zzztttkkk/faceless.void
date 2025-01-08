package vld

import (
	"fmt"
	"reflect"

	"github.com/zzztttkkk/lion"
)

type _StructBuilder[T any] struct {
	_CommonBuilder[T, _StructBuilder[T]]
}

func getSchemeByType(gotype reflect.Type) _IScheme {
	ss, ok := schemes[gotype]
	if !ok {
		panic(fmt.Errorf("fv.vld: empty scheme for `%s`", gotype))
	}
	return ss
}

func StructMeta[T any]() *_StructBuilder[T] {
	gotype := lion.Typeof[T]()
	if gotype.Kind() != reflect.Struct {
		panic(fmt.Errorf("fv.vld: %s is not a struct type", gotype))
	}
	obj := &_StructBuilder[T]{}
	return obj.set("scheme", getSchemeByType(gotype))
}

func Struct[T any](ptr *T) *_StructBuilder[T] {
	return StructMeta[T]().updateptr(ptr)
}

func (builder _StructBuilder[T]) ToPointer() *_PointerBuilder[*T] {
	return PointerMeta[T]().Ele(builder.Build())
}
