package fv_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	fv "github.com/zzztttkkk/faceless.void"
)

type OnUserCreated struct {
	Id    uint64
	Email string
}

func (o *OnUserCreated) UpdateByContextValue(ctx context.Context) bool {
	return true
}

var _ fv.IEvent = (*OnUserCreated)(nil)

func TestEventBus(T *testing.T) {
	ebus := fv.NewEventBus(nil)

	ebus.Register(reflect.TypeOf(OnUserCreated{}), func(at int64, evt fv.IEvent) {
		eptr := evt.(*OnUserCreated)
		fmt.Println(at, eptr)
	})

	ebus.Emit(context.Background(), &OnUserCreated{Id: 1, Email: "test@test.com"}, nil)
}
