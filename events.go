package fv

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type IEvent interface {
	UpdateByContextValue(ctx context.Context) bool
}

type EventListener func(at int64, evt IEvent)

type _ListenerEle struct {
	wrapped EventListener
	raw     EventListener
}

type EventBus struct {
	lock      sync.RWMutex
	listeners map[reflect.Type][]_ListenerEle
	onpanic   func(err any, at int64, evt any)
}

func NewEventBus(onpanic func(err any, at int64, evt any)) *EventBus {
	return &EventBus{listeners: map[reflect.Type][]_ListenerEle{}, onpanic: onpanic}
}

func (ebus *EventBus) AddListener(evttype reflect.Type, fnc EventListener) *EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()

	ebus.listeners[evttype] = append(ebus.listeners[evttype], _ListenerEle{
		wrapped: func(at int64, evt IEvent) {
			defer func() {
				if ebus.onpanic == nil {
					return
				}
				e := recover()
				if e != nil {
					ebus.onpanic(e, at, evt)
					return
				}
			}()
			fnc(at, evt)
		},
		raw: fnc,
	})
	return ebus
}

func (ebus *EventBus) RemoveListener(evttype reflect.Type, fnc EventListener) *EventBus {
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

func (ebus *EventBus) RemoveAllListener(evttype reflect.Type) *EventBus {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()
	delete(ebus.listeners, evttype)
	return ebus
}

type EventEmitOpts struct {
	Concurrency bool
}

var (
	_DefaultEventEmitOpts = &EventEmitOpts{}
)

func (ebus *EventBus) emit(evttype reflect.Type, evt IEvent, opts *EventEmitOpts) bool {
	var now = time.Now().UnixNano()
	if opts == nil {
		opts = _DefaultEventEmitOpts
	}

	ebus.lock.RLock()
	handlers := ebus.listeners[evttype]
	if len(handlers) < 1 {
		ebus.lock.RUnlock()
		return false
	}
	ebus.lock.RUnlock()

	if !opts.Concurrency {
		for _, fnc := range handlers {
			fnc.wrapped(now, evt)
		}
		return true
	}

	for _, fnc := range handlers {
		go fnc.wrapped(now, evt)
	}
	return true
}

func (ebus *EventBus) Emit(ctx context.Context, evttype reflect.Type, evt IEvent, opts *EventEmitOpts) bool {
	if !evt.UpdateByContextValue(ctx) {
		return false
	}
	return ebus.emit(evttype, evt, opts)
}
