package vld_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/vld"
	"github.com/zzztttkkk/lion"
)

type E struct {
	Email string
	A2    string
}

func init() {
	vld.SchemeOf[E]().Scope(func(ctx context.Context, mptr *E) {
		vld.String(&mptr.Email).Email().With(ctx)
	})
}

type Params struct {
	A  int
	B  string
	C  []string
	D  map[int64]string
	EP *E
	ES []E
	EM map[string]E
}

func init() {
	lion.AppendType[map[int64]string]()
	lion.AppendType[E]()

	vld.SchemeOf[Params]().Scope(func(ctx context.Context, mptr *Params) {
		vld.Int(&mptr.A).Min(1).Max(23).With(ctx)

		vld.String(&mptr.B).NoEmpty().With(ctx)

		vld.Slice(&mptr.C).NoEmpty().
			Ele(vld.StringMeta().NoEmpty().Build()).
			With(ctx)

		vld.Map(&mptr.D).
			Ele(vld.StringMeta().RegexpString(`^\d+$`).Build()).
			Key(vld.IntMeta[int64]().Min(12).Build()).
			With(ctx)

		vld.Struct(&mptr.EP).With(ctx)
		vld.Slice(&mptr.ES).NoEmpty().Ele(vld.StructMeta[E]().Build()).With(ctx)
		vld.Map(&mptr.EM).Ele(vld.StructMeta[E]().Build()).With(ctx)
	})
}

func TestVld(t *testing.T) {
	var params Params
	params.A = 3
	params.B = "xx"
	params.C = []string{"ccc"}
	params.D = map[int64]string{
		13: "3444",
	}
	params.EP = &E{
		Email: "a@x.com",
	}
	params.ES = []E{
		{Email: "vvv@xdd.com"},
	}
	params.EM = map[string]E{
		"xx": {Email: "xxxx@q.com"},
	}
	err := vld.Vld(context.Background(), &params)
	fmt.Println(err)
}

func TestSlice(t *testing.T) {
	var x []int64 = []int64{1212, 456, 67}

	eachslice(unsafe.Pointer(&x), lion.Typeof[int64](), func(euptr unsafe.Pointer) {
		fmt.Println(euptr, *((*int64)(euptr)))
	})

	var es []*E = []*E{
		{Email: "xxxx"},
		{A2: "a2"},
	}
	eachslice(unsafe.Pointer(&es), lion.Typeof[*E](), func(euptr unsafe.Pointer) {
		fmt.Println(euptr, *((**E)(euptr)))
	})

	eachslicet(&es, func(eptr **E) {
		fmt.Println(*eptr)
	})
}

func eachslice(sliceptr unsafe.Pointer, eletype reflect.Type, elefnc func(euptr unsafe.Pointer)) {
	sh := *(*reflect.SliceHeader)(sliceptr)
	if sh.Data == 0 {
		panic("nil slice")
	}
	elesize := eletype.Size()
	begin := unsafe.Pointer(sh.Data)
	for i := 0; i < sh.Len; i++ {
		elefnc(unsafe.Add(begin, i*int(elesize)))
	}
}

func eachslicet[T any](sliceptr *[]T, elefnc func(eptr *T)) {
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(sliceptr))
	if sh.Data == 0 {
		panic("nil slice")
	}
	elesize := (lion.Typeof[T]()).Size()
	begin := unsafe.Pointer(sh.Data)
	for i := 0; i < sh.Len; i++ {
		elefnc((*T)(unsafe.Add(begin, i*int(elesize))))
	}
}
