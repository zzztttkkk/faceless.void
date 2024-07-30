package fv

import (
	"context"
	"reflect"
	"time"
)

type IEvent interface {
	UpdateByContextValue(ctx context.Context) bool
}

type EventListener func(at int64, evt IEvent)

type EventBus struct {
	listeners map[reflect.Type][]EventListener
	onpanic   func(err any, at int64, evt any)
}

func NewEventBus(onpanic func(err any, at int64, evt any)) *EventBus {
	return &EventBus{listeners: map[reflect.Type][]EventListener{}, onpanic: onpanic}
}

func (ebus *EventBus) AddListener(evttype reflect.Type, fnc EventListener) {
	ebus.listeners[evttype] = append(ebus.listeners[evttype], func(at int64, evt IEvent) {
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
	})
}

type EventEmitOpts struct {
	Concurrency bool
}

var (
	_DefaultEventEmitOpts = &EventEmitOpts{}
)

func (ebus *EventBus) emit(evttype reflect.Type, evt IEvent, opts *EventEmitOpts) {
	var now = time.Now().UnixNano()

	handlers := ebus.listeners[evttype]
	if len(handlers) < 1 {
		panic("no listener")
	}

	if opts == nil {
		opts = _DefaultEventEmitOpts
	}

	if !opts.Concurrency {
		for _, fnc := range handlers {
			fnc(now, evt)
		}
		return
	}

	for _, fnc := range handlers {
		go fnc(now, evt)
	}
}

func (ebus *EventBus) Emit(ctx context.Context, evt IEvent, opts *EventEmitOpts) {
	if !evt.UpdateByContextValue(ctx) {
		return
	}

	ev := reflect.ValueOf(evt)
	if ev.Kind() != reflect.Pointer {
		panic("event must be a struct pointer")
	}
	ev = ev.Elem()
	if ev.Kind() != reflect.Struct {
		panic("event must be a struct pointer")
	}

	ebus.emit(ev.Type(), evt, opts)
}

func (ebus *EventBus) EmitWithType(ctx context.Context, evttype reflect.Type, evt IEvent, opts *EventEmitOpts) {
	if !evt.UpdateByContextValue(ctx) {
		return
	}
	ebus.emit(evttype, evt, opts)
}
