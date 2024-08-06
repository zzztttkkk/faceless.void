package fv

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/zzztttkkk/faceless.void/internal"
)

type difnc struct {
	deps      []reflect.Type
	fnc       reflect.Value
	done      bool
	exectimes int
}

type _DIContainer struct {
	lock         sync.RWMutex
	execed       bool
	fncs         []*difnc
	valpool      map[reflect.Type]reflect.Value
	tokenvalpool map[reflect.Type]map[string]reflect.Value
}

//goland:noinspection ALL
type __fv_private_token_value__ interface {
	__fv_private_token_value__() (string, reflect.Type, reflect.Value)
}

var tokenValueInterfaceType = reflect.TypeOf((*__fv_private_token_value__)(nil)).Elem()

type TokenValue[T any] struct {
	token string
	val   T
}

//goland:noinspection ALL
func (tv TokenValue[T]) __fv_private_token_value__() (string, reflect.Type, reflect.Value) {
	vv := reflect.ValueOf(tv.val)
	return tv.token, vv.Type(), vv
}

func NewTokenValue[T any](token string, val T) TokenValue[T] {
	return TokenValue[T]{token: token, val: val}
}

//goland:noinspection ALL
type __fv_private_token_value_getter__ interface {
	__fv_private_token_value_getter_get_type__() reflect.Type
}

type TokenValueGetter[T any] struct {
	typehint *T
	Fnc      func(string, reflect.Type) any
}

var (
	tokenValueGetterInterfaceType = reflect.TypeOf((*__fv_private_token_value_getter__)(nil)).Elem()
)

//goland:noinspection ALL
func (g TokenValueGetter[T]) __fv_private_token_value_getter_get_type__() reflect.Type {
	return reflect.TypeOf(g.typehint).Elem()
}

func (g TokenValueGetter[T]) Get(token string) T {
	return g.Fnc(token, g.__fv_private_token_value_getter_get_type__()).(T)
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewDIContainer() *_DIContainer {
	return &_DIContainer{
		valpool:      make(map[reflect.Type]reflect.Value),
		tokenvalpool: make(map[reflect.Type]map[string]reflect.Value),
	}
}

func (dic *_DIContainer) errorf(v string, args ...any) error {
	return internal.ErrNamespace{Namespace: "di"}.Errorf(v, args...)
}

func (dic *_DIContainer) appendone(v reflect.Value) {
	var vtype reflect.Type
	var tokenptr *string
	if v.CanConvert(tokenValueInterfaceType) {
		var token string
		token, vtype, v = (v.Interface().(__fv_private_token_value__)).__fv_private_token_value__()
		tokenptr = &token
	} else {
		vtype = v.Type()
	}

	if tokenptr != nil {
		tvm := dic.tokenvalpool[vtype]
		if tvm == nil {
			tvm = make(map[string]reflect.Value)
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

func (dic *_DIContainer) append(v reflect.Value) {
	dic.lock.Lock()
	defer dic.lock.Unlock()

	if v.Kind() == reflect.Slice && v.Type().Elem().AssignableTo(tokenValueInterfaceType) {
		for i := 0; i < v.Len(); i++ {
			dic.appendone(v.Index(i))
		}
		return
	}
	dic.appendone(v)
}

func (dic *_DIContainer) get(k reflect.Type) (reflect.Value, bool) {
	dic.lock.RLock()
	defer dic.lock.RUnlock()
	v, ok := dic.valpool[k]
	return v, ok
}

type errTokenNotFound struct {
	rtype reflect.Type
	token string
}

func (e *errTokenNotFound) Error() string {
	return fmt.Sprintf("fv.di: type `%s` token `%s` not found", e.rtype, e.token)
}

var _ error = (*errTokenNotFound)(nil)

func IsTokenNotFound(e error) bool {
	var ev *errTokenNotFound
	return errors.As(e, &ev)
}

func (dic *_DIContainer) getbytoken(token string, k reflect.Type) reflect.Value {
	dic.lock.RLock()
	defer dic.lock.RUnlock()
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

func (dic *_DIContainer) Prepare(fnc any) {
	rv := reflect.ValueOf(fnc)
	if rv.IsNil() || rv.Kind() != reflect.Func {
		panic(dic.errorf("`%s` is not a function", fnc))
	}

	if rv.Type().NumIn() != 0 {
		panic(dic.errorf("`%s` can not require arguments", fnc))
	}

	if rv.Type().NumOut() < 1 {
		panic(dic.errorf("`%s` has no return value", fnc))
	}

	for _, v := range rv.Call(nil) {
		dic.append(v)
	}
}

func (dic *_DIContainer) Register(fnc any) {
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
}

func (dic *_DIContainer) Run() {
	dic.lock.Lock()
	if dic.execed {
		dic.lock.Unlock()
		panic(dic.errorf("container already executed"))
	}
	dic.execed = true
	dic.lock.Unlock()

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
		for _, ele := range remains {
			ele.exectimes++

			var args []reflect.Value
			for _, argtype := range ele.deps {
				var av reflect.Value
				if argtype.AssignableTo(tokenValueGetterInterfaceType) {
					getter := reflect.New(argtype).Elem()
					getter.FieldByName("Fnc").Set(reflect.ValueOf(func(s string, r reflect.Type) any {
						return dic.getbytoken(s, r).Interface()
					}))
					av = getter
				} else {
					var ok bool
					av, ok = dic.get(argtype)
					if !ok {
						break
					}
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
			panic(dic.errorf("can not resolve dependencies"))
		}
	}

	clear(dic.fncs)
	clear(dic.valpool)
	clear(dic.tokenvalpool)
}
