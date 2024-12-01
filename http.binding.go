package fv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type IBinding interface {
	Binding(ctx context.Context) error
}

var (
	iBindingType = reflect.TypeOf((*IBinding)(nil)).Elem()
)

type _Getter struct {
	ctx     context.Context
	req     *http.Request
	cookies []*http.Cookie
}

func (getter *_Getter) init(ctx context.Context, req *http.Request) context.Context {
	ctx = context.WithValue(ctx, internal.CtxKeyForBindingGetter, getter)
	getter.ctx = ctx
	getter.req = req
	return ctx
}

type BindingSrcKind int

const (
	BindingSrcForm = BindingSrcKind(iota)
	BindingSrcHeader
	BindingSrcCookie
	BindingSrcQuery
	BindingSrcPostForm
	BindingSrcPathValue
	BindingSrcOsEnv
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func (getter *_Getter) getvalues(where BindingSrcKind, key string) ([]string, bool) {
	var vs []string
	var ok bool
	switch where {
	case BindingSrcQuery:
		{
			vs, ok = getter.req.URL.Query()[key]
			break
		}
	case BindingSrcPathValue:
		{
			v := getter.req.PathValue(key)
			if v != "" {
				vs = append(vs, v)
				ok = true
			}
		}
	case BindingSrcHeader:
		{
			vs, ok = getter.req.Header[key]
			break
		}
	case BindingSrcCookie:
		{
			if getter.cookies == nil {
				getter.cookies = getter.req.Cookies()
			}
			for _, item := range getter.cookies {
				if item.Name == key {
					vs = append(vs, item.Value)
					ok = true
				}
			}
			break
		}
	case BindingSrcForm:
		{
			getter.req.ParseMultipartForm(defaultMaxMemory)
			vs, ok = getter.req.Form[key]
			break
		}
	case BindingSrcPostForm:
		{
			getter.req.ParseMultipartForm(defaultMaxMemory)
			vs, ok = getter.req.PostForm[key]
			break
		}
	case BindingSrcOsEnv:
		{
			v, ok := os.LookupEnv(key)
			if ok {
				vs = append(vs, v)
				return vs, true
			}
			return nil, false
		}
	}
	return vs, ok
}

func (getter *_Getter) getvaluesbynames(where BindingSrcKind, key string, alias ...string) ([]string, bool) {
	vs, ok := getter.getvalues(where, key)
	if ok {
		return vs, true
	}
	for _, name := range alias {
		vs, ok = getter.getvalues(where, name)
		if ok {
			return vs, true
		}
	}
	return nil, false
}

func (getter *_Getter) Int(where BindingSrcKind, base int, key string, alias ...string) (int64, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return 0, false
	}
	iv, err := strconv.ParseInt(vs[0], base, 64)
	if err != nil {
		return 0, false
	}
	return iv, true
}

func (getter *_Getter) Ints(where BindingSrcKind, base int, key string, alias ...string) ([]int64, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok {
		return nil, false
	}
	var ints = make([]int64, 0, len(vs))
	for _, v := range vs {
		iv, err := strconv.ParseInt(v, base, 64)
		if err != nil {
			return nil, false
		}
		ints = append(ints, iv)
	}
	return ints, true
}

func (getter *_Getter) Uint(where BindingSrcKind, key string, alias ...string) (uint64, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return 0, false
	}
	iv, err := strconv.ParseUint(vs[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return iv, true
}

func (getter *_Getter) Uints(where BindingSrcKind, base int, key string, alias ...string) ([]uint64, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok {
		return nil, false
	}
	var ints = make([]uint64, 0, len(vs))
	for _, v := range vs {
		iv, err := strconv.ParseUint(v, base, 64)
		if err != nil {
			return nil, false
		}
		ints = append(ints, iv)
	}
	return ints, true
}

func (getter *_Getter) Bool(where BindingSrcKind, key string, alias ...string) (bool, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return false, false
	}
	iv, err := strconv.ParseBool(vs[0])
	if err != nil {
		return false, false
	}
	return iv, true
}

func (getter *_Getter) Bools(where BindingSrcKind, key string, alias ...string) ([]bool, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return nil, false
	}
	bools := make([]bool, len(vs))
	for idx, v := range vs {
		iv, err := strconv.ParseBool(v)
		if err != nil {
			return nil, false
		}
		bools[idx] = iv
	}
	return bools, true
}

func (getter *_Getter) Time(where BindingSrcKind, layout string, key string, alias ...string) (time.Time, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return time.Time{}, false
	}
	tv, err := time.Parse(layout, vs[0])
	if err != nil {
		return time.Time{}, false
	}
	return tv, true
}

func (getter *_Getter) Times(where BindingSrcKind, layout string, key string, alias ...string) ([]time.Time, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return nil, false
	}
	times := make([]time.Time, len(vs))
	for idx, v := range vs {
		tv, err := time.Parse(layout, v)
		if err != nil {
			return nil, false
		}
		times[idx] = tv
	}
	return times, true
}

func (getter *_Getter) String(where BindingSrcKind, key string, alias ...string) (string, bool) {
	vs, ok := getter.getvaluesbynames(where, key, alias...)
	if !ok || len(vs) < 1 {
		return "", false
	}
	return vs[0], true
}

func (getter *_Getter) Strings(where BindingSrcKind, key string, alias ...string) ([]string, bool) {
	return getter.getvaluesbynames(where, key, alias...)
}

type _BindingInstance struct {
	ptr    int64
	info   *_TypeInfo
	fields [](func(context.Context) error)
}

func (ins *_BindingInstance) Error(ctx context.Context) error {
	for _, fnc := range ins.fields {
		err := fnc(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ins *_BindingInstance) nameof(ptr unsafe.Pointer) string {
	offset := int64(uintptr(ptr)) - ins.ptr
	for idx := range ins.info.offsets {
		ele := &ins.info.offsets[idx]
		if ele.offset == offset {
			return ele.name
		}
	}
	panic(fmt.Errorf("fv.binding: can not find field info by offset(%d), %v", offset, ins.info.type_))
}

func BindingWithType(ptr unsafe.Pointer, vtype reflect.Type) _BindingInstance {
	return _BindingInstance{
		ptr:  int64(uintptr(ptr)),
		info: typeinfos[vtype],
	}
}

func Binding[T any](ptr *T) _BindingInstance {
	return BindingWithType(unsafe.Pointer(ptr), reflect.TypeOf(ptr).Elem())
}
