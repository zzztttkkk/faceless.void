package evts_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/zzztttkkk/faceless.void/evts"
	"github.com/zzztttkkk/lion"
)

func TestEvent(t *testing.T) {
	bus := evts.EventBusBuilder().Build()

	type AddEvt struct {
		A int64
		B int64
	}

	evts.Register(bus, func(at int64, evt *AddEvt) {
		fmt.Println("Type", evt)
	})

	bus.AddListener(lion.Typeof[AddEvt](), func(at int64, evtany any) {
		time.Sleep(time.Second)
		fmt.Println("Any", evtany)
	})

	fmt.Println(bus.Emit(&AddEvt{A: 4, B: 6}))

	bus.Close(true)
}
