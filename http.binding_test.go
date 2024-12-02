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
	A string `bnd:",aa=we,b=45,c"`
	B int16  `bnd:"b"`
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
	var bnd = BindingWithType(unsafe.Pointer(&abv), typeofABParams)
	bnd.String(&abv.A)
	bnd.Int16(&abv.B)
	bnd.Strings(&abv.C, nil)
	bnd.Int32Slice(&abv.D)

	fmt.Println(bnd.Error(ctx))
	fmt.Println(abv, ErrorKindBindingMissingRequired, ErrorKindBindingParseFailed)

	var f1 func()

	fmt.Println(reflect.ValueOf(f1).IsNil())
}
