package account

import (
	"context"
	"encoding/json"
	"hello/sitea/modules/internal"
	"net/http"
	"reflect"

	fv "github.com/zzztttkkk/faceless.void"
	"github.com/zzztttkkk/faceless.void/vld"
)

type RegisterParams struct {
	Name     string
	Email    string
	Password string
}

func init() {
	fv.RegisterTypes(reflect.TypeOf(RegisterParams{}))
}

// Binding implements fv.IBinding.
func (params *RegisterParams) Binding(ctx context.Context) error {
	bnd := fv.Binding(params)
	bnd.String(&params.Email).Validate(vld.Strings.Email().Func())
	bnd.String(&params.Name)
	bnd.String(&params.Password)
	return bnd.Error(ctx)
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
	internal.Delgates.Register.Set(func(ctx context.Context, email, pwd string) (string, error) {
		var params RegisterParams
		params.Email = email
		params.Password = pwd
		r, e := Register(ctx, &params)
		if e != nil {
			return "", nil
		}
		return r.Id, nil
	})

	fv.Endpoint().
		Func(
			func(ctx context.Context, req *http.Request, respw http.ResponseWriter) error {
				var params RegisterParams
				err := params.Binding(ctx)
				if err != nil {
					return err
				}
				result, err := Register(ctx, &params)
				if err != nil {
					return err
				}
				enc := json.NewEncoder(respw)
				return enc.Encode(result)
			},
		).
		Register()
}
