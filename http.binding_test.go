package fv

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"unsafe"

	"github.com/zzztttkkk/faceless.void/internal"
)

type ABParams struct {
	A string
	B int16
	C []string
	D []int32
}

var (
	typeofABParams = reflect.TypeOf(ABParams{})
)

func init() {
	RegisterTypes(typeofABParams)
}

func TestBinding(t *testing.T) {
	var getter _Getter
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.URL.RawQuery = "A=aaaa&B=123&C=a&C=b&D=1&D=456"

	ctx := context.WithValue(req.Context(), internal.CtxKeyForHttpRequest, req)
	ctx = getter.init(ctx, req)

	var abv ABParams
	var bnd = BindingWithType(typeofABParams, unsafe.Pointer(&abv))
	bnd.String(&abv.A)
	bnd.Int16(&abv.B)
	bnd.Strings(&abv.C, nil)
	bnd.Int32Slice(&abv.D)

	fmt.Println(bnd.Error(ctx))
	fmt.Println(abv)
}
