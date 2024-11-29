package fv

import (
	"fmt"
	"reflect"
	"strings"
)

type _MemOffsetFieldInfo struct {
	offset int64
	field  reflect.StructField
	name   string
}

type _TypeInfo struct {
	type_   reflect.Type
	offsets []_MemOffsetFieldInfo
}

var (
	typeinfos      = map[reflect.Type]*_TypeInfo{}
	commonTagNames = []string{
		"bnd", "binding",
		"json", "toml",
		"db",
		"vld", "validate",
	}
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

	info := &_TypeInfo{
		type_: typ,
	}
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
			name:   tag(&ft, commonTagNames...).Name,
		})
	}
}

func RegisterTypes(types ...reflect.Type) {
	for _, typ := range types {
		addOneType(typ)
	}
}

type Tag struct {
	Name    string
	Options map[string]string
}

func tag(fv *reflect.StructField, names ...string) Tag {
	var tv string
	for _, name := range names {
		tv = fv.Tag.Get(name)
		if tv != "" {
			break
		}
	}

	tag := Tag{}

	parts := strings.Split(tv, ",")
	for idx := range parts {
		parts[idx] = strings.TrimSpace(parts[idx])
	}

	name, remains := parts[0], parts[1:]
	if name == "" {
		name = fv.Name
	}
	tag.Name = name
	for _, v := range remains {
		idx := strings.Index(v, "=")
		if tag.Options == nil {
			tag.Options = map[string]string{}
		}
		if idx > -1 {
			tag.Options[strings.TrimSpace(v[:idx])] = strings.TrimSpace(v[idx+1:])
		} else {
			tag.Options[v] = ""
		}
	}
	return tag
}
