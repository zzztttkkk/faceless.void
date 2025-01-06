package vld_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/zzztttkkk/faceless.void/vld"
	"github.com/zzztttkkk/lion"
)

type E struct {
	A1 string
	A2 string
}

func init() {
	vld.SchemeOf[E]().Scope(func(ctx context.Context, mptr *E) {})
}

type Params struct {
	A  int
	B  string
	C  []string
	D  map[int64]string
	ES []*E
	E  *E
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
	err := vld.Vld(context.Background(), &params)
	fmt.Println(err)
}
