package evts_test

import (
	"fmt"
	"testing"

	"github.com/zzztttkkk/faceless.void/evts"
	"github.com/zzztttkkk/lion"
)

func TestEvent(t *testing.T) {
	type AddEvt struct {
		A int64
		B int64
	}
	bus := evts.EventBusBuilder().Build()
	bus.AddListener(lion.Typeof[AddEvt](), evts.Wrap(func(at int64, evt *AddEvt) {
		fmt.Println("Type", evt)
	}))
	_ = bus.Emit(&AddEvt{A: 4, B: 6})
	bus.Close(true)
}
