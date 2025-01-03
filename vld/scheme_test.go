package vld_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/zzztttkkk/faceless.void/vld"
	"github.com/zzztttkkk/lion"
)

type E struct {
	A1 string
	A2 string
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

	scheme := vld.SchemeOf[Params]()
	defer scheme.Finish()

	model := vld.Ptr[Params]()

	scheme.
		Field(&model.B, vld.String().NoEmpty().Build()).
		Field(&model.C, &vld.VldFieldMeta{MinLength: sql.Null[int]{V: 1, Valid: true}, Optional: true}).
		Field(&model.D, &vld.VldFieldMeta{Regexp: regexp.MustCompile(`^\d+$`), MapKey: &vld.VldFieldMeta{MinInt: sql.Null[int64]{V: 12, Valid: true}}})

	vld.IntPtr(&model.A).Min(1).Max(23).Finish(scheme)
}

func TestVld(t *testing.T) {
	var params Params
	params.D = map[int64]string{
		12: "3444r",
	}
	err := vld.Vld(context.Background(), &params)
	fmt.Println(err)
}
