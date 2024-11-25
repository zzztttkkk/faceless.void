package fv

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"unsafe"
)

type ABParams struct {
	A string
	B int16
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
	req.URL.RawQuery = "A=aaaa&B=123"

	ctx := context.WithValue(req.Context(), ctxKeyForHttpRequest, req)
	ctx = getter.init(ctx, req)

	var abv ABParams
	var bnd = Binding(typeofABParams, unsafe.Pointer(&abv))
	bnd.String(&abv.A)
	bnd.Int16(&abv.B)

	fmt.Println(bnd.Error(ctx))
	fmt.Println(abv)
}
