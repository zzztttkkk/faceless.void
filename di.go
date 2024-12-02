package fv

import (
	"fmt"
	"reflect"

	"github.com/zzztttkkk/faceless.void/internal"
)

type difnc struct {
	deps []reflect.Type
	fnc  reflect.Value
	done bool
}

type _DIContainer[T comparable] struct {
	execed       bool
	fncs         []*difnc
	valpool      map[reflect.Type]reflect.Value
	tokenvalpool map[reflect.Type]map[T]reflect.Value
}

//goland:noinspection ALL
type __fv_private_token_value__ interface {
	__fv_private_token_value__() (any, reflect.Type, reflect.Value)
}

var tokenValueInterfaceType = reflect.TypeOf((*__fv_private_token_value__)(nil)).Elem()

type TokenValue[T comparable] struct {
	token T
	val   any
}

func (tv TokenValue[T]) __fv_private_token_value__() (any, reflect.Type, reflect.Value) {
	vv := reflect.ValueOf(tv.val)
	return tv.token, vv.Type(), vv
}

var (
	_ __fv_private_token_value__ = TokenValue[int]{}
)

func ValueWithToken[T comparable](token T, val any) TokenValue[T] {
	return TokenValue[T]{token: token, val: val}
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewContainer[T comparable]() *_DIContainer[T] {
	return &_DIContainer[T]{
		valpool:      make(map[reflect.Type]reflect.Value),
		tokenvalpool: make(map[reflect.Type]map[T]reflect.Value),
	}
}

func (dic *_DIContainer[T]) errorf(v string, args ...any) error {
	return internal.ErrNamespace{Namespace: "di"}.Errorf(v, args...)
}

func (dic *_DIContainer[T]) appendone(v reflect.Value) {
	var vtype reflect.Type
	var tokenptr *T
	if v.Type().Implements(tokenValueInterfaceType) {
		var anytoken any
		anytoken, vtype, v = (v.Interface().(__fv_private_token_value__)).__fv_private_token_value__()
		var token = anytoken.(T)
		tokenptr = &(token)
	} else {
		vtype = v.Type()
	}

	if tokenptr != nil {
		tvm := dic.tokenvalpool[vtype]
		if tvm == nil {
			tvm = make(map[T]reflect.Value)
		}
		_, ok := tvm[*tokenptr]
		if ok {
			panic(dic.errorf("type `%s`, token `%s` is already registered", vtype, *tokenptr))
		}
		tvm[*tokenptr] = v
		dic.tokenvalpool[vtype] = tvm
		return
	}

	_, ok := dic.valpool[vtype]
	if ok {
		panic(dic.errorf("`%s` is already registered", vtype))
	}
	dic.valpool[vtype] = v
}

func (dic *_DIContainer[T]) append(v reflect.Value) {
	if v.Kind() == reflect.Slice && v.Type().Elem().Implements(tokenValueInterfaceType) {
		for i := 0; i < v.Len(); i++ {
			dic.appendone(v.Index(i))
		}
		return
	}
	dic.appendone(v)
}

func (dic *_DIContainer[T]) get(k reflect.Type) (reflect.Value, bool) {
	v, ok := dic.valpool[k]
	return v, ok
}

type errTokenNotFound struct {
	rtype reflect.Type
	token any
}

func (e *errTokenNotFound) Error() string {
	return fmt.Sprintf("fv.di: type `%s` token `%s` not found", e.rtype, e.token)
}

var _ error = (*errTokenNotFound)(nil)

func IsTokenNotFound(e error) bool {
	_, ok := e.(*errTokenNotFound)
	return ok
}

func (dic *_DIContainer[T]) getbytoken(token T, k reflect.Type) reflect.Value {
	tvm, ok := dic.tokenvalpool[k]
	if !ok {
		panic(&errTokenNotFound{token: token, rtype: k})
	}
	v, ok := tvm[token]
	if !ok {
		panic(&errTokenNotFound{token: token, rtype: k})
	}
	return v
}

func (dic *_DIContainer[T]) Register(fnc any) *_DIContainer[T] {
	rv := reflect.ValueOf(fnc)
	if rv.IsNil() || rv.Kind() != reflect.Func {
		panic(dic.errorf("`%s` is not a function", fnc))
	}

	ele := &difnc{fnc: rv}
	for i := 0; i < rv.Type().NumIn(); i++ {
		argtype := rv.Type().In(i)
		ele.deps = append(ele.deps, argtype)
	}
	dic.fncs = append(dic.fncs, ele)
	return dic
}

func (dic *_DIContainer[T]) GetByToken(dest any, token T) *_DIContainer[T] {
	dv := reflect.ValueOf(dest)
	dt := dv.Type().Elem()
	val := dic.getbytoken(token, dt)
	dv.Elem().Set(val)
	return dic
}

func (dic *_DIContainer[T]) Run() {
	if dic.execed {
		panic(dic.errorf("container already executed"))
	}
	dic.execed = true

	for {
		var remains []*difnc
		for _, ele := range dic.fncs {
			if ele.done {
				continue
			}
			remains = append(remains, ele)
		}

		if len(remains) < 1 {
			break
		}

		donecount := 0
		var lastTokenMissingErr *errTokenNotFound
		var lastTokenMissingFunc reflect.Value
		for _, ele := range remains {
			var args []reflect.Value
			for _, argtype := range ele.deps {
				var av reflect.Value
				var ok bool
				av, ok = dic.get(argtype)
				if !ok {
					break
				}
				args = append(args, av)
			}

			call := func() (rvs []reflect.Value, ok bool) {
				defer func() {
					rv := recover()
					if rv == nil {
						return
					}
					e, a := rv.(error)
					if a && IsTokenNotFound(e) {
						ok = false
						rvs = nil
						lastTokenMissingErr = e.(*errTokenNotFound)
						lastTokenMissingFunc = ele.fnc
						return
					}
					panic(rv)
				}()
				return ele.fnc.Call(args), true
			}

			if len(args) == len(ele.deps) {
				rvs, ok := call()
				if !ok {
					continue
				}
				for _, rv := range rvs {
					dic.append(rv)
				}
				ele.done = true
				donecount++
			}
		}

		if donecount < 1 {
			// TODO
			if lastTokenMissingErr != nil {
				fmt.Println(lastTokenMissingErr.Error(), lastTokenMissingFunc)
			}
			for _, ele := range remains {
				if ele.done || ele.fnc == lastTokenMissingFunc {
					continue
				}
				fmt.Println(ele.fnc)
			}
			panic(fmt.Errorf("fv.dic: can not resolve dependencies"))
		}
	}
}
