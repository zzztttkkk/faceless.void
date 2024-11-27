package fv

import (
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/zzztttkkk/faceless.void/internal"
)

type EventListener func(at int64, evt any)

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
	listeners    map[any][]_ListenerEle
	onpanic      func(err any, at int64, evttype any, evt any)
	workers      int
	taskchannels [](chan eventTask)
}

type _EventBusBuilder struct {
	pairs []internal.Pair[string]
}

func EventBusBuilder() *_EventBusBuilder {
	return &_EventBusBuilder{}
}

func (builder *_EventBusBuilder) OnPanic(fnc func(err any, at int64, evttype any, evtptr any)) *_EventBusBuilder {
	builder.pairs = append(builder.pairs, internal.Pair[string]{Key: "onpanic", Val: fnc})
	return builder
}

func (builder *_EventBusBuilder) Workers(v int) *_EventBusBuilder {
	builder.pairs = append(builder.pairs, internal.Pair[string]{Key: "workers", Val: v})
	return builder
}

func (builder *_EventBusBuilder) Build() *_EventBus {
	ins := &_EventBus{
		listeners: make(map[any][]_ListenerEle),
	}
	for _, pair := range builder.pairs {
		switch pair.Key {
		case "onpanic":
			{
				ins.onpanic = pair.Val.(func(any, int64, any, any))
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

func (ebus *_EventBus) AddListener(evttype any, fnc EventListener) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()

	wrapped := fnc
	if ebus.onpanic != nil {
		wrapped = func(at int64, evt any) {
			defer func() {
				e := recover()
				if e != nil {
					ebus.onpanic(e, at, evttype, evt)
				}
			}()
			fnc(at, evt)
		}
	}
	ebus.listeners[evttype] = append(ebus.listeners[evttype], _ListenerEle{
		wrapped: wrapped,
		raw:     fnc,
	})
	return ebus
}

func (ebus *_EventBus) RemoveListener(evttype any, fnc EventListener) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()

	ptr := reflect.ValueOf(fnc).Pointer()
	var nls []_ListenerEle
	for _, l := range ebus.listeners[evttype] {
		if reflect.ValueOf(l.raw).Pointer() == ptr {
			continue
		}
		nls = append(nls, l)
	}
	ebus.listeners[evttype] = nls
	return ebus
}

func (ebus *_EventBus) RemoveAllListener(evttype any) *_EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()
	delete(ebus.listeners, evttype)
	return ebus
}

func (ebus *_EventBus) Emit(evttype any, evt any) {
	var now = time.Now().UnixNano()

	ebus.lock.RLock()
	handlers := ebus.listeners[evttype]
	if len(handlers) < 1 {
		ebus.lock.RUnlock()
		return
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
		return
	}

	for _, fnc := range handlers {
		go fnc.wrapped(now, evt)
	}
	return
}
