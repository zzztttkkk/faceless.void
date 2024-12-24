package internalvld

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/zzztttkkk/lion"
)

type Params struct {
	A int
	B string
	C []string
	D map[int64]string
}

func init() {
	lion.AppendType[map[int64]string]()

	mptr := Ptr[Params]()
	SchemeOf[Params]().
		Field(&mptr.A, &VldFieldMeta{MinInt: sql.Null[int64]{V: 0, Valid: true}}).
		Field(&mptr.B, &VldFieldMeta{}).
		Field(&mptr.C, &VldFieldMeta{MinLength: sql.Null[int]{V: 1, Valid: true}, Optional: true}).
		Field(&mptr.D, &VldFieldMeta{Regexp: regexp.MustCompile(`^\d+$`), MapKey: &VldFieldMeta{MinInt: sql.Null[int64]{V: 12, Valid: true}}}).
		Finish()
}

func TestVld(t *testing.T) {
	var params Params

	params.D = map[int64]string{
		12: "3444r",
	}
	err := Vld(&params)
	fmt.Println(err)
}
