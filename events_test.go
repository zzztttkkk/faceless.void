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

var OnUserCreatedType = reflect.TypeOf(OnUserCreated{})

func (o *OnUserCreated) UpdateByContextValue(ctx context.Context) bool {
	return true
}

var _ fv.IEvent = (*OnUserCreated)(nil)

func TestEventBus(T *testing.T) {
	ebus := fv.NewEventBus(nil)

	fnc := func(at int64, evt fv.IEvent) {
		eptr := evt.(*OnUserCreated)
		fmt.Println(at, eptr)
	}

	ebus.AddListener(OnUserCreatedType, fnc)
	ebus.RemoveListener(OnUserCreatedType, fnc)

	ebus.Emit(context.Background(), OnUserCreatedType, &OnUserCreated{Id: 1, Email: "test@test.com"}, nil)
}
