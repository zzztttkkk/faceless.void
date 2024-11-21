package fv_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	fv "github.com/zzztttkkk/faceless.void"
)

type OnUserCreated struct {
	Id    uint64
	Email string
}

var OnUserCreatedType = reflect.TypeOf(OnUserCreated{})

func (o *OnUserCreated) UpdateByContext(ctx context.Context) {
}

var _ fv.IEvent = (*OnUserCreated)(nil)

func TestEventBus(T *testing.T) {
	fnc := func(at int64, evt fv.IEvent) {
		eptr := evt.(*OnUserCreated)
		fmt.Println(at, eptr)
	}

	fv.On(OnUserCreatedType, fnc)

	handled := fv.Emit(context.Background(), OnUserCreatedType, &OnUserCreated{Id: 1, Email: "test@test.com"}, nil)
	fmt.Println(handled)

	time.Sleep(time.Millisecond * 10)
}
