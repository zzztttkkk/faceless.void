package fv

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type IBinding interface {
	Binding(ctx context.Context, req *http.Request) error
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
	ctx = context.WithValue(ctx, ctxKeyForBindingGetter, getter)
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
	}
	return vs, ok
}

func (getter *_Getter) Int(where BindingSrcKind, key string) (int64, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return 0, false
	}
	iv, err := strconv.ParseInt(vs[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return iv, true
}

func (getter *_Getter) MustInt(where BindingSrcKind, key string, validator func(int64) error) int64 {
	var err error
	iv, ok := getter.Int(where, key)
	if !ok {
		err = fmt.Errorf("missing required params: %v, %s", where, key)
	} else {
		if validator != nil {
			err = validator(iv)
		}
	}
	if err != nil {
		panic(err)
	}
	return iv
}

func (getter *_Getter) Uint(where BindingSrcKind, key string) (uint64, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return 0, false
	}
	iv, err := strconv.ParseUint(vs[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return iv, true
}

func (getter *_Getter) Bool(where BindingSrcKind, key string) (bool, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return false, false
	}
	iv, err := strconv.ParseBool(vs[0])
	if err != nil {
		return false, false
	}
	return iv, true
}

func (getter *_Getter) Time(where BindingSrcKind, key string, layout string) (time.Time, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return time.Time{}, false
	}
	tv, err := time.Parse(layout, vs[0])
	if err != nil {
		return time.Time{}, false
	}
	return tv, true
}

func (getter *_Getter) Duration(where BindingSrcKind, key string) (time.Duration, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return 0, false
	}
	tv, err := time.ParseDuration(vs[0])
	if err != nil {
		return 0, false
	}
	return tv, true
}

func (getter *_Getter) String(where BindingSrcKind, key string) (string, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return "", false
	}
	return vs[0], true
}

func (getter *_Getter) MustString(where BindingSrcKind, key string, validator func(string) error) string {
	var err error
	sv, ok := getter.String(where, key)
	if !ok {
		err = fmt.Errorf("")
	} else {
		if validator != nil {
			err = validator(sv)
		}
	}
	if err != nil {
		panic(err)
	}
	return sv
}

func (getter *_Getter) Any(where BindingSrcKind, key string, parse func(string, ...any) (any, error), args ...any) (any, bool) {
	vs, ok := getter.getvalues(where, key)
	if !ok || len(vs) < 1 {
		return nil, false
	}
	av, err := parse(vs[0], args...)
	if err != nil {
		return nil, false
	}
	return av, true
}

func (getter *_Getter) Unmarshal(dest any) error {
	return nil
}

type _InstanceGetter struct {
	getter *_Getter
	names  _FieldNames
}

func (obj *_InstanceGetter) String(dest *string, where BindingSrcKind, validator func(string) error) {
	key := obj.names.Name(unsafe.Pointer(dest))
	*dest = obj.getter.MustString(where, key, validator)
}

func (getter *_Getter) Instance(t reflect.Type, ptr unsafe.Pointer) *_InstanceGetter {
	return &_InstanceGetter{
		getter: getter,
		names:  NewFieldNames(t, ptr, ""),
	}
}
