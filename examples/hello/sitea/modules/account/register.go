package account

import (
	"context"
	"net/http"
	"reflect"
	"unsafe"

	fv "github.com/zzztttkkk/faceless.void"
	"github.com/zzztttkkk/faceless.void/vld"
)

type RegisterParams struct {
	Name     string
	Email    string
	Password string
}

var (
	typeOfRegisterParams = reflect.TypeOf(RegisterParams{})
)

func init() {
	fv.RegisterTypes(typeOfRegisterParams)
}

// Binding implements fv.IBinding.
func (params *RegisterParams) Binding(ctx context.Context, req *http.Request) error {
	bnd := fv.BindingGetter(ctx).Instance(typeOfRegisterParams, unsafe.Pointer(params))
	bnd.String(&params.Password, fv.BindingSrcForm, vld.String().MinLen(1).MaxLen(10).Finish())
	return nil
}

var _ fv.IBinding = (*RegisterParams)(nil)

type RegisterResult struct {
	Id   string
	Name string
}

func Register(ctx context.Context, params *RegisterParams) (*RegisterResult, error) {
	return nil, nil
}

func init() {
	fv.Endpoint().Register(Register)
}
