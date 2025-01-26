package vld_test

import (
	"context"
	"fmt"
	"testing"

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
	A      int
	B      string
	C      []string
	D      map[int64]string
	EPtr   *E
	ESlice []E
	EMmap  map[string]E
}

func init() {
	lion.AppendType[map[int64]string]()
	lion.AppendType[E]()
	lion.AppendType[map[string]E]()

	vld.SchemeOf[Params]().Scope(func(ctx context.Context, mptr *Params) {
		noemptystring := vld.StringMeta().NoEmpty().Build()

		vld.Int(&mptr.A).Min(1).Max(23).With(ctx)

		vld.String(&mptr.B).NoEmpty().With(ctx)

		vld.Slice(&mptr.C).NoEmpty().
			Ele(noemptystring).
			With(ctx)

		vld.Map(&mptr.D).NoEmpty().
			Ele(vld.StringMeta().RegexpString(`^\d+$`).Build()).
			Key(vld.IntMeta[int64]().Min(12).Build()).
			With(ctx)

		emate := vld.StructMeta[E]().Build()

		vld.Pointer(&mptr.EPtr).Ele(emate).With(ctx)
		vld.Slice(&mptr.ESlice).NoEmpty().Ele(emate).With(ctx)
		vld.Map(&mptr.EMmap).NoEmpty().Ele(emate).Key(noemptystring).With(ctx)
	})

	vld.StringMeta().EnumSlice(vld.AllErrorKinds)
	vld.IntMeta[int64]().EnumSlice(vld.AllErrorKinds)
}

func TestVld(t *testing.T) {
	var params Params
	params.A = 3
	params.B = "xx"
	params.C = []string{"ccc"}
	params.D = map[int64]string{
		12: "3444",
	}
	params.EPtr = &E{
		Email: "a@x.com",
	}
	params.ESlice = []E{
		{Email: "vvv@xdd.com"},
	}
	params.EMmap = map[string]E{
		"xx": {Email: "xxxx@q.com"},
	}
	err := vld.Validate(context.Background(), &params)
	fmt.Println(err)
}

type EPTest struct {
	EV   E
	EPtr *E
}

func init() {
	vld.SchemeOf[EPTest]().Scope(func(ctx context.Context, mptr *EPTest) {
		vld.Struct(&mptr.EV).With(ctx)

		vld.Pointer(&mptr.EPtr).Ele(vld.StructMeta[E]().Build()).With(ctx)
	})
}

func TestEPTest(t *testing.T) {
	ept := EPTest{
		EV:   E{Email: "xxx@qq.com"},
		EPtr: &E{Email: "ff@w.com"},
	}
	fmt.Println(vld.Validate(context.Background(), &ept))
}

type ESTest struct {
	ESlice    []E
	EPtrSlice []*E
}

func init() {
	vld.SchemeOf[ESTest]().Scope(func(ctx context.Context, mptr *ESTest) {
		emate := vld.StructMeta[E]().Build()
		eptrmeta := vld.StructMeta[E]().ToPointer().Build()

		vld.Slice(&mptr.ESlice).NoEmpty().Ele(emate).Func(func(ctx context.Context, v []E) error {
			fmt.Println(v)
			return nil
		}).With(ctx)

		vld.Slice(&mptr.EPtrSlice).NoEmpty().Ele(eptrmeta).Func(func(ctx context.Context, v []*E) error {
			for _, e := range v {
				fmt.Println(e)
			}
			return nil
		}).With(ctx)
	})
}

func TestESTVld(t *testing.T) {
	var est = ESTest{
		ESlice: []E{
			{Email: "yyy@q.com"},
		},
		EPtrSlice: []*E{
			{Email: "xxx@q.com"},
		},
	}
	fmt.Println(vld.Validate(context.Background(), &est))
}

type EMTest struct {
	EValMap map[string]E
	EPtrMap map[string]*E
}

func init() {
	vld.SchemeOf[EMTest]().Scope(func(ctx context.Context, mptr *EMTest) {
		emb := vld.StructMeta[E]()

		vld.Map(&mptr.EValMap).Ele(emb.Build()).With(ctx)

		vld.Map(&mptr.EPtrMap).Ele(emb.ToPointer().Build()).With(ctx)
	})
}

func TestEMTest(t *testing.T) {
	var emt = EMTest{
		EValMap: map[string]E{
			"x": {Email: "xxx@w.com"},
		},
		EPtrMap: map[string]*E{
			"y": {Email: "ccc@s.com"},
		},
	}
	fmt.Println(vld.Validate(context.Background(), &emt))
}
