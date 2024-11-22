package fv

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _FieldNames struct {
	begin    int64
	tagname  string
	typeinfo *_TypeInfo
}

func NewFieldNames(t reflect.Type, ptr unsafe.Pointer, tagname string) _FieldNames {
	typeinfo, ok := typeinfos[t]
	if !ok {
		panic(fmt.Errorf("type `%v` is not found, please register it", t))
	}
	begin := (int64)((uintptr)(ptr))
	return _FieldNames{
		begin:    begin,
		tagname:  tagname,
		typeinfo: typeinfo,
	}
}

func (names *_FieldNames) Name(ptr unsafe.Pointer) string {
	var offset = (int64)((uintptr)(ptr)) - (int64)(names.begin)
	for idx := range names.typeinfo.offsets {
		ele := &names.typeinfo.offsets[idx]
		if ele.offset == offset {
			return ele.field.Name
		}
	}
	panic(fmt.Errorf("can not found field info when offset is %d", offset))
}
