package evts

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type EventListener func(at int64, evtany any)

type _ListenerEle struct {
	wrapped EventListener
	raw     EventListener
}

type eventTask struct {
	evt     any
	at      int64
	handler *_ListenerEle
}

type _EventBus struct {
	lock         sync.RWMutex
	listeners    map[unsafe.Pointer][]_ListenerEle
	onpanic      func(err any, at int64, evt any)
	workers      int
	taskchannels [](chan eventTask)

	closed bool
	wg     sync.WaitGroup
}

type _EventBusBuilder struct {
	pairs []internal.Pair[string]
}

func EventBusBuilder() *_EventBusBuilder {
	return &_EventBusBuilder{}
}

func (builder *_EventBusBuilder) OnPanic(fnc func(err any, at int64, evtptr any)) *_EventBusBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("onpanic", fnc))
	return builder
}

func (builder *_EventBusBuilder) Workers(v int) *_EventBusBuilder {
	builder.pairs = append(builder.pairs, internal.PairOf("workers", v))
	return builder
}

func (builder *_EventBusBuilder) Build() *_EventBus {
	ins := &_EventBus{
		listeners: make(map[unsafe.Pointer][]_ListenerEle),
	}
	for _, pair := range builder.pairs {
		switch pair.Key {
		case "onpanic":
			{
				ins.onpanic = pair.Val.(func(any, int64, any))
				break
			}
		case "workers":
			{
				ins.workers = pair.Val.(int)
				break
			}
		}
	}
	if ins.workers > 0 {
		for i := 0; i < ins.workers; i++ {
			_tc := make(chan eventTask)
			ins.taskchannels = append(ins.taskchannels, _tc)
			go func(tc chan eventTask) {
				for task := range tc {
					task.handler.wrapped(task.at, task.evt)
				}
			}(_tc)
		}
	}
	return ins
}

func evtPointTypeUptr(evttype reflect.Type) unsafe.Pointer {
	if evttype.Kind() != reflect.Struct {
		panic(fmt.Errorf("fv.evts: event type must be a struct, %s", evttype))
	}
	return reflect.ValueOf(reflect.PointerTo(evttype)).UnsafePointer()
}

func (ebus *_EventBus) AddListener(evttype reflect.Type, fnc EventListener) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()

	etypeuptr := evtPointTypeUptr(evttype)
	ebus.listeners[etypeuptr] = append(ebus.listeners[etypeuptr], _ListenerEle{
		wrapped: func(at int64, evt any) {
			ebus.wg.Add(1)
			defer func() {
				ebus.wg.Done()
				if ebus.onpanic != nil {
					if rv := recover(); rv != nil {
						ebus.onpanic(rv, at, evt)
					}
				}
			}()
			fnc(at, evt)
		},
		raw: fnc,
	})
	return ebus
}

func (ebus *_EventBus) RemoveListener(evttype reflect.Type, fnc EventListener) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()

	typeuptr := evtPointTypeUptr(evttype)

	ptr := reflect.ValueOf(fnc).Pointer()
	var nls []_ListenerEle
	for _, l := range ebus.listeners[typeuptr] {
		if reflect.ValueOf(l.raw).Pointer() == ptr {
			continue
		}
		nls = append(nls, l)
	}
	ebus.listeners[typeuptr] = nls
	return ebus
}

func (ebus *_EventBus) RemoveAllListener(evttype reflect.Type) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()
	delete(ebus.listeners, evtPointTypeUptr(evttype))
	return ebus
}

type anyface struct {
	typeptr unsafe.Pointer
	valptr  unsafe.Pointer
}

var (
	ErrEmptyHandlers = errors.New("fv.evts: empty handlers")
	ErrBusClosed     = errors.New("fv.evts: bus closed")
)

func (ebus *_EventBus) Emit(evt any) error {
	var now = time.Now().UnixNano()

	typeptr := (*anyface)(unsafe.Pointer(&evt)).typeptr

	ebus.lock.RLock()
	if ebus.closed {
		ebus.lock.RUnlock()
		return ErrBusClosed
	}
	handlers := ebus.listeners[typeptr]
	if len(handlers) < 1 {
		ebus.lock.RUnlock()
		return ErrEmptyHandlers
	}
	ebus.lock.RUnlock()

	tcc := len(ebus.taskchannels)
	if tcc > 0 {
		for _, fnc := range handlers {
			idx := rand.Intn(tcc)
			ebus.taskchannels[idx] <- eventTask{
				at:      now,
				handler: &fnc,
				evt:     evt,
			}
		}
		return nil
	}

	for _, fnc := range handlers {
		go fnc.wrapped(now, evt)
	}
	return nil
}

func (bus *_EventBus) Close(gracefully bool) {
	bus.lock.Lock()
	bus.closed = true
	bus.lock.Unlock()

	if gracefully {
		runtime.Gosched()
		bus.wg.Wait()
	}
}

func Wrap[T any](handler func(at int64, evt *T)) EventListener {
	return func(at int64, evtany any) {
		valptr := (*anyface)(unsafe.Pointer(&evtany)).valptr
		handler(at, (*T)(valptr))
	}
}
