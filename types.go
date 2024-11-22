package fv

import (
	"fmt"
	"reflect"
	"strings"
)

type _MemOffsetFieldInfo struct {
	offset int64
	field  reflect.StructField
}

type _TypeInfo struct {
	offsets []_MemOffsetFieldInfo
}

var (
	typeinfos = map[reflect.Type]*_TypeInfo{}
)

func addOneType(typ reflect.Type) {
	_, ok := typeinfos[typ]
	if ok {
		return
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Errorf("type `%v` is not a struct type", typ))
	}
	if typ.Name() == "" || strings.ToUpper(typ.Name()[:1]) != typ.Name()[:1] {
		panic(fmt.Errorf("type `%v` is not a public struct type", typ))
	}

	info := &_TypeInfo{}
	typeinfos[typ] = info

	addMemFields(typ, info)
}

func addMemFields(typ reflect.Type, info *_TypeInfo) {
	fc := typ.NumField()

	ev := reflect.New(typ).Elem()
	begin := ev.Addr().Pointer()

	for i := 0; i < fc; i++ {
		ft := typ.Field(i)
		if !ft.IsExported() || ft.Anonymous {
			continue
		}

		fv := ev.FieldByIndex(ft.Index)
		faddr := fv.Addr().Pointer()

		info.offsets = append(info.offsets, _MemOffsetFieldInfo{
			offset: int64(faddr) - int64(begin),
			field:  ft,
		})
	}
}

func RegisterTypes(types ...reflect.Type) {
	for _, typ := range types {
		addOneType(typ)
	}
}
