package fv_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

func TestEvt(t *testing.T) {
	type AAA struct {
		A string
		B int
	}

	typeofAAA := reflect.TypeOf(AAA{})

	bus := fv.EventBus().Workers(1).OnPanic(func(err any, at int64, evttype any, evt any) {
		switch evttype {
		case typeofAAA:
			{
				fmt.Println(">>>>>>>>>>>>>>>>>>>")
				break
			}
		}
		fmt.Println("Panic: ", err, at, evttype, evt)
	}).Build().
		AddListener(typeofAAA, func(at int64, evtany any) {
			evt := (evtany).(*AAA)
			fmt.Println(evt, at)
		})

	evt := AAA{A: "yyy", B: 34}
	bus.Emit(typeofAAA, &evt)
	time.Sleep(time.Second)
}
