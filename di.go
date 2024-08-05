package fv

import (
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
	lock    sync.RWMutex
	execed  bool
	fncs    []*difnc
	valpool map[reflect.Type]reflect.Value
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewDIContainer() *_DIContainer {
	return &_DIContainer{
		valpool: make(map[reflect.Type]reflect.Value),
	}
}

func (dic *_DIContainer) errorf(v string, args ...any) error {
	return internal.ErrNamespace{Namespace: "di"}.Errorf(v, args...)
}

func (dic *_DIContainer) append(v reflect.Value) {
	dic.lock.Lock()
	defer dic.lock.Unlock()

	k := v.Type()
	_, ok := dic.valpool[k]
	if ok {
		panic(dic.errorf("`%s` is already registered", k))
	}
	dic.valpool[k] = v
}

func (dic *_DIContainer) get(k reflect.Type) (reflect.Value, bool) {
	dic.lock.RLock()
	defer dic.lock.RUnlock()
	v, ok := dic.valpool[k]
	return v, ok
}

func (dic *_DIContainer) Pre(fnc any) {
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
		ele.deps = append(ele.deps, rv.Type().In(i))
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
				av, ok := dic.get(argtype)
				if !ok {
					break
				}
				args = append(args, av)
			}

			if len(args) == len(ele.deps) {
				for _, rv := range ele.fnc.Call(args) {
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
}
